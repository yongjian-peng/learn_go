package logger

import (
	"fmt"
	"path"
	"runtime"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logs sync.Map

type LogConfig struct {
	StorageLocation string `yaml:"storage_location"`
	MaxAge          int    `yaml:"max_age"`
	MaxBackups      int    `yaml:"max_backups"`
	MaxSize         int    `yaml:"max_size"`
	LogLevel        string `yaml:"log_level"`
}

var defaultConfig LogConfig

func SetLogConfig(logConfig LogConfig) {
	defaultConfig = logConfig
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// GetLogger 获取日志
func GetLogger(logName string) *zap.Logger {
	log, ok := logs.Load(logName)
	if ok {
		return log.(*zap.Logger)
	}

	logConfig := defaultConfig
	hook := lumberjack.Logger{
		Filename:   fmt.Sprintf("%s%s.log", logConfig.StorageLocation, logName), // 日志文件路径
		MaxSize:    logConfig.MaxSize,                                           // 最大bytes
		MaxBackups: logConfig.MaxBackups,                                        // 最多保留多少个备份
		MaxAge:     logConfig.MaxAge,                                            // days
		Compress:   true,                                                        // 是否压缩 disabled by default
	}

	w := zapcore.AddSync(&hook)

	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
	var level zapcore.Level
	switch logConfig.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := NewEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)

	//logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger := zap.New(core)
	//存储全局log引用
	logs.Store(logName, logger)
	return logger
}

func ApiInfo(fileName, rquestId string, msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String("requetstId:", rquestId))
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	GetLogger(fileName).Info(msg, fields...)
}

func ApiError(fileName, rquestId string, msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String("requetstId:", rquestId))
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	GetLogger(fileName).Error(msg, fields...)
}

func ApiWarn(fileName, rquestId string, msg string, fields ...zapcore.Field) {
	fields = append(fields, zap.String("requetstId:", rquestId))
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	GetLogger(fileName).Warn(msg, fields...)
}

func getCallerInfoForLog() (callerFields []zap.Field) {

	pc, file, line, ok := runtime.Caller(2) // 回溯两层，拿到写日志的调用方的函数信息
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) //Base函数返回路径的最后一个元素，只保留函数名

	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return
}
