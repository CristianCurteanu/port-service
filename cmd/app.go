package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	server   *http.Server
	listener net.Listener
}

func NewApp(listener net.Listener, handler http.Handler) *App {
	return &App{
		listener: listener,
		server: &http.Server{
			Handler: handler,
		},
	}
}

func (a *App) Run() error {
	return a.server.Serve(a.listener)
}

func (a *App) CloseWithContext(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.CloseWithContext(ctx)
}

type Domain struct {
	Path        string
	Middlewares []gin.HandlerFunc
	Routes      []Route
}

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

func NewRouter(groups ...Domain) http.Handler {
	router := gin.Default()

	for _, group := range groups {
		newGroup := router.Group(group.Path)
		for _, middleware := range group.Middlewares {
			newGroup.Use(middleware)
		}

		for _, route := range group.Routes {
			newGroup.Handle(route.Method, route.Path, route.Handler)
		}
	}

	return router
}

func BuildApp(port int, modules ...Domain) (*App, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	handler := NewRouter(modules...)
	app := NewApp(listener, handler)
	return app, nil
}
