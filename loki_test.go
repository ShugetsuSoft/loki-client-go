package loki_client_go

import (
	"github.com/ShugetsuSoft/loki-client-go/lib"
	"testing"
	"time"
)

func TestLokiClient(t *testing.T) {
	loki := NewLokiClient("http://localhost:3100/")
	loki.WriteLog(lib.Label{
		"type": "test",
		"mode": "consecutive",
	}, "Log1", time.Now())
	loki.WriteLog(lib.Label{
		"type": "test",
		"mode": "consecutive",
	}, "Log2", time.Now())
	loki.WriteLog(lib.Label{
		"type": "test",
		"mode": "consecutive",
	}, "Log3", time.Now())
	loki.WriteLog(lib.Label{
		"type": "test",
		"mode": "consecutive",
	}, "Log4", time.Now())
	err := loki.Push()
	if err != nil {
		t.Fatal(err)
	}
}
