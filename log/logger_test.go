package log

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/BurntSushi/toml"

	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	os.Mkdir("bin", os.ModePerm)

	c := &Config{
		Level:            DebugLevel,
		EnableCaller:     true,
		Development:      true,
		EnableStdOutput:  true,
		EnableFileOutput: true,
		FileOutput: FileOutput{
			FileName:   "bin/test.log",
			MaxSize:    1,
			MaxAge:     7,
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   true,
		},
	}

	bytes, err := yaml.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("bin/config-test.yaml", bytes, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal(bytes, c); err != nil {
		t.Fatal(err)
	}
	t.Log(c)

	file, err := os.Create("bin/config-test.toml")
	if err := toml.NewEncoder(file).Encode(c); err != nil {
		t.Fatal(err)
	}
	if _, err := toml.DecodeFile("bin/config-test.toml", c); err != nil {
		t.Fatal(err)
	}
	t.Log(c)
}

func TestStdLogger(t *testing.T) {
	logger, err := CreateLogger(&Config{
		Level:           DebugLevel,
		EnableCaller:    true,
		Development:     true,
		EnableStdOutput: true,
	})
	if err != nil {
		t.Fatal(err)
	}

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

	logger, err := CreateLogger(&Config{
		Level:            DebugLevel,
		EnableCaller:     true,
		Development:      true,
		EnableStdOutput:  false,
		EnableFileOutput: true,
		FileOutput: FileOutput{
			FileName:   "bin/test.log",
			MaxSize:    1,
			MaxAge:     7,
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

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

	logger, err := CreateLogger(&Config{
		Level:            DebugLevel,
		EnableCaller:     true,
		Development:      true,
		EnableStdOutput:  true,
		EnableFileOutput: true,
		FileOutput: FileOutput{
			FileName:   "bin/test.log",
			MaxSize:    1,
			MaxAge:     7,
			MaxBackups: 0,
			LocalTime:  true,
			Compress:   false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

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
	logger, err := CreateLogger(&Config{
		Level:           DebugLevel,
		EnableCaller:    true,
		Development:     true,
		EnableStdOutput: true,
	})
	if err != nil {
		t.Fatal(err)
	}

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
	logger, err := CreateLogger(&Config{
		Level:           DebugLevel,
		EnableCaller:    true,
		Development:     true,
		EnableStdOutput: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	logger = logger.WithFieldsLazy(zap.String("lazy1", "lazy1")).WithFieldsLazy(zap.String("lazy2", "lazy2"))
	logger.Info("lazy")
	logger.InfoFields("lazy", zap.String("test", "test"))
}
