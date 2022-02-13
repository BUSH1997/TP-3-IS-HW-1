package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ProxyRequest struct {
	Method string
	Host string
	UserAgent string
	Accept string
	ProxyConnection string
}

func main() {
	http.HandleFunc("/", ProxyHandler)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println(r.Header.Get("HOST"))
	fmt.Println(r.Header.Get("User-Agent"))
	fmt.Println(r.Header.Get("Accept"))
	fmt.Println(r.Header.Get("Proxy-Connection"))
	fmt.Println(r.RequestURI)

	proxyRequest := ProxyRequest{
		Method:          r.Method,
		Host:            r.RequestURI,
		UserAgent:       r.Header.Get("User-Agent"),
		Accept:          r.Header.Get("Accept"),
	}

	resp, err := runGetFullReq(proxyRequest)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	io.Copy(w, resp.Body)
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func runGetFullReq(request ProxyRequest) (*http.Response, error) {

	req := &http.Request{
		Method: request.Method,
		Header: http.Header{
			"User-Agent": {request.UserAgent},
			"Accept": {request.Accept},
			"Host": {request.Host},
		},
	}

	req.URL, _ = url.Parse(request.Host)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error happend", err)
		return nil, err
	}

	return resp, nil
}
