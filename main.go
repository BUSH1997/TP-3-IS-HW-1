package main

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type ProxyHandler struct {
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		handleHTTPS(w, r)
	}
}

func main() {
	handler := &ProxyHandler{}

	server := http.Server{
		Addr:    ":8000",
		Handler: handler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf(err.Error())
	}
}

func handleHTTPS(w http.ResponseWriter, r *http.Request) {
	destConn, err := connectHandshake(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("cannot establish connection")
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Connection established"))
	if err != nil {
		log.Fatal("cannot write data")
	}

	srcConn, err := connectHijacker(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("cannot establish connection")
	}

	go transferData(destConn, srcConn)
	go transferData(srcConn, destConn)
}

func connectHandshake(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return nil, err
	}

	return conn, nil
}

func connectHijacker(w http.ResponseWriter) (net.Conn, error) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "hijacking not supported", http.StatusInternalServerError)
		return nil, errors.New("hijacking not supported")
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return conn, nil
}

func transferData(dest io.WriteCloser, src io.ReadCloser) {
	defer func(dest io.WriteCloser) {
		err := dest.Close()
		if err != nil {

		}
	}(dest)
	defer func(src io.ReadCloser) {
		err := src.Close()
		if err != nil {

		}
	}(src)

	_, err := io.Copy(dest, src)
	if err != nil {
		return
	}
}
