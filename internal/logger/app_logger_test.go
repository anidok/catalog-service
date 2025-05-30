package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type LoggerTestSuite struct {
	suite.Suite
	buf    *bytes.Buffer
	logger *log.Logger
}

func TestLoggerSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (s *LoggerTestSuite) SetupTest() {
	s.buf = &bytes.Buffer{}
	logger = &logrus.Logger{
		Out:       s.buf,
		Formatter: &logrus.JSONFormatter{},
		Level:     logrus.DebugLevel,
	}
	NonContext = &ContextLogger{
		entry: logrus.NewEntry(logger),
	}
	s.logger = log.Default()
	log.SetOutput(s.buf)
}

func (s *LoggerTestSuite) TearDownTest() {
	log.SetOutput(os.Stderr)
	log.SetFlags(s.logger.Flags())
}

func (s *LoggerTestSuite) TestContextLogger() {
	ctx := context.Background()
	l := NewContextLogger(ctx, "test_method")
	s.NotNil(l)

	// Test different log levels
	testCases := []struct {
		name     string
		logFunc  func()
		contains string
	}{
		{
			name: "Info logging",
			logFunc: func() {
				l.Info("test info")
			},
			contains: "test info",
		},
		{
			name: "Error logging",
			logFunc: func() {
				l.Error("test error", errors.New("error occurred"))
			},
			contains: "error occurred",
		},
		{
			name: "Debug logging",
			logFunc: func() {
				l.Debug("test debug")
			},
			contains: "test debug",
		},
		{
			name: "Warn logging",
			logFunc: func() {
				l.Warn("test warn")
			},
			contains: "test warn",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.buf.Reset()
			tc.logFunc()
			s.Contains(s.buf.String(), tc.contains)
		})
	}
}

func (s *LoggerTestSuite) TestLoggerWithFields() {
	ctx := context.Background()
	l := NewContextLogger(ctx, "test_method")
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 2,
	}

	testCases := []struct {
		name     string
		logFunc  func()
		validate func(map[string]interface{})
	}{
		{
			name: "InfoWithFields",
			logFunc: func() {
				l.InfoWithFields("test info with fields", fields)
			},
			validate: func(data map[string]interface{}) {
				s.Equal("value1", data["key1"])
				s.Equal(float64(2), data["key2"])
			},
		},
		{
			name: "ErrorWithFields",
			logFunc: func() {
				l.ErrorWithFields("test error with fields", fields, errors.New("test error"))
			},
			validate: func(data map[string]interface{}) {
				s.Equal("test error", data[LogError])
				s.Equal("value1", data["key1"])
				s.Equal(float64(2), data["key2"])
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.buf.Reset()
			tc.logFunc()
			var data map[string]interface{}
			err := json.Unmarshal(s.buf.Bytes(), &data)
			s.NoError(err)
			tc.validate(data)
		})
	}
}

func (s *LoggerTestSuite) TestFormatMessages() {
	ctx := context.Background()
	l := NewContextLogger(ctx, "test_method")

	testCases := []struct {
		name     string
		logFunc  func()
		contains string
	}{
		{
			name: "Infof",
			logFunc: func() {
				l.Infof("test %s", "formatted")
			},
			contains: "test formatted",
		},
		{
			name: "Errorf",
			logFunc: func() {
				l.Errorf(errors.New("base error"), "test %s", "formatted")
			},
			contains: "test formatted",
		},
		{
			name: "Debugf",
			logFunc: func() {
				l.Debugf("test %s", "formatted")
			},
			contains: "test formatted",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.buf.Reset()
			tc.logFunc()
			s.Contains(s.buf.String(), tc.contains)
		})
	}
}
func (s *LoggerTestSuite) TestLoggerLevels() {
	tests := []struct {
		level     string
		shouldLog bool
	}{
		{"debug", true},
		{"info", true},
		{"warn", true},
		{"error", true},
		{"fatal", true},
	}

	for _, tt := range tests {
		s.Run(tt.level, func() {
			lvl, _ := logrus.ParseLevel(tt.level)
			logger.SetLevel(lvl)
			s.Equal(lvl, logger.GetLevel())
		})
	}
}

func (s *LoggerTestSuite) TestWarnAndDebugWithFields() {
	ctx := context.Background()
	l := NewContextLogger(ctx, "test_method")
	testErr := errors.New("warn error")
	fields := map[string]interface{}{
		"test_key": "test_value",
		"count":    42,
	}

	testCases := []struct {
		name     string
		logFunc  func()
		validate func(map[string]interface{})
	}{
		{
			name: "WarnWithFields",
			logFunc: func() {
				l.WarnWithFields("test warning", fields, testErr)
			},
			validate: func(data map[string]interface{}) {
				s.Equal("test_value", data["test_key"])
				s.Equal(float64(42), data["count"])
				s.Contains(data["msg"].(string), "test warning")
				s.Equal("warning", data["level"])
				// Print log output for debugging
				fmt.Printf("Log output: %+v\n", data)
			},
		},
		{
			name: "DebugWithFields",
			logFunc: func() {
				l.DebugWithFields("test debug", fields)
			},
			validate: func(data map[string]interface{}) {
				s.Equal("test_value", data["test_key"])
				s.Equal(float64(42), data["count"])
				s.Equal("test debug", data["msg"])
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.buf.Reset()
			tc.logFunc()
			var data map[string]interface{}
			err := json.Unmarshal(s.buf.Bytes(), &data)
			s.NoError(err)
			tc.validate(data)
		})
	}
}

func (s *LoggerTestSuite) TestLoggerFormat() {
	tests := []struct {
		format   string
		validate func(string)
	}{
		{
			format: "json",
			validate: func(output string) {
				var data map[string]interface{}
				err := json.Unmarshal([]byte(output), &data)
				s.NoError(err)
				s.Equal("info", data["level"])
				s.Equal("test message", data["msg"])
			},
		},
		{
			format: "text",
			validate: func(output string) {
				s.Contains(output, "level=info")
				s.Contains(output, "msg=\"test message\"")
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.format, func() {
			// Create new buffer for each test
			testBuf := &bytes.Buffer{}

			// Create and configure logger
			testLogger := &logrus.Logger{
				Out:   testBuf,
				Level: logrus.InfoLevel,
			}

			if tt.format == "json" {
				testLogger.Formatter = &logrus.JSONFormatter{}
			} else {
				testLogger.Formatter = &logrus.TextFormatter{DisableTimestamp: true}
			}

			// Set global logger
			logger = testLogger
			l := NewContextLogger(context.Background(), "test_method")

			// Log test message
			l.Info("test message")

			// Validate output
			output := testBuf.String()
			s.NotEmpty(output)
			tt.validate(output)
		})
	}
}
