package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/illuminatingKong/kongming-kit/runner"
	"net/http"
)

type RouterInjector interface {
	Inject(router *gin.RouterGroup)
}

type HttpServiceEngine struct {
	ProjectName string
	Addr        string
	Middlewares []gin.HandlerFunc
	Mode        string
	RouterGroup map[string]RouterInjector
	VPath       string
	*gin.Engine
}

func (hse *HttpServiceEngine) InjectApiRouter(root *gin.RouterGroup) {
	for name, r := range hse.RouterGroup {
		r.Inject(root.Group(name))
	}
}

func (hse *HttpServiceEngine) Init() *HttpServiceEngine {
	gin.SetMode(hse.Mode)
	g := gin.New()
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Invalid path: %s", c.Request.URL.Path)
	})
	g.Use(hse.Middlewares...)
	g.HandleMethodNotAllowed = true
	g.NoMethod(func(c *gin.Context) {
		c.String(http.StatusMethodNotAllowed, "Method not allowed: %s %s", c.Request.Method, c.Request.URL.Path)
	})
	hse.Engine = g
	apiRouters := hse.Group(hse.VPath)
	hse.InjectApiRouter(apiRouters)
	return hse
}

func (hse *HttpServiceEngine) Stop(ctx context.Context, server *http.Server) error {
	log := runner.GetLogger()
	log.Infof("HttpService engine ready stop")
	var err error
	if e := server.Shutdown(ctx); e != nil {
		err = e
	}
	return err
}

func (hse *HttpServiceEngine) Start(ctx context.Context, server *http.Server) error {
	log := runner.GetLogger()
	log.Infof("HttpService engine ready start at  %s", server.Addr)
	var err error
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return err
}
