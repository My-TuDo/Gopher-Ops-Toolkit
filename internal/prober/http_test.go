package prober

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPProbe_Success(t *testing.T) {
	// 启动本地服务器， 返回 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	p := HTTPProber{
		Keyword: "ok",
	}
	res := p.Probe(server.URL, 3*time.Second)
	if res.Status != "健康" {
		t.Errorf("期望健康，得到 %s", res.Status)
	}
}
