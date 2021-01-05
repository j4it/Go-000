package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
Target:
基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal
信号的注册和处理，要保证能够一个退出，全部注销退出
*/
func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleIndex))
	srv := &http.Server{
		Addr:    ":18090",
		Handler: mux,
	}
	g, ctx := errgroup.WithContext(context.Background())
	//	done := make(chan struct{})
	g.Go(func() error {
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				fmt.Println("http server: ", err)
			}
		}()
		select {
		case <-ctx.Done():
			defer func() {
				//				close(done)
				fmt.Println("shutting down server...")
			}()
			time.Sleep(5 * time.Second)
			return srv.Shutdown(ctx)
		}
	})
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

		for {
			s := <-c

			fmt.Println("get a signal:", s.String())
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
				return fmt.Errorf("SIG:%s, start to shutdown srv", s.String())
			default:
				fmt.Println("got SIG:%s", s.String())
			}
		}
	})
	if err := g.Wait(); err != nil {
		fmt.Println("errgroup got err: ", err)
	}

	//	<-done
	fmt.Println("close all")
}

func newSrv(done chan struct{}) {}

func handleSignal(done chan struct{}) {

}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}
