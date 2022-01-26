package logger

import (
	"bytes"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MemorySink implements zap.Sink by writing all messages to a buffer.
type MemorySink struct {
	*bytes.Buffer
}

// Implement Close and Sync as no-ops to satisfy the interface. The Write
// method is provided by the embedded buffer.
func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

var sink *MemorySink

func getLogger(level string) (*Logger, error) {
	sink = &MemorySink{new(bytes.Buffer)}
	_ = zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
		return sink, nil
	})

	return NewLogger(level, []string{"memory://"})
}

func TestLogger(t *testing.T) {
	defer sink.Close()

	t.Run("log without any errors", func(t *testing.T) {
		logg, err := getLogger("debug")
		require.NoError(t, err)
		logg.Debug("log without any errors")
	})

	t.Run("correct log", func(t *testing.T) {
		logg, err := getLogger("debug")
		require.NoError(t, err)

		logg.Debug("correct log",
			"ip", "66.249.65.3",
			"method", "GET",
			"path", "/hello?q=1",
			"httpVersion", "HTTP/1.1",
			"status", 200,
			"latency", 30,
			"userAgent", "Mozilla/5.0",
		)

		var values map[string]interface{}
		err = json.Unmarshal(sink.Bytes(), &values)
		require.NoError(t, err)
		require.Equal(t, "66.249.65.3", values["ip"])
		require.Equal(t, "DEBUG", values["lvl"])
		require.Equal(t, "correct log", values["msg"])
	})

	t.Run("no log with smaller level", func(t *testing.T) {
		logg, err := getLogger("error")
		require.NoError(t, err)

		logg.Debug("no log",
			"ip", "66.249.65.3",
			"method", "GET",
			"path", "/hello?q=1",
			"httpVersion", "HTTP/1.1",
			"status", 200,
			"latency", 30,
			"userAgent", "Mozilla/5.0",
		)
		require.Equal(t, "", sink.String())
	})
}
