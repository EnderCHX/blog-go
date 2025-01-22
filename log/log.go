package log

import (
	"blog-go/config"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

var (
	cWhiteBlue    = "\033[1;37;44m"
	cWhiteRed     = "\033[1;37;41m"
	cWhiteYellow  = "\033[1;37;43m"
	cWhiteBlack   = "\033[1;37;40m"
	cWhiteGreen   = "\033[1;37;42m"
	cWhiteCyan    = "\033[1;37;46m"
	cWhiteMagenta = "\033[1;37;45m"
	cEnd          = "\033[0m"
)

type Encoder struct {
	zapcore.Encoder
	separator string
	title     string
	color     bool
}

func (e *Encoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf := buffer.NewPool().Get()

	ifColor := func(s, c string) string {
		if !e.color {
			return s
		}
		return c + s + cEnd
	}

	// 标题
	buf.AppendString(ifColor(e.title, cWhiteBlue))
	buf.AppendString(ifColor(e.separator, cWhiteYellow))

	// 时间
	buf.AppendString(ifColor(entry.Time.Format("2006-01-02 15:04:05"), cWhiteGreen))
	buf.AppendString(ifColor(e.separator, cWhiteYellow))

	// 日志级别
	if entry.Level == zapcore.DebugLevel {
		buf.AppendString(ifColor("[DEBUG]", cWhiteBlack))
	} else if entry.Level == zapcore.InfoLevel {
		buf.AppendString(ifColor("[INFO] ", cWhiteBlue))
	} else if entry.Level == zapcore.WarnLevel {
		buf.AppendString(ifColor("[WARN] ", cWhiteYellow))
	} else if entry.Level == zapcore.ErrorLevel {
		buf.AppendString(ifColor("[ERROR]", cWhiteRed))
	}
	buf.AppendString(ifColor(e.separator, cWhiteYellow))

	// 调用者
	if entry.Caller.Defined {
		buf.AppendString(entry.Caller.TrimmedPath())
		buf.AppendString(ifColor(e.separator, cWhiteYellow))
	}

	// 消息
	buf.AppendString(ifColor(entry.Message, cWhiteCyan))

	// 字段
	for _, field := range fields {
		buf.AppendString(ifColor(e.separator, cWhiteYellow))
		buf.AppendString(ifColor(field.Key, cWhiteMagenta))
		buf.AppendString(ifColor("=", cWhiteMagenta))
		if field.Type == 15 {
			buf.AppendString(ifColor(field.String, cWhiteMagenta))
		} else if field.Type == 11 {
			buf.AppendString(ifColor(fmt.Sprintf("%v", field.Integer), cWhiteMagenta))
		} else {
			buf.AppendString(ifColor("nil", cWhiteMagenta))
		}
		buf.AppendString(ifColor(" ", cWhiteMagenta))
	}

	buf.AppendString("\n")
	return buf, nil
}
func init() {
	var level zap.AtomicLevel
	config_ := config.ConfigContext
	if config_.LogCongfig.LogLevel == "debug" {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else if config_.LogCongfig.LogLevel == "info" {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	lumberjackLogger := &lumberjack.Logger{
		Filename:   config_.LogCongfig.LogPath,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
	}

	core := zapcore.NewTee(
		zapcore.NewCore(
			&Encoder{Encoder: encoder, separator: "|", title: "[BLOG]"},
			zapcore.AddSync(lumberjackLogger),
			level,
		),
		zapcore.NewCore(
			&Encoder{Encoder: encoder, separator: "", title: "[BLOG]", color: true},
			zapcore.AddSync(os.Stdout),
			level,
		),
	)

	Logger = zap.New(core)
	defer Logger.Sync()

	Logger.Info("日志初始化成功")
}

func GetLogger() *zap.Logger {
	return Logger
}
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
