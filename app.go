package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
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
		base := `<option>` + env + `</option>` + "\n"
		envDropdowns += base
	}
	for _, node := range app.config.Nodes {
		base := `<option>` + node.Hostname + `</option>` + "\n"
		hostDropdowns += base
	}

	obj := gin.H{
		"envDropdowns":  template.HTML(envDropdowns),
		"hostDropdowns": template.HTML(hostDropdowns),
	}
	ctx.HTML(200, "index.tmpl", obj)
}
