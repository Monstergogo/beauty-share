package logger

import (
	"context"
	"github.com/Monstergogo/beauty-share/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

type LogConf struct {
	LogFilepath string `json:"log_filepath"`
	ErrFilepath string `json:"err_filepath"`
}

func GetLogger() *zap.Logger {
	return logger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filepath string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    10, // 文件最大大小（MB）
		MaxBackups: 5,  // 保留旧文件最大个数
		MaxAge:     30, // 保留旧文件最大天数
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func InitLogger(conf LogConf) {
	encoder := getEncoder()
	// 记录全量日志
	coreLog := zapcore.NewCore(encoder, getLogWriter(conf.LogFilepath), zapcore.DebugLevel)
	var core, coreErr zapcore.Core
	if len(conf.ErrFilepath) > 0 {
		// err日志除了保存在conf.LogFilepath, 还会单独保存在一个文件
		coreErr = zapcore.NewCore(encoder, getLogWriter(conf.ErrFilepath), zapcore.ErrorLevel)
		core = zapcore.NewTee(coreLog, coreErr)
	} else {
		core = coreLog
	}
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func LogWithTraceId(ctx context.Context, level zapcore.Level, msg string, fields ...zapcore.Field) {
	traceId := util.GetTraceIdFromCtx(ctx)
	fieldsNew := make([]zap.Field, 0, len(fields)+1)
	fieldsNew = append(fieldsNew, zap.String("trace_id", traceId))
	fieldsNew = append(fieldsNew, fields...)
	switch level {
	case zapcore.WarnLevel:
		logger.Warn(msg, fieldsNew...)
	case zapcore.ErrorLevel:
		logger.Error(msg, fieldsNew...)
	case zap.DebugLevel:
		logger.Debug(msg, fieldsNew...)
	default:
		logger.Info(msg, fieldsNew...)
	}
}
