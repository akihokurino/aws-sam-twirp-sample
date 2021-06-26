package main

import (
	"bytes"
	"club-server/proto/go/pb"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/akrylysov/algnhsa"

	"github.com/twitchtv/twirp"

	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func main() {
	env, err := readSSM(os.Getenv("SSM_PATH"))
	if err != nil {
		panic(err)
	}
	envMap, err := godotenv.Parse(bytes.NewBufferString(env))
	if err != nil {
		panic(err)
	}
	for k, v := range envMap {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}

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

func readSSM(path string) (string, error) {
	config := aws.NewConfig()
	sess, err := session.NewSession(config)
	if err != nil {
		return "", err
	}

	svc := ssm.New(sess, &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	res, _ := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: aws.Bool(true),
	})

	val := *res.Parameter.Value

	return val, nil
}

type hello struct {
}

func (h *hello) World(ctx context.Context, req *pb.Empty) (*pb.Text, error) {
	return &pb.Text{
		Text: "Hello World 2",
	}, nil
}
