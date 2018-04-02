// package main

// import (
// 	_ "fmt"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/nat-ventura/message-party/routes"
// 	"github.com/nat-ventura/message-party/routes/users"

// 	client "github.com/nat-ventura/message-party/client"
// 	server "github.com/nat-ventura/message-party/server"

// 	"github.com/go-kit/kit/log"
// 	"github.com/go-kit/kit/log/level"
// 	"github.com/julienschmidt/httprouter"
// )

// func main() {
// 	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
// 	logger = level.NewFilter(logger, level.AllowAll())

// 	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

// 	handler := routes.Handler{
// 		Logger: logger,
// 		User:   users.NewServer(logger),
// 	}

// 	port := ":3000"

// 	logger.Log("transport", "HTTP", "port", port, "msg", "listening")
// 	err := http.ListenAndServe(port, start(&handler))
// 	if err != nil {
// 		logger.Log("error starting server: ", err)
// 	}

// 	signalChan := make(chan os.Signal)
// 	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
// 	sig := <-signalChan

// 	level.Info(logger).Log("server stopping", "party server", "sig", sig)
// }

// func start(handler *routes.Handler) *httprouter.Router {
// 	router := httprouter.New()

// 	return router
// }
