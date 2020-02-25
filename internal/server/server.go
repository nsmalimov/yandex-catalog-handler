package server

import (
	"encoding/json"
	"fmt"
	"log"

	"yandex-catalog-handler/internal/concator"
	"yandex-catalog-handler/internal/consumer"
	"yandex-catalog-handler/internal/entity"
	"yandex-catalog-handler/internal/loader"
	"yandex-catalog-handler/internal/result"
	"yandex-catalog-handler/pkg/config"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
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
	results, err := h.resultService.GetAll()

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte(err.Error()))

		return
	}

	b, err := json.Marshal(results)

	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte(err.Error()))

		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(b)
}

func (h *ServerHandler) StartCalc(ctx *fasthttp.RequestCtx) {
	go func() {
		resultMain := entity.Result{
			Cause: "Ok",
		}

		err := h.loaderService.Load()
		if err != nil {
			errS := fmt.Sprintf("Error when try loaderService.Load[h.resultService.Create[(h *ServerHandler) StartCalc], err: %s", err)
			resultMain.Cause = errS
			log.Printf("%s\n", errS)

			err = h.resultService.Create(resultMain)

			if err != nil {
				log.Printf("Error when try h.resultService.Create[h.resultService.Create[(h *ServerHandler) StartCalc], err: %s", err)
			}
			return
		}

		err = h.concatorService.Concate()
		if err != nil {
			errS := fmt.Sprintf("Error when try concatorService.Concate[h.resultService.Create[(h *ServerHandler) StartCalc], err: %s", err)
			resultMain.Cause = errS
			log.Printf("%s\n", errS)

			err = h.resultService.Create(resultMain)

			if err != nil {
				log.Printf("Error when try h.resultService.Create[(h *ServerHandler) StartCalc], err: %s", err)
			}
			return
		}

		b, err := json.Marshal(h.concatorService.Result)
		if err != nil {
			log.Printf("Error when try json.Marshal[(h *ServerHandler) StartCalc], err: %s", err)
		}

		resultMain.Results = b

		err = h.resultService.Create(resultMain)

		if err != nil {
			log.Printf("Error when try h.resultService.Create[(h *ServerHandler) StartCalc], err: %s", err)
		}

		log.Println("Done")

		return
	}()

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte("Started. Wait."))
}

func (s *Server) Run() (err error) {
	router := fasthttprouter.New()
	router.GET("/ping", s.serverHandler.Ping)
	router.POST("/set_delta_time", s.serverHandler.SetDeltaTine)
	router.GET("/get_operate_logs", s.serverHandler.GetOperateLogs)

	router.GET("/start_calc", s.serverHandler.StartCalc)

	router.GET("/", fasthttp.FSHandler(fmt.Sprintf("%s/%s", s.cfg.WebFolderPath, "index.html"), 1))

	router.GET("/static/*filepath", fasthttp.FSHandler(s.cfg.WebFolderPath, 1))

	port := fmt.Sprintf(":%d", s.cfg.Port)

	log.Printf("Ready to start on port: %s\n", port)

	err = fasthttp.ListenAndServe(port, router.Handler)

	return
}
