package httpclient

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
	urlParse "net/url"
	"strings"
	"time"
)

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

	fmt.Println(66666666666666666, r.StatusCode, r.ReadAllString())
	return r.ReadAllString(), nil
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
	return r.ReadAllString(), nil
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
		for key, value := range headers.(g.Map) {
			client.SetHeader(key, value.(string))
		}
	}

	query, ok := data["query"]
	if ok {
		v := urlParse.Values{}
		for key, value := range query.(g.Map) {
			v.Add(key, value.(string))
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
	return client, url, body.(g.Map), nil
}
