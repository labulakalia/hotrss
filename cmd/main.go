package main

import (
	"context"
	"flag"
	"fmt"
	"hotrss/internal/crawler"
	"hotrss/internal/route"
)

func main() {
	baseurl := flag.String("baseurl", "http://127.0.0.1:8080", "your feed base ip or url eg: http://1.1.1.1:8080 http://labulaka521.top")
	port := flag.Int("port", 8080, "http server run port")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = context.WithValue(ctx, "baseurl", *baseurl)
	crawler.RegistryCrawlers(ctx)

	r := route.InitRoute()
	r.Run(fmt.Sprintf("0.0.0.0:%d", *port))

}
