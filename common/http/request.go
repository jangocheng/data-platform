package HttpRequest

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
)

type Request struct {
	cli               *http.Client
	debug             bool
	url               string
	method            string
	timeout           time.Duration
	headers           map[string]string
	cookies           map[string]string
	data              interface{}
	params            map[string]string
	disableKeepAlives bool
	disableRedirect   bool
	tlsClientConfig   *tls.Config
	proxy             string
	json              bool
}

type RequestConfig struct {
	Url 	 		string
	Params      	map[string]string
	Data        	interface{}
	Headers     	map[string]string
	Cookies     	map[string]string
	Timeout     	time.Duration
	DisableAlive  	bool
	DisRedirect     bool
	DisTlsVerify    bool
	Proxy           string
}

func NewRequest(cnf ...RequestConfig) *Request {
	if len(cnf) == 0 {
		return &Request{}
	}
	conf := cnf[0]
 	req := &Request{url: conf.Url,timeout:conf.Timeout,
		headers:conf.Headers,cookies:conf.Cookies,
		disableKeepAlives:conf.DisableAlive,data:conf.Data,
		disableRedirect:conf.DisRedirect,params:conf.Params,
		tlsClientConfig:&tls.Config{InsecureSkipVerify:conf.DisTlsVerify}}
	return req
}

func (r *Request) SetDisRedirect(v bool) *Request {
	r.disableRedirect = v
	return r
}

func (r *Request) SetProxy(proxy string) *Request {
	r.proxy = proxy
	return r
}

func (r *Request) SetKeepAlives(v bool) *Request {
	r.disableKeepAlives = v
	return r
}

func (r *Request) SetTLSClient(v *tls.Config) *Request {
	r.tlsClientConfig = v
	return r
}

// Debug model
func (r *Request) SetDebug(v bool) *Request {
	r.debug = v
	return r
}

func (r *Request) SetTimeout(d time.Duration) *Request {
	r.timeout = d
	return r
}

// Set headers
func (r *Request) SetHeaders(h map[string]string) *Request {
	r.headers = h
	return r
}

// Set cookies
func (r *Request) SetCookies(c map[string]string) *Request {
	r.cookies = c
	return r
}

// Set cookies
func (r *Request) SetBody(d interface{}) *Request {
	r.data = d
	return r
}

// Set parmas
func (r *Request) SetParams(p map[string]string) *Request {
	r.params = p
	return r
}

func GetProxyFunc(r *Request) func(*http.Request) (*url.URL, error) {
	if r.proxy != "" {
		return func(r2 *http.Request) (*url.URL, error) {
			return url.Parse(r.proxy)
		}
	} else {
		return func(r2 *http.Request) (*url.URL, error) {
			return nil, nil
		}
	}
}

func CheckRedireciFunc(r *Request) func(*http.Request, []*http.Request) error {
	if r.disableRedirect {
		return func(req *http.Request, via []*http.Request) error {
			return errors.New("disable redirect")
		}
	} else {
		return nil
	}
}

// Build client
func (r *Request) buildClient() *http.Client {
	if r.cli == nil {
		r.cli = &http.Client{
			Transport: &http.Transport{
				Proxy: GetProxyFunc(r),
				DialContext: (&net.Dialer{
					Timeout:   r.timeout * time.Second,
				}).DialContext,
				TLSClientConfig:       r.tlsClientConfig,
				DisableKeepAlives:     r.disableKeepAlives,
				ResponseHeaderTimeout: r.timeout * time.Second,
			},
			CheckRedirect: CheckRedireciFunc(r),
		}
	}
	return r.cli
}

