package log

import (
	"io/ioutil"
	"log"
	"testing"

	"go.uber.org/zap"         //nolint:depguard
	"go.uber.org/zap/zapcore" //nolint:depguard
)

type nullOutput struct{}

func (n2 nullOutput) Write(p []byte) (n int, err error) {
	return ioutil.Discard.Write(p)
}

func (n2 nullOutput) Sync() error {
	return nil
}

func TestLoggerFunctionality(t *testing.T) {
	// this test only use to verify the output manually
	output := nullOutput{}
	//nolint:gocritic
	// output := os.Stderr // uncomment for debugging

	oLog := DefaultLogger
	defer func() {
		DefaultLogger = oLog
	}()
	// standard logger
	l := log.New(output, "", log.LstdFlags|log.Lshortfile)
	l.Printf("stdlib: This is a logger test foo:%s int:%d", "bar", 12)

	// zap with abstraction
	zapDefaultConfig := Configuration{
		LogLevel: LevelInfo,
		Format:   FormatJson,
		Output:   output,
	}
	zapWithAbstraction, err := NewLogger(zapDefaultConfig)
	if err != nil {
		t.Fatalf("failed to initiate logger: %v", err)
	}
	zapWithAbstraction.WithFields(Fields{"foo": "bar", "int": 12}).Infof("zapWithAbstraction: This is a logger test")

	// zap sugar no abstraction
	writer := zapcore.Lock(output)
	core := zapcore.NewCore(getEncoder(FormatJson), writer, zap.InfoLevel)
	zapSugar := zap.New(core,
		zap.AddCaller(),
	).Sugar()
	zapSugar.Infow("zapSugar: This is a logger test",
		"foo", "bar",
		"int", 12,
	)

	// zap normal
	writerNormal := zapcore.Lock(output)
	coreNormal := zapcore.NewCore(getEncoder(FormatJson), writerNormal, zap.InfoLevel)
	zapNormal := zap.New(coreNormal,
		zap.AddCaller(),
	)
	zapNormal.Info("zapNormal: This is a logger test",
		zap.String("foo", "bar"),
		zap.Int("int", 12),
	)

	// zap without json encoder
	// zap with abstraction
	zapPlainConfig := Configuration{
		LogLevel: LevelInfo,
		Format:   FormatPlain,
		Output:   output,
	}
	zapPlainWithAbstraction, err := NewLogger(zapPlainConfig)
	if err != nil {
		t.Fatalf("failed to initiate logger: %v", err)
	}
	zapPlainWithAbstraction.Infof("zapPlainWithAbstraction: logging with plain")

	// test default functionality
	DefaultLogger = zapWithAbstraction
	Infof("this is info message")
	Debugf("this is debug message")
	Warnf("this is warn message")
	Errorf("this is error message")
	Criticalf("this is critical message")
	WithFields(Fields{"key": "value"}).Infof("this is info message with fields")
}

func BenchmarkStdlib(b *testing.B) {
	b.ReportAllocs()
	l := log.New(nullOutput{}, "", log.LstdFlags|log.Lshortfile)
	for i := 0; i < b.N; i++ {
		l.Printf("This is a logger test foo:%s int:%d", "bar", 12)
	}
}

func BenchmarkZapWithAbstraction(b *testing.B) {
	b.ReportAllocs()
	defaultConfig := Configuration{
		LogLevel: LevelInfo,
		Format:   FormatJson,
		Output:   nullOutput{},
	}

	l, err := NewLogger(defaultConfig)
	if err != nil {
		b.Fatalf("failed to initiate logger: %v", err)
	}
	for i := 0; i < b.N; i++ {
		l.WithFields(Fields{"foo": "bar", "int": 12}).Infof("This is a logger test")
	}
}

func BenchmarkZapSugarNoAbstraction(b *testing.B) {
	b.ReportAllocs()
	writer := zapcore.Lock(nullOutput{})
	core := zapcore.NewCore(getEncoder(FormatJson), writer, zap.InfoLevel)

	logger := zap.New(core,
		zap.AddCaller(),
	).Sugar()

	for i := 0; i < b.N; i++ {
		logger.Infow("This is a logger test",
			"foo", "bar",
			"int", 12,
		)
	}
}

func BenchmarkZapNormalNoAbstraction(b *testing.B) {
	b.ReportAllocs()
	writer := zapcore.Lock(nullOutput{})
	core := zapcore.NewCore(getEncoder(FormatJson), writer, zap.InfoLevel)

	logger := zap.New(core,
		zap.AddCaller(),
	)

	for i := 0; i < b.N; i++ {
		logger.Info("This is a logger test",
			zap.String("foo", "bar"),
			zap.Int("int", 12),
		)
	}
}
