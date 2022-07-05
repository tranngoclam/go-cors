package main

import (
	chimdw "github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net"
	"net/http"
	"strings"
)

const addr = ":8081"

var (
	allowedMethods = []string{}
	allowedOrigins = []string{"*"}
)

func main() {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)
	mw := chainMiddleware()

	srv := http.Server{
		Addr: ":8081",
		Handler: chimdw.Recoverer(cors.New(cors.Options{
			AllowedMethods: allowedMethods,
			AllowedOrigins: allowedOrigins,
			Debug:          true,
		}).Handler(healthcheckMiddleware(mw(mux)))),
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(srv.Serve(lis))
}

type middleware func(h http.Handler) http.Handler

func chainMiddleware(mws ...middleware) middleware {
	n := len(mws)
	return func(h http.Handler) http.Handler {
		chainer := func(mw middleware, handler http.Handler) http.Handler {
			return mw(handler)
		}

		chained := h
		for i := n - 1; i >= 0; i-- {
			chained = chainer(mws[i], chained)
		}

		return chained
	}
}

func healthcheckMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.EqualFold(r.URL.Path, "/healthcheck") {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
