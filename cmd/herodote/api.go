package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	_ "net/http/pprof"

	"github.com/ViBiOh/herodote/pkg/herodote"
	"github.com/ViBiOh/httputils/v4/pkg/cors"
	"github.com/ViBiOh/httputils/v4/pkg/httputils"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/owasp"
	"github.com/ViBiOh/httputils/v4/pkg/prometheus"
	"github.com/ViBiOh/httputils/v4/pkg/recoverer"
	"github.com/ViBiOh/httputils/v4/pkg/renderer"
	"github.com/ViBiOh/httputils/v4/pkg/server"
)

//go:embed templates static
var content embed.FS

func main() {
	config, err := newConfig()
	if err != nil {
		logger.Fatal(fmt.Errorf("config: %s", err))
	}

	go func() {
		fmt.Println(http.ListenAndServe("localhost:9999", http.DefaultServeMux))
	}()

	ctx := context.Background()

	client, err := newClients(ctx, config)
	if err != nil {
		logger.Fatal(fmt.Errorf("client: %s", err))
	}

	defer client.Close(ctx)

	adapter := newAdapters(client)

	appServer := server.New(config.appServer)
	promServer := server.New(config.promServer)
	prometheusApp := prometheus.New(config.prometheus)

	herodoteApp, err := herodote.New(config.herodote, adapter.adapter, client.tracer.GetTracer("herodote"))
	logger.Fatal(err)

	rendererApp, err := renderer.New(config.renderer, content, herodote.FuncMap, client.tracer.GetTracer("renderer"))
	logger.Fatal(err)

	rendererHandler := rendererApp.Handler(herodoteApp.TemplateFunc)

	go promServer.Start(client.health.ContextEnd(), "prometheus", prometheusApp.Handler())
	go appServer.Start(client.health.ContextEnd(), "http", httputils.Handler(rendererHandler, client.health, recoverer.Middleware, prometheusApp.Middleware, client.tracer.Middleware, owasp.New(config.owasp).Middleware, cors.New(config.cors).Middleware))

	client.health.WaitForTermination(appServer.Done())
	server.GracefulWait(appServer.Done(), promServer.Done())
}
