package main

import (
	"context"
	"flag"
	"fmt"
	"hotrss/internal/crawler"
	"hotrss/internal/route"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	baseurl := flag.String("baseurl", "http://127.0.0.1:8080", "your feed base ip or url eg: http://1.1.1.1:8080 http://labulaka521.top")
	port := flag.Int("port", 8080, "http server run port")
	logfile := flag.String("logfile", "/var/log/hotrss.log", "log file")
	flag.Parse()

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	lublogger := lumberjack.Logger{
		Filename: *logfile,
	}
	log.Logger = zerolog.New(&lublogger).With().Caller().Timestamp().Logger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = context.WithValue(ctx, interface{}("baseurl"), *baseurl)
	crawler.RegistryCrawlers(ctx)

	r := route.InitRoute(*logfile)

	r.Run(fmt.Sprintf("0.0.0.0:%d", *port))
}
