package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	//if len(os.Args) != 2 {
	//	log.Printf("need port number\n")
	//	os.Exit(1)
	//}
	//p := os.Args[0]
	p := "80"
	l, err := net.Listen("tcp", ":"+p)
	if err != nil {
		log.Fatalf("failed to listen port: %s %v", p, err)
	}
	if err := run(context.Background(), l); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		// 인수로 받은 net.Listener를 이용하므로 Addr 필드는 지정하지 않는다.
		//Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		// ListenAndServe 메서드가 아닌 Serve 메서드로 변경한다.
		if err := s.Serve(l); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown %+v", err)
	}
	return eg.Wait()
}
