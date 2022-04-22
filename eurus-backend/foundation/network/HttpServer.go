package network

import (
	"eurus-backend/foundation/log"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type IHttpConfig interface {
	//Mandatory field. Empty string to bind all network interface
	GetServerIp() string
	//set 0 will use default port 80
	GetServerPort() uint

	GetAccessLogFilePath() string
	//If true, GetAccessLog() will be ignored
	IsUseDefaultAccessLog() bool

	IsEnableHttpAccessLog() bool

	GetAccessLog() *logrus.Logger
	//Mandatory field
	GetErrorLog() *logrus.Logger
}

type HttpServer struct {
	config   IHttpConfig
	Router   *chi.Mux
	Logger   *logrus.Logger
	IOWriter io.Writer // used in Logger *logrus.Logger, used when enableAccesslog
}

func NewServer(config IHttpConfig) (*HttpServer, error) {
	server := new(HttpServer)
	server.config = config
	server.Router = chi.NewRouter()
	server.Router.Use(middleware.Timeout(30 * time.Second))

	if config.IsUseDefaultAccessLog() {
		logger, err := createAccessLog(config)
		if err != nil {
			log.GetDefaultLogger().Error("Unable to create default access log file: ", err.Error())
		} else {
			server.Router.Use(middleware.CleanPath)
			server.Logger = logger
			server.IOWriter = logger.Out
			server.Router.Use(NewStructuredLogger(logger))
		}

	} else if config.GetAccessLog() != nil {
		server.Router.Use(NewStructuredLogger(config.GetAccessLog()))
	} else {
		logger, err := createAccessLog(config)
		if err != nil {
			log.GetDefaultLogger().Error("Unable to create default access log file: ", err.Error())
		} else {
			server.Router.Use(middleware.CleanPath)
			server.Router.Use(NewStructuredLogger(logger))
		}
		server.IOWriter = logger.Out
		server.Logger = logger
		if config.IsEnableHttpAccessLog() {
			logger.Out = server.IOWriter
		} else {
			logger.Out = ioutil.Discard
		}
	}

	return server, nil
}

func (me *HttpServer) HttpAuthContext(next http.Handler) http.Handler {

	return http.HandlerFunc(func(responseWriter http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(responseWriter, req)
	})
}

func (me *HttpServer) SetIndexResponse(handler func(writer http.ResponseWriter, req *http.Request)) {
	me.Router.Get("/", handler)
}

func createAccessLog(config IHttpConfig) (*logrus.Logger, error) {
	logFilePath := config.GetAccessLogFilePath()

	if logFilePath == "" {
		logFilePath = "./accessLog.log"
	}
	return log.NewLogger("HTTPAccessLog", logFilePath, logrus.DebugLevel)
}

func (me *HttpServer) Listen() error {
	fmt.Println("Server start listening at ", me.config.GetServerIp(), ":", me.config.GetServerPort())
	log.GetLogger(log.Name.Root).Infoln("Server start listening at ", me.config.GetServerIp(), ":", me.config.GetServerPort())

	var port uint = me.config.GetServerPort()
	if port == 0 {
		port = 80
	}
	addr := me.config.GetServerIp() + ":" + strconv.FormatUint(uint64(port), 10)
	err := http.ListenAndServe(addr, me.Router)
	if err != nil {
		me.config.GetErrorLog().Error("Unable to listen on HTTP port: ", me.config.GetServerPort())
	}
	return err
}

// StructuredLogger is a simple, but powerful implementation of a custom structured
// logger backed on logrus. I encourage users to copy it, adapt it and make it their
// own. Also take a look at https://github.com/pressly/lg for a dedicated pkg based
// on this work, designed for context-based http routers.

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLogger struct {
	Logger *logrus.Logger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()
	contentLen := r.ContentLength
	logFields["content_length"] = strconv.FormatInt(contentLen, 10)
	contentType := r.Header.Get("Content-Type")
	logFields["content_type"] = contentType

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Infoln("request started")

	return entry
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Infoln("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// Helper methods used by the application to get the request-scoped
// logger entry and set additional fields between handlers.
//
// This is a useful pattern to use to set state on the entry as it
// passes through the handler chain, which at any point can be logged
// with a call to .Print(), .Info(), etc.

func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}

func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}

func (me *HttpServer) EnableAccessLog() {
	me.Logger.Out = me.IOWriter
}

func (me *HttpServer) DisableAccessLog() {
	me.Logger.Out = ioutil.Discard
}
