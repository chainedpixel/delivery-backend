package unit

import (
	config "github.com/MarlonG1/delivery-backend/configs"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	"github.com/sirupsen/logrus"
	"testing"
)

type TestHook struct {
	Entries []*logrus.Entry
}

func NewTestHook() *TestHook {
	return &TestHook{
		Entries: make([]*logrus.Entry, 0),
	}
}

func (hook *TestHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *TestHook) Fire(entry *logrus.Entry) error {
	hook.Entries = append(hook.Entries, entry)
	return nil
}

func TestLogging(t *testing.T) {
	testCases := []struct {
		name    string
		level   logrus.Level
		message string
		details map[string]interface{}
	}{
		{
			name:    "Testing Debug Log Level",
			level:   logrus.DebugLevel,
			message: "This is a debug message",
		},
		{
			name:    "Testing Info Log Level",
			level:   logrus.InfoLevel,
			message: "This is an info message",
		},
		{
			name:    "Testing Warn Log Level",
			level:   logrus.WarnLevel,
			message: "This is a warn message",
		},
		{
			name:    "Testing Error Log Level",
			level:   logrus.ErrorLevel,
			message: "This is an error message",
		},
		{
			name:    "Testing Fatal Log Level",
			level:   logrus.FatalLevel,
			message: "This is a fatal message",
		},
		{
			name:    "Testing Log Level with Details",
			level:   logrus.InfoLevel,
			message: "This is an info message with details",
			details: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:    "Testing Not Existing Log Level",
			level:   logrus.InfoLevel,
			message: "This is when the log level is not existing, then it should be Info",
		},
	}

	for _, tc := range testCases {
		envConfig, err := config.NewEnvConfig()
		if err != nil {
			t.Error(err)
			return
		}

		_ = logs.InitLogger(envConfig)
		testHook := NewTestHook()
		logs.Logger.AddHook(testHook)

		switch tc.level {
		case logrus.DebugLevel:
			logs.Debug(tc.message, tc.details)
		case logrus.InfoLevel:
			logs.Info(tc.message, tc.details)
		case logrus.WarnLevel:
			logs.Warn(tc.message, tc.details)
		case logrus.ErrorLevel:
			logs.Error(tc.message, tc.details)
		case logrus.FatalLevel:
			logs.Fatal(tc.message, tc.details)
		default:
			logs.Info(tc.message, tc.details)
		}

		if len(testHook.Entries) == 0 {
			t.Errorf("Expected at least one log entry, but got none")
		}

		lastEntry := testHook.Entries[len(testHook.Entries)-1]
		if lastEntry.Level != tc.level {
			t.Errorf("Expected log level %s, but got %s", tc.level, lastEntry.Level)
		}
		if lastEntry.Message != tc.message {
			t.Errorf("Expected log message %s, but got %s", tc.message, lastEntry.Message)
		}

		for k, v := range tc.details {
			if lastEntry.Data[k] != v {
				t.Errorf("Expected detail %s to be %v, but got %v", k, v, lastEntry.Data[k])
			}
		}
	}
}
