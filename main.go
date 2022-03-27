package main

import (
    "log"
    "net/http"
    "net/http/httputil"
    "os"
    "fmt"
)

var (
    backend string
    listenport string
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
}

func main() {
    proxy := &httputil.ReverseProxy{
        Transport: roundTripper(rt),
        Director: func(req *http.Request) {
            req.URL.Scheme = "http"
            req.URL.Host = backend
            req.Host = backend
            req.Header.Set("Host", backend) 
        },
    }
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", listenport), proxy))
}

func rt(req *http.Request) (*http.Response, error) {
    log.Printf("request received. url=%s", req.URL)
    req.Header.Set("Host", backend) 
    defer log.Printf("request complete. url=%s", req.URL)

    return http.DefaultTransport.RoundTrip(req)
}


// roundTripper makes func signature a http.RoundTripper
type roundTripper func(*http.Request) (*http.Response, error)

func (f roundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }
