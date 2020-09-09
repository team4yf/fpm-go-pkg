package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

//ResponseWrapper the wrapper of the http response
type ResponseWrapper struct {
	StatusCode int
	Success    bool
	Err        error
	Body       []byte
	Header     http.Header
}

type HTTPAuthType string
type HTTPAuthData map[string]interface{}

const (
	TypeHTTPAuthBasic = "basic"
)

//HTTPAuth http auth struct
type HTTPAuth struct {
	Type HTTPAuthType
	Data HTTPAuthData
}

//ConvertBody Convert the body to another struct
func (rsp *ResponseWrapper) ConvertBody(data interface{}) (err error) {
	if json.Unmarshal(rsp.Body, data); err != nil {
		return
	}
	return
}

//GetStringBody get the string of the body
func (rsp *ResponseWrapper) GetStringBody() string {

	return (string)(rsp.Body)
}

//Get send a get request with timeout
func Get(url string, timeout int) ResponseWrapper {
	return get(url, timeout, nil, nil)
}

//GetWithAuth send a get request with timeout
func GetWithAuth(url string, timeout int, auth *HTTPAuth) ResponseWrapper {
	return get(url, timeout, nil, auth)
}

//GetWithHeader send a get request with header and timeout
func GetWithHeader(url string, timeout int, header map[string]string) ResponseWrapper {
	return get(url, timeout, header, nil)
}

//GetWithHeaderAndAuth send a get request with header and timeout
func GetWithHeaderAndAuth(url string, timeout int, header map[string]string, auth *HTTPAuth) ResponseWrapper {
	return get(url, timeout, header, auth)
}

//PostParams post a form request with timeout
func PostParams(url string, params string, timeout int) ResponseWrapper {
	return post(url, "application/x-www-form-urlencoded", []byte(params), timeout, nil, nil)
}

//PostParamsWithHeader post a form request with timeout
func PostParamsWithHeader(url string, params string, timeout int, header map[string]string) ResponseWrapper {
	return post(url, "application/x-www-form-urlencoded", []byte(params), timeout, header, nil)
}

//PostParamsWithHeaderAndAuth post a form request with timeout
func PostParamsWithHeaderAndAuth(url string, params string, timeout int, header map[string]string, auth *HTTPAuth) ResponseWrapper {
	return post(url, "application/x-www-form-urlencoded", []byte(params), timeout, header, auth)
}

//PostJSON post a json data request with timeout
func PostJSON(url string, body []byte, timeout int) ResponseWrapper {
	return post(url, "application/json", body, timeout, nil, nil)
}

//PostJSONWithHeade post a json data request with timeout
func PostJSONWithHeade(url string, body []byte, timeout int, header map[string]string) ResponseWrapper {
	return post(url, "application/json", body, timeout, header, nil)
}

//PostJSONWithHeaderAndAuth post a json data request with timeout
func PostJSONWithHeaderAndAuth(url string, body []byte, timeout int, header map[string]string, auth *HTTPAuth) ResponseWrapper {
	return post(url, "application/json", body, timeout, header, auth)
}

//post post json & header with timeout
func post(url, contentType string, body []byte, timeout int, header map[string]string, auth *HTTPAuth) ResponseWrapper {
	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return createRequestError(err)
	}
	req.Header.Set("Content-Type", contentType)
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	return request(req, timeout, auth)
}

func get(url string, timeout int, header map[string]string, auth *HTTPAuth) ResponseWrapper {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return createRequestError(err)
	}
	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	return request(req, timeout, auth)
}

func request(req *http.Request, timeout int, auth *HTTPAuth) ResponseWrapper {
	wrapper := ResponseWrapper{StatusCode: 0, Header: make(http.Header)}
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	setRequestHeader(req)
	if auth != nil {
		switch auth.Type {
		case TypeHTTPAuthBasic:
			req.SetBasicAuth(auth.Data["username"].(string), auth.Data["password"].(string))
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		wrapper.Err = fmt.Errorf("执行HTTP请求错误-%s", err.Error())
		return wrapper
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wrapper.Err = fmt.Errorf("读取HTTP请求返回值失败-%s", err.Error())
		return wrapper
	}
	wrapper.StatusCode = resp.StatusCode
	wrapper.Body = body
	wrapper.Header = resp.Header

	return wrapper
}

func setRequestHeader(req *http.Request) {
	req.Header.Set("User-Agent", "fpm-go-pkg")
}

func createRequestError(err error) ResponseWrapper {
	errorMessage := errors.Wrap(err, "创建HTTP请求错误")
	return ResponseWrapper{StatusCode: 0, Success: false, Err: errorMessage}
}
