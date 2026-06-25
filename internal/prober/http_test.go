package prober

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 测试 HTTPProber 的成功探测且返回的状态码为 200 且 body 包含关键字的情况
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

// 测试 HTTPProber 的关键字不匹配的情况
func TestHTTPProbe_KeywordMismatch(t *testing.T) {
	// 启动本地服务器， 返回 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	p := HTTPProber{
		Keyword: "unhealthy",
	}
	res := p.Probe(server.URL, 3*time.Second)
	if res.Status != "不健康" {
		t.Errorf("期望不健康，得到 %s", res.Status)
	}
}

// 测试 HTTPProber 的返回状态码不为 200 的情况
func TestHTTPProbe_StatusCodeNot200(t *testing.T) {
	// 启动本地服务器， 返回 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	p := HTTPProber{}
	res := p.Probe(server.URL, 3*time.Second)
	if res.Status != "不健康" {
		t.Errorf("期望不健康，得到 %s", res.Status)
	}
}
