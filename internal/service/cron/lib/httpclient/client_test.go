package httpclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// newJSONOnlyServer 模拟只接受 application/json 的接口：非 JSON Content-Type 返回 415。
func newJSONOnlyServer(t *testing.T, capture func(contentType, body string)) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		contentType := r.Header.Get("Content-Type")
		if capture != nil {
			capture(contentType, string(body))
		}
		if !strings.HasPrefix(contentType, "application/json") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// 用户只填请求体、没设置 Content-Type 头时，应默认按 JSON 发送（415 回归测试）。
func TestPostDataDefaultsToJSON(t *testing.T) {
	var gotContentType, gotBody string
	srv := newJSONOnlyServer(t, func(contentType, body string) {
		gotContentType = contentType
		gotBody = body
	})

	params := `{"data":"{\"foo\":\"bar\"}"}`
	if _, err := Post(srv.URL, params, 5*time.Second); err != nil {
		t.Fatalf("Post 返回错误: %v (Content-Type=%q body=%q)", err, gotContentType, gotBody)
	}
	if !strings.HasPrefix(gotContentType, "application/json") {
		t.Fatalf("期望默认 JSON Content-Type, 实际 %q", gotContentType)
	}
	if gotBody != `{"foo":"bar"}` {
		t.Fatalf("期望 body %q, 实际 %q", `{"foo":"bar"}`, gotBody)
	}
}

// 显式设置 JSON Content-Type 头（行结构）应正常发送。
func TestPostJSONWithExplicitHeader(t *testing.T) {
	var gotContentType, gotBody string
	srv := newJSONOnlyServer(t, func(contentType, body string) {
		gotContentType = contentType
		gotBody = body
	})

	params := `{"headers":[{"key":"Content-Type","value":"application/json","enabled":true}],"data":"{\"foo\":\"bar\"}"}`
	if _, err := Post(srv.URL, params, 5*time.Second); err != nil {
		t.Fatalf("Post 返回错误: %v", err)
	}
	if !strings.HasPrefix(gotContentType, "application/json") {
		t.Fatalf("期望 JSON Content-Type, 实际 %q", gotContentType)
	}
	if gotBody != `{"foo":"bar"}` {
		t.Fatalf("期望 body %q, 实际 %q", `{"foo":"bar"}`, gotBody)
	}
}

// 显式设置带 charset 的 JSON Content-Type 也应正常。
func TestPostJSONWithCharsetHeader(t *testing.T) {
	srv := newJSONOnlyServer(t, nil)
	params := `{"headers":[{"key":"Content-Type","value":"application/json; charset=utf-8","enabled":true}],"data":"{\"foo\":\"bar\"}"}`
	if _, err := Post(srv.URL, params, 5*time.Second); err != nil {
		t.Fatalf("Post 返回错误: %v", err)
	}
}

// 显式设置表单 Content-Type 时应保持表单编码，不强制 JSON。
func TestPostFormWithExplicitHeader(t *testing.T) {
	var gotContentType, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		gotContentType = r.Header.Get("Content-Type")
		gotBody = string(body)
		if !strings.HasPrefix(gotContentType, "application/x-www-form-urlencoded") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	params := `{"headers":[{"key":"Content-Type","value":"application/x-www-form-urlencoded","enabled":true}],"data":"foo=bar"}`
	if _, err := Post(srv.URL, params, 5*time.Second); err != nil {
		t.Fatalf("Post 返回错误: %v", err)
	}
	if gotBody != "foo=bar" {
		t.Fatalf("期望表单 body %q, 实际 %q", "foo=bar", gotBody)
	}
}

// 禁用的 Content-Type 行不应生效，未显式设置时默认按 JSON 发送。
func TestDisabledHeaderRowNotApplied(t *testing.T) {
	var gotContentType string
	srv := newJSONOnlyServer(t, func(contentType, _ string) {
		gotContentType = contentType
	})

	params := `{"headers":[{"key":"Content-Type","value":"text/plain","enabled":false}],"data":"hello"}`
	if _, err := Post(srv.URL, params, 5*time.Second); err != nil {
		t.Fatalf("Post 返回错误: %v (Content-Type=%q)", err, gotContentType)
	}
	if !strings.HasPrefix(gotContentType, "application/json") {
		t.Fatalf("期望默认 JSON Content-Type, 实际 %q", gotContentType)
	}
}

// Do 应返回状态码、耗时与响应大小。
func TestDoReturnsStatusDurationAndSize(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))
	t.Cleanup(srv.Close)

	result, err := Do("GET", srv.URL, "", 5*time.Second)
	if err != nil {
		t.Fatalf("Do 返回错误: %v", err)
	}
	if result.StatusCode != http.StatusOK {
		t.Fatalf("期望状态码 200, 实际 %d", result.StatusCode)
	}
	if result.Output != "hello" {
		t.Fatalf("期望输出 %q, 实际 %q", "hello", result.Output)
	}
	if result.Size != len("hello") {
		t.Fatalf("期望大小 %d, 实际 %d", len("hello"), result.Size)
	}
	if result.Duration <= 0 {
		t.Fatalf("期望耗时大于 0, 实际 %v", result.Duration)
	}
}

// 4xx/5xx 时 Do 返回结果（带响应体）与错误。
func TestDoErrorStatusReturnsResultAndError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("boom"))
	}))
	t.Cleanup(srv.Close)

	result, err := Do("GET", srv.URL, "", 5*time.Second)
	if err == nil {
		t.Fatal("期望返回错误")
	}
	if result == nil {
		t.Fatal("期望返回结果")
	}
	if result.StatusCode != http.StatusInternalServerError {
		t.Fatalf("期望状态码 500, 实际 %d", result.StatusCode)
	}
	if result.Output != "boom" {
		t.Fatalf("期望错误响应体 %q, 实际 %q", "boom", result.Output)
	}
	if !strings.Contains(err.Error(), "状态码: 500") {
		t.Fatalf("错误信息缺少状态码: %v", err)
	}
}

// query 行结构应拼接到 URL，禁用的行跳过。
func TestQueryRowsAppended(t *testing.T) {
	var gotRawQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRawQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	params := `{"query":[{"key":"page","value":"1","enabled":true},{"key":"size","value":"2","enabled":true},{"key":"skip","value":"3","enabled":false}]}`
	if _, err := Get(srv.URL, params, 5*time.Second); err != nil {
		t.Fatalf("Get 返回错误: %v", err)
	}
	if gotRawQuery != "page=1&size=2" {
		t.Fatalf("期望 query %q, 实际 %q", "page=1&size=2", gotRawQuery)
	}
}

// headers 中存在空 key 的行（UI 添加后未填写）不应导致 invalid header field name 错误。
func TestEmptyHeaderKeySkipped(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	params := `{"headers":[{"key":"","value":"x","enabled":true}],"data":"hello"}`
	if _, err := Post(srv.URL, params, 5*time.Second); err != nil {
		t.Fatalf("Post 返回错误: %v", err)
	}
}
