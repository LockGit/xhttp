package process

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"xhttp/storage"
)

var HttpRequestFactory = map[string]IHttp{}

func init() {
	register(http.MethodPost, &PostJson{})
	register(http.MethodGet, &Get{})
}

type IHttp interface {
	DoRequest(api *storage.APIChildren, requestParams map[string]interface{}) (v string, err error)
}

func register(method string, o IHttp) {
	if _, ok := HttpRequestFactory[method]; !ok {
		HttpRequestFactory[method] = o
	}
}

type PostJson struct {
}

func (p *PostJson) DoRequest(api *storage.APIChildren, requestParams map[string]interface{}) (v string, err error) {
	url := api.Url
	var bs []byte
	bs, err = json.Marshal(requestParams)
	if err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(strings.ToUpper(api.Method), url, strings.NewReader(string(bs)))
	c := http.Client{
		Timeout: time.Duration(api.Timeout) * time.Second,
	}
	var resp *http.Response
	resp, err = c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return string(bodyBytes), nil
}

type Get struct {
}

func (g *Get) DoRequest(api *storage.APIChildren, requestParams map[string]interface{}) (v string, err error) {
	url := api.Url
	var req *http.Request
	req, err = http.NewRequest(strings.ToUpper(api.Method), url, nil)
	q := req.URL.Query()
	for k, val := range requestParams {
		q.Add(k, val.(string))
	}
	req.URL.RawQuery = q.Encode()
	c := http.Client{
		Timeout: time.Duration(api.Timeout) * time.Second,
	}
	var resp *http.Response
	log.Println("req.URL:::::::", req.URL, req.Method)
	resp, err = c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return string(bodyBytes), nil
}

func execRequest(api *storage.APIChildren, requestParams map[string]interface{}, data *Response) (err error) {
	method := strings.ToUpper(api.Method)
	if exec, ok := HttpRequestFactory[method]; ok {
		var val string
		val, err = exec.DoRequest(api, requestParams)
		if err != nil {
			return
		}
		data.Set(api.Name, val)
		return
	}
	return errors.New("no support method:" + method)
}
