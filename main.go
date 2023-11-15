package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

var (
	backend    string
	listenport string
	cert       string
	key        string
)

func init() {
	backend = os.Getenv("BACKEND_ADDRESS")
	if backend == "" {
		panic("BACKEND_ADDRESS environment variable not set")
	}
	listenport = os.Getenv("LISTEN_PORT")
	if listenport == "" {
		listenport = "8080"
	}
	cert = os.Getenv("SERVER_CERT")
	key = os.Getenv("SERVER_KEY")
}

func main() {
	proxy := &httputil.ReverseProxy{
		Transport: roundTripper(rt),
		Director: func(req *http.Request) {
			req.URL.Scheme = "https"
			req.URL.Host = backend
			req.Host = backend
			req.Header.Set("Host", backend)
		},
	}
	if cert != "" && key != "" {
		log.Println("starting tls server")
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%s", listenport), cert, key, proxy))
	} else {
		log.Println("starting http server")
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", listenport), proxy))
	}
}

func rt(req *http.Request) (*http.Response, error) {
	log.Printf("request received. url=%s", req.URL)
	req.Header.Set("Host", backend)
	defer log.Printf("request complete. url=%s", req.URL)

	var InsecureTransport http.RoundTripper = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   true,
	}

	if req.ContentLength > 0 {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Printf("failed to read body: %s", err)
			return &http.Response{}, err
		}
		rdr := ioutil.NopCloser(bytes.NewBuffer(b))
		//fmt.Printf("%s\n", b)
		if strings.Contains(string(b), "PowerOffVM_Task") {
			log.Println("poweroff event detected going to sleep for a bit")
			time.Sleep(1200 * time.Second)
			return &http.Response{}, fmt.Errorf("die horribly becuase we detected poweroff event")
		}
		req.Body = rdr
	}
	return InsecureTransport.RoundTrip(req)
	//return http.DefaultTransport.RoundTrip(req)
}

// roundTripper makes func signature a http.RoundTripper
type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
