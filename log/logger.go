package log

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
	"unicode"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

const (
	ACCESS = "ACCESS"
	APPLOG = "APPLOG"
)

var logHandle *LoggerHandle
var loginstance *log.Logger

type Type string
type t2pLogInfo struct {
	//can change after init
	servicename string
	traceID     string
	spanID      string
	token       string
	apiKey      string
	pathURI     string
	method      string
	// message      string
	// info         string
	// logType      string
	environment  string
	responseTime int
	// caller    string
	// indexname string

	// 	"level":"LOG_LEVEL" -with function call
	// 	"LogDateTime":""    -with function call
}

type LoggerHandle struct {
	log     *log.Logger
	logInfo *t2pLogInfo
}

func New(serviceName string, env string, logPathPrefix ...string) *LoggerHandle {

	if logHandle != nil {
		return logHandle
	}

	prefix := ""
	if len(logPathPrefix) > 0 {
		prefix = removeTrailSlash(logPathPrefix[0], "")
	}

	info := &t2pLogInfo{
		servicename: serviceName,
	}

	info.environment = env
	log := newLogger(info, env, prefix)
	logHandle = &LoggerHandle{
		log:     log,
		logInfo: info,
	}

	return logHandle
}

func newLogger(info *t2pLogInfo, env string, prefix string) *log.Logger {
	if loginstance == nil {
		loginstance = log.New()
	}

	ecsConsoleFormatter := ecslogrus.Formatter{}
	// 	FieldMap: ecslogrus.FieldMap{
	// 		ecslogrus.FieldKeyMsg:         "json.msg",
	// 		ecslogrus.FieldKeyLevel:       "json.level",
	// 		ecslogrus.FieldKeyTime:        "LogDateTime",
	// 		ecslogrus.FieldKeyLogrusError: "json.logrus_error",
	// 		ecslogrus.FieldKeyFunc:        "json.func",
	// 		ecslogrus.FieldKeyFile:        "json.caller",
	// 	},
	// 	TimestampFormat: "2006-01-02 15:04:05.000",
	// 	// CallerPrettyfier: caller(),
	// }

	jsonFormatter := log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyLevel:       "log.level",
			log.FieldKeyMsg:         "message",
			log.FieldKeyTime:        "LogDateTime",
			log.FieldKeyLogrusError: "json.logrus_error",
			log.FieldKeyFunc:        "json.func",
			log.FieldKeyFile:        "json.caller",
		},
		TimestampFormat: "2006-01-02 15:04:05.000",
		// CallerPrettyfier: caller(),
	}

	loginstance.SetFormatter(&ecsConsoleFormatter)
	loginstance.SetLevel(log.TraceLevel)

	if os.Getenv("IS_LAMBDA") == "" {
		if exists("/logbeat") {
			path := fmt.Sprintf("/logbeat/%v.log", info.servicename)
			if prefix != "" {
				if !strings.Contains(`DEVELOP`+`PRODUCTION`+`UAT`+`SIT`+`PRE_PRODUCTION`, env) {
					path = fmt.Sprintf("%v/%v.log", prefix, info.servicename)
				}
			}

			writer, _ := rotatelogs.New(
				path+".%Y%m%d",
				rotatelogs.WithLinkName(path),
				rotatelogs.WithMaxAge(24*time.Hour),
				rotatelogs.WithRotationTime(time.Hour),
			)

			loginstance.Hooks.Add(lfshook.NewHook(
				lfshook.WriterMap{
					log.InfoLevel:  writer,
					log.ErrorLevel: writer,
					log.DebugLevel: writer,
					log.WarnLevel:  writer,
					log.TraceLevel: writer,
					log.PanicLevel: writer,
					log.FatalLevel: writer,
				},
				&jsonFormatter,
			))
		}
	}

	// loginstance.SetReportCaller(true)

	return loginstance
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (l *LoggerHandle) indexname(ltype string) string {
	return fmt.Sprintf("%v-%v-%v", strings.ToLower(l.logInfo.servicename), strings.ToLower(l.logInfo.environment), strings.ToLower(ltype))
}

func (l *LoggerHandle) fields(ltype string, info string) log.Fields {
	return log.Fields{
		"indexname":         l.indexname(ltype),
		"flb_key":           "applog",
		"json.apikey":       l.logInfo.apiKey,
		"json.caller":       fileInfo(3),
		"json.environment":  l.logInfo.environment,
		"json.info":         info,
		"json.method":       l.logInfo.method,
		"json.pathuri":      l.logInfo.pathURI,
		"json.responseTime": l.logInfo.responseTime,
		"json.servicename":  l.logInfo.servicename,
		"json.spanid":       l.logInfo.spanID,
		"json.token":        l.logInfo.token,
		"json.traceid":      l.logInfo.traceID,
		"json.type":         ltype,
		"logversion":        "v2",
	}
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func (l *LoggerHandle) Access(msg string, info string) {
	l.log.WithFields(l.fields(ACCESS, info)).Info(msg)
}

func (l *LoggerHandle) Warn(msg string, info string) {
	l.log.WithFields(l.fields(APPLOG, info)).Warn(msg)
}

func (l *LoggerHandle) Debug(msg string, info string) {
	l.log.WithFields(l.fields(APPLOG, info)).Debug(msg)
}

func (l *LoggerHandle) Error(msg string, info string) {
	l.log.WithFields(l.fields(APPLOG, info)).Error(msg)
}

func (l *LoggerHandle) Info(msg string, info string) {
	l.log.WithFields(l.fields(APPLOG, info)).Info(msg)
}

func (l *LoggerHandle) Trace(msg string, info string) {
	l.log.WithFields(l.fields(APPLOG, info)).Trace(msg)
}

func (l *LoggerHandle) SetResponseTime(timeuse int) *LoggerHandle {
	l.logInfo.responseTime = timeuse
	return l
}

func (l *LoggerHandle) SetPathURI(path string) *LoggerHandle {
	l.logInfo.pathURI = path
	return l
}

func (l *LoggerHandle) SetMethod(method string) *LoggerHandle {
	l.logInfo.method = method
	return l
}

func (l *LoggerHandle) SetAPIKey(key string) *LoggerHandle {
	l.logInfo.apiKey = key
	return l
}

func (l *LoggerHandle) SetToken(token string) *LoggerHandle {
	l.logInfo.token = token
	return l
}

func (l *LoggerHandle) SetSpanID(span string) *LoggerHandle {
	l.logInfo.spanID = span
	return l
}

func (l *LoggerHandle) SetTraceID(trace string) *LoggerHandle {
	l.logInfo.traceID = trace
	return l
}

func removeTrailSlash(path string, withPrefix string) string {

	path = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, path)

	path = strings.Trim(path, " ")

	if strings.Trim(withPrefix, " ") != "" {
		for strings.HasSuffix(path, withPrefix) {
			path = strings.TrimSuffix(path, withPrefix)
		}
	} else {
		for strings.HasSuffix(path, string(os.PathSeparator)) {
			path = strings.TrimSuffix(path, string(os.PathSeparator))
		}
	}

	if strings.HasSuffix(path, "\\") {
		return removeTrailSlash(path, "\\")
	}

	if strings.HasSuffix(path, "/") {
		return removeTrailSlash(path, "/")
	}

	path = strings.Trim(path, " ")

	rex, err := regexp.MatchString("[ ]", path)

	if err != nil || rex {
		panic("error Path incorrect: " + path)
	}
	return path
}
