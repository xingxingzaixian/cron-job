package httpclient

import (
	"encoding/json"
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
	var data g.Map
	err := json.Unmarshal([]byte(params), &data)
	if err != nil {
		return nil, url, nil, err
	}
	client := g.Client()
	client.SetTimeout(timeout)

	headers, ok := data["headers"]
	if ok {
		client.SetHeaderMap(headers.(g.MapStrStr))
	}

	query, ok := data["query"]
	if ok {
		v := urlParse.Values{}
		for key, value := range query.(g.MapStrStr) {
			v.Add(key, value)
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
