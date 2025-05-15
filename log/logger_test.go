package log

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestStdLogger(t *testing.T) {
	logger := NewLogger(&Config{
		Level:        DebugLevel,
		EnableCaller: true,
		Development:  true,
		Cores: []CoreConfig{
			NewStdCoreConfig(),
		},
	})

	logger.Debug(DebugLevel)
	logger.Info(InfoLevel)
	logger.Warn(WarnLevel)
	logger.Error(ErrorLevel)
	//logger.DPanic(DPanicLevel)
	//logger.Panic(PanicLevel)
	//logger.Fatal(FatalLevel)
}

func TestFileLogger(t *testing.T) {
	os.Mkdir("bin", os.ModePerm)

	logger := NewLogger(&Config{
		Level:        DebugLevel,
		EnableCaller: true,
		Development:  true,
		Cores: []CoreConfig{
			NewFileCoreConfig(&FileCoreParams{
				Path:       "bin/test.log",
				MaxSize:    1,
				MaxAge:     7,
				MaxBackups: 0,
				LocalTime:  true,
				Compress:   false,
			}),
		},
	})

	for i := 0; i < 1e4; i++ {
		logger.Debug(DebugLevel)
		logger.Info(InfoLevel)
		logger.Warn(WarnLevel)
		logger.Error(ErrorLevel)
		//logger.DPanic(DPanicLevel)
		//logger.Panic(PanicLevel)
		//logger.Fatal(FatalLevel)
	}
}

func TestStdFileLogger(t *testing.T) {
	os.Mkdir("bin", os.ModePerm)

	logger := NewLogger(&Config{
		Level:        DebugLevel,
		EnableCaller: true,
		Development:  true,
		Cores: []CoreConfig{
			NewStdCoreConfig(),
			NewFileCoreConfig(&FileCoreParams{
				Path:       "bin/test.log",
				MaxSize:    1,
				MaxAge:     7,
				MaxBackups: 0,
				LocalTime:  true,
				Compress:   false,
			}),
		},
	})

	for i := 0; i < 1e4; i++ {
		logger.Debug(DebugLevel)
		logger.Info(InfoLevel)
		logger.Warn(WarnLevel)
		logger.Error(ErrorLevel)
		//logger.DPanic(DPanicLevel)
		//logger.Panic(PanicLevel)
		//logger.Fatal(FatalLevel)
	}
}

func TestLoggerMisc(t *testing.T) {
	logger := NewLogger(&Config{
		Level:        DebugLevel,
		EnableCaller: true,
		Development:  true,
		Cores: []CoreConfig{
			NewStdCoreConfig(),
		},
	})

	logger.Warn("misc")

	loggerWith := logger.WithFields(zap.String("s1", "s1"), zap.String("s2", "s2"))
	loggerWith.Warn("misc")

	loggerNamed := logger.Named("name")
	loggerNamed.Warn("misc")

	loggerNamespace := logger.WithFields(zap.Namespace("misc"))
	loggerNamespace.InfoFields("misc", zap.Int64("int64", 1))
	loggerNamespace.InfoFields("misc")
}

func TestLogger_WithFieldsLazy(t *testing.T) {
	logger := NewLogger(&Config{
		Level:        DebugLevel,
		EnableCaller: true,
		Development:  true,
		Cores: []CoreConfig{
			NewStdCoreConfig(),
		},
	})

	logger = logger.WithFieldsLazy(zap.String("lazy1", "lazy1")).WithFieldsLazy(zap.String("lazy2", "lazy2"))
	logger.Info("lazy")
	logger.InfoFields("lazy", zap.String("test", "test"))
}
