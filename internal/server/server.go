package server

import (
	"fmt"
	"log"
	"net/http"

	"yandex-catalog-handler/internal/concator"
	"yandex-catalog-handler/internal/consumer"
	"yandex-catalog-handler/internal/entity"
	"yandex-catalog-handler/internal/loader"
	"yandex-catalog-handler/internal/result"
	"yandex-catalog-handler/pkg/config"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	_ "net/http/pprof"
)

type ServerHandler struct {
	loaderService   *loader.Loader
	concatorService *concator.Concator
	resultService   *result.Service
}

type Server struct {
	serverHandler   ServerHandler
	cfg             config.Config
	consumerService *consumer.Consumer
}

func NewServer(cfg config.Config,
	consumerService *consumer.Consumer,
	loaderService *loader.Loader,
	concatorService *concator.Concator,
	resultService *result.Service) *Server {
	return &Server{
		serverHandler: ServerHandler{
			loaderService:   loaderService,
			concatorService: concatorService,
			resultService:   resultService,
		},
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

func (h *ServerHandler) StartCalc(ctx *fasthttp.RequestCtx) {
	// go func() {
	resultMain := entity.Result{
		Cause: "Ok",
	}

	err := h.loaderService.Load()
	if err != nil {
		errS := fmt.Sprintf("Error when try loaderService.Load, err: %s", err)
		resultMain.Cause = errS
		log.Printf("%s\n", errS)

		err = h.resultService.Create(resultMain)

		if err != nil {
			log.Printf("Error when try h.resultService.Create, err: %s", err)
		}

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("Started. Wait."))
	}

	resultConcator, err := h.concatorService.Concate()
	if err != nil {
		errS := fmt.Sprintf("Error when try concatorService.Concate, err: %s", err)
		resultMain.Cause = errS
		log.Printf("%s\n", errS)

		err = h.resultService.Create(resultMain)

		if err != nil {
			log.Printf("Error when try h.resultService.Create, err: %s", err)
		}

		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("Started. Wait."))
	}

	resultMain.Results = resultConcator

	err = h.resultService.Create(resultMain)

	if err != nil {
		log.Printf("Error when try h.resultService.Create, err: %s", err)
	}

	log.Println("Done")

	// return
	// }()

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte("Started. Wait."))
}

func (s *Server) Run() (err error) {
	router := fasthttprouter.New()
	router.GET("/ping", s.serverHandler.Ping)
	router.POST("/set_delta_time", s.serverHandler.SetDeltaTine)
	router.POST("/get_operate_logs", s.serverHandler.SetDeltaTine)

	router.GET("/start_calc", s.serverHandler.StartCalc)

	router.GET("/index", fasthttp.FSHandler(s.cfg.WebFolderPath, 1))

	router.GET("/get_price", fasthttp.FSHandler("/Users/nurislam_alimov/IdeaProjects/yandex-catalog-handler/data/66343037-3430-3935-2D35-3163632D3131&FranchiseeId=383450", 2))

	router.GET("/static/js/*filepath", fasthttp.FSHandler("/Users/nurislam_alimov/IdeaProjects/yandex-catalog-handler/data/66343037-3430-3935-2D35-3163632D3131&FranchiseeId=383450", 0))

	port := fmt.Sprintf(":%d", s.cfg.Port)

	log.Printf("Ready to start on port: %s\n", port)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	err = fasthttp.ListenAndServe(port, router.Handler)

	return
}
