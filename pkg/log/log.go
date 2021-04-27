package log

import (
    "fmt"
    "go.uber.org/zap/zapcore"
    "os"
    "strings"
)
import "go.uber.org/zap"

//https://studygolang.com/articles/17394

var logger *zap.Logger

func init() {
    var err error
    mode := "dev"
    if m := strings.TrimSpace(os.Getenv("mode")); m != "" {
        mode = m
    }

    if mode == "dev" {
        encoderConfig := zapcore.EncoderConfig{
            TimeKey:        "time",
            LevelKey:       "level",
            NameKey:        "logger",
            CallerKey:      "caller",
            MessageKey:     "msg",
            StacktraceKey:  "stacktrace",
            LineEnding:     zapcore.DefaultLineEnding,
            EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
            EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
            EncodeDuration: zapcore.SecondsDurationEncoder,
            EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
        }

        // 设置日志级别
        atom := zap.NewAtomicLevelAt(zap.DebugLevel)

        config := zap.Config{
            Level:            atom,                                                // 日志级别
            Development:      true,                                                // 开发模式，堆栈跟踪
            Encoding:         "console",                                           // 输出格式 console 或 json
            EncoderConfig:    encoderConfig,                                       // 编码器配置
            InitialFields:    map[string]interface{}{"serviceName": "spikeProxy"}, // 初始化字段，如：添加一个服务器名称
            OutputPaths:      []string{"stdout"},                                  // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
            ErrorOutputPaths: []string{"stderr"},
        }

        // 构建日志
        logger, err = config.Build()
        if err != nil {
            panic(fmt.Sprintf("log 初始化失败: %v", err))
        }
        //logger, err = zap.NewDevelopment()
        //if err != nil {
        //    panic(err)
        //}
        logger = zap.New(logger.Core(), zap.AddCaller(), zap.AddCallerSkip(1))
    } else if mode == "prod" {
        logger, err = zap.NewProduction()
    } else {
        panic(fmt.Errorf("unsupported logger mode: %s", mode))
    }
    if err != nil {
        panic(err)
    }

    fmt.Printf("setting up logger mode: %s\n", mode)
}

func Debug(msg string, fields ...zap.Field) {
    logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
    logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
    logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
    logger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
    logger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
    logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
    logger.Fatal(msg, fields...)
}
