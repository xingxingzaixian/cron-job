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
		{
		  "url": "http://www.baidu.com",
		  "method": "GET"
	    }
		{
		  "headers": {},
		  "query": {},
		  "data": {},
		}
*/
func Get(url, params string, timeout time.Duration) (string, error) {
	client, url, _, err := parseClient(url, params, timeout)
	if err != nil {
		return "", err
	}

	r, err := client.Get(gctx.GetInitCtx(), url)
	if err != nil {
		return "", err
	}
	defer r.Close()

	return readResponse(r)
}

func Post(url, params string, timeout time.Duration) (string, error) {
	client, url, data, err := parseClient(url, params, timeout)
	if err != nil {
		return "", err
	}

	r, err := client.Post(gctx.GetInitCtx(), url, data)
	if err != nil {
		return "", err
	}
	defer r.Close()
	return readResponse(r)
}

// readResponse 校验HTTP状态码并读取响应体（带大小限制）
func readResponse(r *gclient.Response) (string, error) {
	if r.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP请求失败，状态码: %d", r.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxResponseBodySize+1))
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}
	if len(body) > maxResponseBodySize {
		return "", fmt.Errorf("响应体超过大小限制(%d字节)", maxResponseBodySize)
	}
	return string(body), nil
}

func parseClient(url, params string, timeout time.Duration) (*gclient.Client, string, g.Map, error) {
	client := g.Client()
	client.SetTimeout(timeout)
	if params == "" {
		return client, url, nil, nil
	}

	var data g.Map
	err := json.Unmarshal([]byte(params), &data)
	if err != nil {
		return nil, url, nil, err
	}

	headers, ok := data["headers"]
	if ok {
		headerMap, ok := headers.(map[string]interface{})
		if !ok {
			return nil, url, nil, fmt.Errorf("headers 必须是 JSON 对象")
		}
		for key, value := range headerMap {
			strValue, ok := value.(string)
			if !ok {
				return nil, url, nil, fmt.Errorf("header[%s] 的值必须是字符串", key)
			}
			client.SetHeader(key, strValue)
		}
	}

	query, ok := data["query"]
	if ok {
		queryMap, ok := query.(map[string]interface{})
		if !ok {
			return nil, url, nil, fmt.Errorf("query 必须是 JSON 对象")
		}
		v := urlParse.Values{}
		for key, value := range queryMap {
			strValue, ok := value.(string)
			if !ok {
				return nil, url, nil, fmt.Errorf("query[%s] 的值必须是字符串", key)
			}
			v.Add(key, strValue)
		}

		if strings.Contains(url, "?") {
			url += "&" + v.Encode()
		} else {
			url += "?" + v.Encode()
		}
	}

	body, ok := data["data"]
	if !ok {
		return client, url, nil, nil
	}
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil, url, nil, fmt.Errorf("data 必须是 JSON 对象")
	}
	return client, url, bodyMap, nil
}
