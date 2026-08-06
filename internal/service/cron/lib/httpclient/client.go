package httpclient

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
	"io"
	urlParse "net/url"
	"strings"
	"time"
)

// maxResponseBodySize 响应体大小上限（10MB），防止大响应耗尽内存
const maxResponseBodySize = 10 * 1024 * 1024

/*
	command 示例：
	{
	  "url": "http://www.baidu.com",
	  "method": "GET"
	}
	params 示例（headers/query 为行结构，data 为原始请求体字符串）：
	{
	  "headers": [{"key":"Content-Type","value":"application/json","enabled":true}],
	  "query": [{"key":"page","value":"1","enabled":true}],
	  "data": "{\"foo\":\"bar\"}"
	}
*/

// Result HTTP请求完整结果
type Result struct {
	Output     string        // 响应体
	StatusCode int           // 状态码
	Duration   time.Duration // 耗时
	Size       int           // 响应体字节数
}

// paramRow 参数行（key-value 行结构，支持启用/停用与描述）
type paramRow struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
	Desc    string `json:"desc,omitempty"`
}

// requestParams 解析后的请求参数
type requestParams struct {
	headers map[string]string
	query   urlParse.Values
	body    string
}

// Do 执行HTTP请求并返回完整结果（含状态码、耗时、响应大小）
func Do(method, targetURL, params string, timeout time.Duration) (*Result, error) {
	client, reqURL, rp, err := parseParams(targetURL, params, timeout)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	var r *gclient.Response
	switch strings.ToUpper(method) {
	case "POST":
		r, err = client.Post(gctx.GetInitCtx(), reqURL, rp.body)
	default:
		r, err = client.Get(gctx.GetInitCtx(), reqURL)
	}
	duration := time.Since(start)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	body, readErr := io.ReadAll(io.LimitReader(r.Body, maxResponseBodySize+1))
	if readErr != nil {
		return nil, fmt.Errorf("读取响应失败: %v", readErr)
	}
	size := len(body)
	if size > maxResponseBodySize {
		return nil, fmt.Errorf("响应体超过大小限制(%d字节)", maxResponseBodySize)
	}

	result := &Result{
		Output:     string(body),
		StatusCode: r.StatusCode,
		Duration:   duration,
		Size:       size,
	}
	if r.StatusCode >= 400 {
		return result, fmt.Errorf("HTTP请求失败，状态码: %d", r.StatusCode)
	}
	return result, nil
}

func Get(url, params string, timeout time.Duration) (string, error) {
	result, err := Do("GET", url, params, timeout)
	if err != nil {
		return "", err
	}
	return result.Output, nil
}

func Post(url, params string, timeout time.Duration) (string, error) {
	result, err := Do("POST", url, params, timeout)
	if err != nil {
		return "", err
	}
	return result.Output, nil
}

func parseParams(targetURL, params string, timeout time.Duration) (*gclient.Client, string, *requestParams, error) {
	client := g.Client()
	client.SetTimeout(timeout)
	rp := &requestParams{}
	if params == "" {
		return client, targetURL, rp, nil
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(params), &data)
	if err != nil {
		return nil, targetURL, nil, err
	}

	// headers：行结构 [{key,value,enabled}]
	if raw, ok := data["headers"]; ok && raw != nil {
		rows, err := parseRows(raw, "headers")
		if err != nil {
			return nil, targetURL, nil, err
		}
		rp.headers = make(map[string]string, len(rows))
		for _, row := range rows {
			if row.Enabled && row.Key != "" {
				rp.headers[row.Key] = row.Value
				client.SetHeader(row.Key, row.Value)
			}
		}
	}

	// query：行结构，拼接到 URL
	if raw, ok := data["query"]; ok && raw != nil {
		rows, err := parseRows(raw, "query")
		if err != nil {
			return nil, targetURL, nil, err
		}
		rp.query = urlParse.Values{}
		for _, row := range rows {
			if row.Enabled && row.Key != "" {
				rp.query.Add(row.Key, row.Value)
			}
		}
		if len(rp.query) > 0 {
			if strings.Contains(targetURL, "?") {
				targetURL += "&" + rp.query.Encode()
			} else {
				targetURL += "?" + rp.query.Encode()
			}
		}
	}

	// data：原始请求体字符串
	if raw, ok := data["data"]; ok && raw != nil {
		bodyStr, ok := raw.(string)
		if !ok {
			return nil, targetURL, nil, fmt.Errorf("data 必须是字符串")
		}
		rp.body = bodyStr
	}

	// 未显式设置 Content-Type 且存在请求体时，默认按 JSON 发送，
	// 避免被编码为表单格式导致 JSON 接口返回 415
	if rp.body != "" && !hasHeaderKey(rp.headers, "Content-Type") {
		client.SetHeader("Content-Type", "application/json")
	}
	return client, targetURL, rp, nil
}

// parseRows 解析 key-value 行数组
func parseRows(raw interface{}, name string) ([]paramRow, error) {
	rawJSON, err := json.Marshal(raw)
	if err != nil {
		return nil, fmt.Errorf("%s 解析失败: %v", name, err)
	}
	var rows []paramRow
	if err := json.Unmarshal(rawJSON, &rows); err != nil {
		return nil, fmt.Errorf("%s 必须是行数组: %v", name, err)
	}
	return rows, nil
}

// hasHeaderKey 判断 headers 中是否包含指定请求头（大小写不敏感）
func hasHeaderKey(headers map[string]string, key string) bool {
	for k := range headers {
		if strings.EqualFold(k, key) {
			return true
		}
	}
	return false
}
