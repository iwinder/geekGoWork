package week03

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunNewServer() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	server := http.Server{
		Addr:    "127.0.0.1:8088",
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if err == http.ErrServerClosed {
				glog.Warning("Server closed under request")
			} else {
				glog.Fatalln(fmt.Sprintf("Server  unexpected，Could not gracefully use the server: %v \n", err))
			}
		}
	}()

	graceTimeOut := 30

	osCall := <-stopChan
	glog.V(3).Infoln(fmt.Sprintf("system call: %+v", osCall))
	timeOut := time.Duration(graceTimeOut) * time.Second
	ctx, cancle := context.WithTimeout(context.Background(), timeOut)
	defer cancle()
	if err := server.Shutdown(ctx); err != nil {
		glog.Fatalln(fmt.Sprintf("Could not gracefully shutdown the server: %v  \n", err))
	}

}

func healthzHandler(respon http.ResponseWriter, r *http.Request) {
	io.WriteString(respon, "OK\n")
}
