package main

import (
	"club-server/proto/go/pb"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/akrylysov/algnhsa"

	"github.com/twitchtv/twirp"
)

func main() {
	hooks := &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			return ctx, nil
		},
		ResponseSent: func(ctx context.Context) {

		},
		Error: func(ctx context.Context, err twirp.Error) context.Context {
			return ctx
		},
	}

	http.Handle(pb.HelloPathPrefix, pb.NewHelloServer(&hello{}, hooks))

	http.HandleFunc("/health_check", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	isLocal, _ := strconv.ParseBool(os.Getenv("LOCAL"))
	if isLocal {
		port := os.Getenv("PORT")
		if len(port) == 0 {
			port = "3000"
		}

		log.Printf("running server port:%s", port)
		_ = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	} else {
		algnhsa.ListenAndServe(http.DefaultServeMux, nil)
	}
}

type hello struct {
}

func (h *hello) World(ctx context.Context, req *pb.Empty) (*pb.Text, error) {
	return &pb.Text{
		Text: "Hello World",
	}, nil
}
