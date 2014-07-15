package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type appContext struct {
	config *appConfig
	gin    *gin.Engine
	// Lots of global state, but the app only needs to handle one deploy at once
	deployStatus   string
	logBuffer      *bytes.Buffer
	currentCommand string
}

type infoResponse struct {
	Status  string `json:"status"`
	Log     string `json:"log"`
	Command string `json:"command"`
}

func NewApp(config *appConfig) (*appContext, error) {
	app := new(appContext)
	app.config = config

	app.deployStatus = "Not deploying"

	app.gin = gin.Default()
	app.gin.LoadHTMLTemplates("templates/*")
	app.gin.GET("/", app.rootHandler)
	app.gin.GET("/status", app.statusHandler)
	app.gin.GET("/info", app.infoHandler)
	app.gin.POST("/", app.formHandler)
	path, _ := os.Getwd()
	path = filepath.Join(path, "resources")
	app.gin.Static("/resources", path)

	app.logBuffer = bytes.NewBufferString("")
	app.currentCommand = "N/A"

	log.Println(config)
	return app, nil
}

func (app *appContext) Start() {
	// TODO: load from config
	app.gin.Run(":8081")
}

func (app *appContext) rootHandler(ctx *gin.Context) {
	envDropdowns := ""
	hostDropdowns := ""
	for _, env := range app.config.Environments {
		base := `<option value="` + env + `">` + env + `</option>` + "\n"
		envDropdowns += base
	}
	for _, node := range app.config.Nodes {
		base := `<option value"` + node.Hostname + `">` + node.Hostname + `</option>` + "\n"
		hostDropdowns += base
	}

	obj := gin.H{
		"envDropdowns":  template.HTML(envDropdowns),
		"hostDropdowns": template.HTML(hostDropdowns),
	}
	ctx.HTML(200, "index.tmpl", obj)
}

func (app *appContext) formHandler(ctx *gin.Context) {
	ctx.Req.ParseForm()
	log.Println(ctx.Req.Form)

	env, ok := ctx.Req.Form["environments"]
	if !ok || len(env) == 0 {
		http.Error(ctx.Writer, "Environment field missing", 400)
		return
	}

	host, ok := ctx.Req.Form["nodes"]
	if !ok || len(host) == 0 {
		http.Error(ctx.Writer, "Node field missing", 400)
		return
	}

	configJson, ok := ctx.Req.Form["json"]
	if !ok || len(configJson) == 0 {
		configJson = []string{""}
	}

	err := app.Deploy(env[0], host[0], configJson[0])
	if err != nil {
		log.Println(err)
		http.Error(ctx.Writer, "Error parsing JSON: "+err.Error(), 400)
	}
	http.Redirect(ctx.Writer, ctx.Req, "/status", 303)
}

func (app *appContext) statusHandler(ctx *gin.Context) {
	obj := gin.H{"deployStatus": app.deployStatus}
	ctx.HTML(200, "status.tmpl", obj)
}

func (app *appContext) infoHandler(ctx *gin.Context) {
	logstr := app.logBuffer.String()
	out := ""
	lines := strings.Split(logstr, "\n")
	for _, line := range lines {
		out += "<p>" + line + "</p>\n"
	}
	resp := infoResponse{
		Status:  app.deployStatus,
		Log:     out,
		Command: app.currentCommand,
	}
	ctx.JSON(200, resp)
}

func (app *appContext) Deploy(environment, host, configJson string) error {
	var minJson string
	if len(configJson) == 0 {
		minJson = ""
	} else {
		var err error
		minJson, err = minifyJson(configJson)
		if err != nil {
			return err
		}
	}
	_ = minJson
	// Use dummy script for testing
	//cmdStr := fmt.Sprintf("knife bootstrap -E %s %s -r %s chef-full --json-attributes '%s' -x %s --sudo -c %s", environment,
	//		host, buildRecipes(app.config.Recipes), minJson, app.config.User, app.config.KnifeRb)
	cmdStr := filepath.Join(os.Getenv("GOPATH"), "/src/github.com/jherman3/preview_deploy/fake_command.sh")
	app.currentCommand = cmdStr
	log.Println(cmdStr)
	cmd := exec.Command(cmdStr)
	go app.processCommand(cmd)
	return nil
}

func (app *appContext) processCommand(cmd *exec.Cmd) {
	app.deployStatus = "In progress"

	cmd.Stdout = app.logBuffer
	cmd.Stderr = app.logBuffer

	cmd.Start()
	go func() {
		cmd.Wait()
		app.deployStatus = "Completed"
	}()
	for app.deployStatus == "In progress" {

	}
}

func minifyJson(configJson string) (string, error) {
	dest := bytes.NewBufferString("")
	err := json.Compact(dest, []byte(configJson))
	if err != nil {
		return "", err
	}
	return dest.String(), nil
}

func buildRecipes(recipes []string) string {
	out := `"`
	for i, recipe := range recipes {
		out += "recipe[" + recipe + "]"
		if i < len(recipes)-1 {
			out += ","
		}
	}
	out += `"`
	return out
}
