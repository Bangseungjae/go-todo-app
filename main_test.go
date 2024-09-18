package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net"
	"net/http"
	"testing"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	log.Printf("url: %s", url)
	rsp, err := http.Get(url)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// HTTP 서버와 반환값을 검증한다.
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// run 함수에 종료 알림을 전송한다.
	cancel()

	// run 함수의 반환값을 검증한다.
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
