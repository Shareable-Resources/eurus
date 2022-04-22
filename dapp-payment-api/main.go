package main

import (
	"dapp-payment-api/config"
	"dapp-payment-api/implementation"
	"dapp-payment-api/oapi"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./config.json", "Config file path")
	flag.Parse()

	c := config.NewServerConfig()
	err := c.ReadFromFile(configFile)
	if err != nil {
		panic(err)
	}

	l, err := InitLogger(c)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(httplog.NewLogger("", httplog.Options{JSON: true})))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   c.CORS.AllowedOrigins,
		AllowedMethods:   c.CORS.AllowedMethods,
		AllowedHeaders:   c.CORS.AllowedHeaders,
		AllowCredentials: c.CORS.AllowCredentials,
		MaxAge:           c.CORS.MaxAge,
	}))

	r.Mount("/", oapi.Handler(implementation.NewPaymentAPI(c, l)))

	l.Infof("Start listening on :%v", c.ListenPort)
	http.ListenAndServe(fmt.Sprintf(":%v", c.ListenPort), r)
}

func InitLogger(c *config.ServerConfig) (*logrus.Logger, error) {
	var extension string
	var fileName string

	path := c.LogFilePath
	lastIndex := strings.LastIndex(path, ".")
	if lastIndex >= 0 {
		extension = path[lastIndex:]
		fileName = path[:lastIndex]
	}

	writer, err := rotatelogs.New(fileName+"_%Y-%m-%d"+extension, rotatelogs.WithLinkName(path), rotatelogs.WithMaxAge(24*time.Hour), rotatelogs.WithRotationTime(24*time.Hour))
	if err != nil {
		return nil, err
	}

	ret := logrus.New()
	ret.SetOutput(writer)
	ret.SetLevel(logrus.DebugLevel)
	ret.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
		ForceColors:     false,
	})
	ret.SetReportCaller(true)

	return ret, nil
}
