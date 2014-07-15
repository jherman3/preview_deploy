package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
)

type appContext struct {
	config *appConfig
	gin    *gin.Engine
}

func NewApp(config *appConfig) (*appContext, error) {
	app := new(appContext)
	app.config = config
	app.gin = gin.Default()
	app.gin.LoadHTMLTemplates("templates/*")
	app.gin.GET("/", app.rootHandler)
	app.gin.GET("/status", app.statusHandler)
	app.gin.POST("/", app.formHandler)
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
	http.Redirect(ctx.Writer, ctx.Req, "/status", 303)
}

func (app *appContext) statusHandler(ctx *gin.Context) {
	obj := gin.H{}
	ctx.HTML(200, "status.tmpl", obj)
}
