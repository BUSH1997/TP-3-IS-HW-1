package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/", ProxyHandler)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8000", nil)
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")
	r.RequestURI = ""

	rowURL := r.URL.String()
	if rowURL[len(rowURL) - 1] == '/' {
		r.URL, _ = url.Parse(rowURL[:len(rowURL) - 1])
	}

	resp, err := runProxyReq(r)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		panic(err)
	}
}

func runProxyReq(r *http.Request) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