// Init headers
func (r *Request) initHeaders(req *http.Request) {
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	if r.json {
		req.Header.Set("Content-Type", "application/json")
		return
	}
	if r.data != nil {
		t := reflect.TypeOf(r.data).String()
		if strings.Contains(t, "map[") {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
}

// Init cookies
func (r *Request) initCookies(req *http.Request) {
	for k, v := range r.cookies {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
}

// Parse query for GET request
func parseQuery(url string) ([]string, error) {
	urlList := strings.Split(url, "?")
	if len(urlList) < 2 {
		return make([]string, 0), nil
	}
	query := make([]string, 0)
	for _, val := range strings.Split(urlList[1], "&") {
		v := strings.Split(val, "=")
		if len(v) < 2 {
			return make([]string, 0), errors.New("query parameter error")
		}
		query = append(query, fmt.Sprintf("%s=%s", v[0], v[1]))
	}
	return query, nil
}

// Build query data
func (r *Request) buildBody() (io.Reader, error) {
	// GET and DELETE request dose not send body
	if r.method == "GET" || r.method == "DELETE" {
		return nil, nil
	}

	if r.data == nil {
		return strings.NewReader(""), nil
	}

	t := reflect.TypeOf(r.data).String()
	tKind := strings.ToLower(reflect.TypeOf(r.data).Kind().String())
	tElemKind := strings.ToLower(reflect.TypeOf(r.data).Elem().Kind().String())
	if t != "string" && !strings.Contains(t, "map[") && t != "[]byte" && tKind != "struct" && tElemKind != "struct" {
		return nil, errors.New("incorrect parameter format.")
	}

	if r.json && (strings.Contains(t, "map[") || tKind == "struct" || tElemKind == "struct") {
		if b, err := json.Marshal(r.data); err != nil {
			return nil, err
		} else {
			return bytes.NewReader(b), nil
		}
	} else if t == "string" {
		return strings.NewReader(r.data.(string)), nil
	} else if t == "[]byte" {
		return bytes.NewReader(r.data.([]byte)), nil
	} else {
		dataMap := make(map[string]interface{})
		dataStr, err := json.Marshal(r.data)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dataStr, &dataMap)
		if err != nil {
			return nil, err
		}
		data := make([]string, 0)
		for k, v := range dataMap {
			if s, ok := v.(string); ok {
				data = append(data, fmt.Sprintf("%s=%v", k, s))
				continue
			}
			b, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			data = append(data, fmt.Sprintf("%s=%s", k, string(b)))
		}
		return strings.NewReader(strings.Join(data, "&")), nil
	}
}


// Build GET request url
func (r *Request) buildUrl(url string) (string, error) {
	query, err := parseQuery(url)
	if err != nil {
		return url, err
	}

	if r.params != nil {
		for k, v := range r.params {
			query = append(query, fmt.Sprintf("%s=%s", k, v))
		}
	}
	list := strings.Split(url, "?")

	urlPath := strings.TrimRight(list[0], "/")

	if len(query) > 0 {
		return fmt.Sprintf("%s?%s", urlPath, strings.Join(query, "&")), nil
	}

	return urlPath, nil
}

func (r *Request) elapsedTime(n int64, resp *Response) {
	end := time.Now().UnixNano() / 1e6
	resp.time = end - n
}

func (r *Request) log() {
	if r.debug {
		fmt.Printf("[HttpRequest]\n")
		fmt.Printf("-------------------------------------------------------------------\n")
		fmt.Printf("Request: %s %s\nHeaders: %v\nCookies: %v\nTimeout: %ds\nReqBody: %v\n\n", r.method, r.url, r.headers, r.cookies, r.timeout, r.data)
		//fmt.Printf("-------------------------------------------------------------------\n\n")
	}
}

// Get is a get http request
func (r *Request) Get(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	return r.request(http.MethodGet, r.url)
}

// Post is a post http request
func (r *Request) Post(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	return r.request(http.MethodPost,  r.url)
}

// Post is a post http request
func (r *Request) PostJson(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	r.json = true
	return r.request(http.MethodPost,  r.url)
}

// Put is a put http request
func (r *Request) Put(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	return r.request(http.MethodPut,  r.url)
}

// Delete is a delete http request
func (r *Request) Delete(url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	return r.request(http.MethodDelete,  r.url)
}

// Upload file
func (r *Request) Upload(filename, fileinput string, url ...string) (*Response, error) {
	if len(url) > 0 {
		r.url = url[0]
	}
	return r.sendFile(r.url, filename, fileinput)
}

// Send file
func (r *Request) sendFile(url, filename, fileinput string) (*Response, error) {
	if url == "" {
		return nil, errors.New("parameter url is required")
	}

	fileBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(fileBuffer)
	fileWriter, er := bodyWriter.CreateFormFile(fileinput, filename)
	if er != nil {
		return nil, er
	}

	f, er := os.Open(filename)
	if er != nil {
		return nil, er
	}
	defer f.Close()

	_, er = io.Copy(fileWriter, f)
	if er != nil {
		return nil, er
	}

	contentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()

	// Build Response
	response := &Response{}

	// Start time
	start := time.Now().UnixNano() / 1e6
	// Count elapsed time
	defer r.elapsedTime(start, response)

	// Debug infomation
	defer r.log()
	r.data = nil

	var (
		err error
		req *http.Request
	)
	r.cli = r.buildClient()
	r.method = "POST"

	req, err = http.NewRequest(r.method, url, fileBuffer)
	if err != nil {
		return nil, err
	}

	r.initHeaders(req)
	r.initCookies(req)
	req.Header.Set("Content-Type", contentType)

	resp, err := r.cli.Do(req)
	if err != nil {
		return nil, err
	}

	response.url = url
	response.resp = resp

	return response, nil
}

// Send http request
func (r *Request) request(method, url string) (*Response, error) {
	// Build Response
	response := &Response{}

	// Start time
	start := time.Now().UnixNano() / 1e6
	// Count elapsed time
	defer r.elapsedTime(start, response)

	if method == "" || url == "" {
		return nil, errors.New("parameter method and url is required")
	}

	// Debug infomation
	defer r.log()
	var (
		err  error
		req  *http.Request
		body io.Reader
	)
	r.cli = r.buildClient()

	method = strings.ToUpper(method)
	r.method = method

	url, err = r.buildUrl(url)
	if err != nil {
		return nil, err
	}
	r.url = url

	body, err = r.buildBody()
	if err != nil {
		return nil, err
	}
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	r.initHeaders(req)
	r.initCookies(req)

	resp, err := r.cli.Do(req)

	if err != nil {
		return nil, err
	}

	response.url = url
	response.resp = resp

	return response, nil
}
