package server

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"yandex-catalog-handler/internal/consumer"
	"yandex-catalog-handler/pkg/config"
)

type ServerHandler struct {
}

type Server struct {
	serverHandler   ServerHandler
	cfg             config.Config
	consumerService *consumer.Consumer
}

func NewServer(cfg config.Config, consumerService *consumer.Consumer) *Server {
	return &Server{
		serverHandler:   ServerHandler{},
		cfg:             cfg,
		consumerService: consumerService,
	}
}

func (h *ServerHandler) Ping(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)

	ctx.SetBody([]byte("Pong"))
}

func (h *ServerHandler) SetDeltaTine(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)

	ctx.SetBody([]byte("this is completely new body contents"))
}

func (h *ServerHandler) GetOperateLogs(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)

	ctx.SetBody([]byte("some logs"))
}

func (s *Server) Run() (err error) {
	router := fasthttprouter.New()
	router.GET("/ping", s.serverHandler.Ping)
	router.POST("/set_delta_time", s.serverHandler.SetDeltaTine)
	router.POST("/get_operate_logs", s.serverHandler.SetDeltaTine)

	router.GET("/index", fasthttp.FSHandler(s.cfg.WebFolderPath, 1))

	fmt.Println(fmt.Sprintf("%s/js", s.cfg.WebFolderPath))

	router.GET("/static/js/*filepath", fasthttp.FSHandler(s.cfg.WebFolderPath, 1))

	err = fasthttp.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Port), router.Handler)

	return
}
