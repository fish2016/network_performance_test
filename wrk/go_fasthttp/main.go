package main

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello from GO!!!\n")
}

func main() {
	fmt.Println("listening on 0.0.0.0:8090")

	handler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/hello":
			hello(ctx)
		default:
			ctx.Error("Not Found", fasthttp.StatusNotFound)
		}
	}

	if err := fasthttp.ListenAndServe(":8090", handler); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
