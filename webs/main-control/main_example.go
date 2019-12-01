package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/micro/cli"
	"github.com/micro/go-micro/web"
	"net/http"
)

func main() {
	var (
		beforeStartCalled bool
		afterStartCalled  bool
		beforeStopCalled  bool
		afterStopCalled   bool
		str               = `<html><body><h1>Hello World</h1></body></html>`
		fn                = func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, str)
			fmt.Fprint(w, r.Context().Value("id"))
			r.WithContext(context.WithValue(r.Context(), "log", "test_log"))

		}
	)

	//loggingMiddleware := func(next http.Handler) http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		// Do stuff here
	//		log.Println(r.RequestURI)
	//		ctx := context.WithValue(r.Context(), "id", "uuid")
	//		// Call the next handler, which can be another middleware in the chain, or the final handler.
	//		next.ServeHTTP(w, r.WithContext(ctx))
	//		log.Println(r.Context().Value("id"))
	//		log.Println(r.Context().Value("log"))
	//	})
	//}

	beforeStart := func() error {
		beforeStartCalled = true
		return nil
	}

	afterStart := func() error {
		afterStartCalled = true
		return nil
	}

	beforeStop := func() error {
		beforeStopCalled = true
		return nil
	}

	afterStop := func() error {
		afterStopCalled = true
		return nil
	}

	service := web.NewService(
		web.Name("go.micro.web.test"),
		web.BeforeStart(beforeStart),
		web.AfterStart(afterStart),
		web.BeforeStop(beforeStop),
		web.AfterStop(afterStop),
		web.Flags(cli.StringFlag{Name: "tt",}),
	)
	_ = service.Init(web.Action(func(context *cli.Context) {
		fmt.Println(context.String("t"))
	}), )
	router := mux.NewRouter()
	router.HandleFunc("/", fn)
	//router.HandleFunc("/", fn)
	//router.Use(loggingMiddleware)
	service.Handle("/", router)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}

}
