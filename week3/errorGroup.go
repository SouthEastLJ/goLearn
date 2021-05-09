package main

import (
	"context"
	"fmt"
	"github.com/go-errors/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main(){
	group,ctext := errgroup.WithContext(context.Background())
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello"))
	})

	server := http.Server{
		Handler: serveMux,
		Addr: ":8080",
	}

	// 单个服务退出
	errorOut := make(chan struct{})
	serveMux.HandleFunc("/error", func(writer http.ResponseWriter, request *http.Request) {
		errorOut <- struct{}{}
	})
	// http server
	group.Go(func() error {
		return server.ListenAndServe()
	})
	// http end
	group.Go(func() error {
		select {
		case <-ctext.Done():
			fmt.Printf("errorGroup exit")
		case <- ctext.Done():
			fmt.Printf("server end")
		}
		timeout,_:=context.WithTimeout(context.Background(),10*time.Second)
		return server.Shutdown(timeout)
	})
	// signal end
	group.Go(func() error {
		quit := make(chan os.Signal,0)
		signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)

		select {
		case <-ctext.Done():
			return ctext.Err()
		case sig := <- quit:
			return errors.Errorf("signal:%v",sig)
		}
	})
	group.Wait()
	
}
