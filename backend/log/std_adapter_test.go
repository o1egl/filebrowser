package log_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/filebrowser/filebrowser/v3/log"
	mock "github.com/filebrowser/filebrowser/v3/log/mock"
)

func TestToStdLogger(t *testing.T) {
	const (
		msg       = "hello world"
		expectMsg = "hello world\n"
	)
	testCases := map[string]struct {
		mockInit func(t *testing.T, loggerMock *mock.MockLogger)
		level    log.Level
	}{
		"debug": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Debugf(expectMsg)
			},
			level: log.LevelDebug,
		},
		"info": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Infof(expectMsg)
			},
			level: log.LevelInfo,
		},
		"warning": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Warnf(expectMsg)
			},
			level: log.LevelWarn,
		},
		"error": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Errorf(expectMsg)
			},
			level: log.LevelError,
		},
		"critical": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Criticalf(expectMsg)
			},
			level: log.LevelCritical,
		},
		"fatal": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Fatalf(expectMsg)
			},
			level: log.LevelFatal,
		},
		"unsupported": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Infof(expectMsg)
			},
			level: log.Level(999),
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl, _ := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			loggerMock := mock.NewMockLogger(ctrl)
			tt.mockInit(t, loggerMock)

			stdLogger := log.ToStdLogger(loggerMock, tt.level)
			stdLogger.Print(msg)
		})
	}
}
