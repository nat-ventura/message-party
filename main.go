package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nat/socket-party/routes"
	"github.com/nat/socket-party/routes/users"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = level.NewFilter(logger, level.AllowAll())

	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	handler := routes.Handler{
		Logger: logger,
		User:   users.NewServer(logger),
	}

	http.ListenAndServe(":8080", start(&handler))

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalChan

	level.Info(logger).Log("Server stopping", "party server", "sig", sig)
}

func start(handler *routes.Handler) *httprouter.Router {
	router := httprouter.New()

	return router
}
