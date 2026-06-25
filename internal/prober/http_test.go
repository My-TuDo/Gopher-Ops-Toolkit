package prober

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 测试创建 HTTP 请求失败
func TestHTTPProbe_BadURL(t *testing.T) {
	p := HTTPProber{}
	res := p.Probe("://bad-url", 3*time.Second)
	if res.Status != "不健康" {
		t.Errorf("期望不健康，得到 %s", res.Status)
	}
}

// 测试 发送请求失败
func TestHTTPProbe_ConnectionError(t *testing.T) {
	p := HTTPProber{}
	res := p.Probe("http://localhost:19999", 3*time.Second)
	if res.Status != "不健康" {
		t.Errorf("期望不健康，得到 %s", res.Status)
	}
}

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

// 测试 HTTPProber 的请求超时的情况
func TestHTTPProbe_Timeout(t *testing.T) {
	// 启动本地服务器，延迟响应
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := HTTPProber{}
	res := p.Probe(server.URL, 1*time.Second)
	if res.Status != "不健康" {
		t.Errorf("期望不健康，得到 %s", res.Status)
	}
}

// 测试 HTTPProber 的自定义请求头的情况
func TestHTTPProbe_CustomHeaders(t *testing.T) {
	// 启动本地服务器，检查请求头
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Header") != "test-value" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := HTTPProber{
		Headers: map[string]string{
			"X-Test-Header": "test-value", // 设置自定义请求头
		},
	}
	res := p.Probe(server.URL, 3*time.Second)
	if res.Status != "健康" {
		t.Errorf("期望健康，得到 %s", res.Status)
	}
}

// 测试 HTTPProber 的数据清洗功能，确保 URL 前后空格被去除
func TestHTTPProbe_DataCleaning(t *testing.T) {
	// 启动本地服务器，返回 200
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	p := HTTPProber{}
	// 在 URL 前后添加空格
	res := p.Probe("   "+server.URL+"   ", 3*time.Second)
	if res.Status != "健康" {
		t.Errorf("期望健康，得到 %s", res.Status)
	}
}
