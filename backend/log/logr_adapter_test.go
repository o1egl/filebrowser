package log_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/filebrowser/filebrowser/v3/log"
	mock "github.com/filebrowser/filebrowser/v3/log/mock"
)

func TestLogrAdapter_Logf(t *testing.T) {
	testCases := map[string]struct {
		mockInit func(t *testing.T, loggerMock *mock.MockLogger)
		format   string
		args     []interface{}
	}{
		"debug": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Debugf("hello %s", "world")
			},
			format: "[DEBUG] hello %s",
			args:   []interface{}{"world"},
		},
		"info": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Infof("hello %s", "world")
			},
			format: "[INFO] hello %s",
			args:   []interface{}{"world"},
		},
		"warn": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Warnf("hello %s", "world")
			},
			format: "[WARN] hello %s",
			args:   []interface{}{"world"},
		},
		"warning": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Warnf("hello %s", "world")
			},
			format: "[WARNING] hello %s",
			args:   []interface{}{"world"},
		},
		"err": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Errorf("hello %s", "world")
			},
			format: "[ERR] hello %s",
			args:   []interface{}{"world"},
		},
		"error": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Errorf("hello %s", "world")
			},
			format: "[ERROR] hello %s",
			args:   []interface{}{"world"},
		},
		"critical": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Criticalf("hello %s", "world")
			},
			format: "[CRITICAL] hello %s",
			args:   []interface{}{"world"},
		},
		"fatal": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Fatalf("hello %s", "world")
			},
			format: "[FATAL] hello %s",
			args:   []interface{}{"world"},
		},
		"no level": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Infof("hello %s", "world")
			},
			format: "hello %s",
			args:   []interface{}{"world"},
		},
		"malformed level": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Infof("[ERROR hello %s", "world")
			},
			format: "[ERROR hello %s",
			args:   []interface{}{"world"},
		},
		"unsupported level": {
			mockInit: func(t *testing.T, loggerMock *mock.MockLogger) {
				loggerMock.EXPECT().Infof("[DBG] hello %s", "world")
			},
			format: "[DBG] hello %s",
			args:   []interface{}{"world"},
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl, _ := gomock.WithContext(context.Background(), t)
			defer ctrl.Finish()
			loggerMock := mock.NewMockLogger(ctrl)
			tt.mockInit(t, loggerMock)

			adapter := log.NewLogrAdapter(loggerMock)
			adapter.Logf(tt.format, tt.args...)
		})
	}
}
