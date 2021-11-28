package http

import (
    "time"
    "sync"

    "net/http"
)

type httpBenchCtx struct {
    wg                   *sync.WaitGroup
    conns_chan            chan bool
    duration_chan         chan bool

    request              *http.Request
    requestBody        [ ]byte 
}

type httpBenchReqBodyType string

const (
    httpBenchReqBodyTypePlain httpBenchReqBodyType = "plain"

    httpBenchReqBodyTypeJson  httpBenchReqBodyType = "json"
)

type HttpBench struct {
    ConcurrentConns       uint
    Duration              time.Duration
    ConnReqs              uint
    ConnReqInterval       time.Duration        
    IdleTimeout           time.Duration

    Secure                bool
    InsecureSkipVerify    bool

    Host                  string
    Port                  string
    Method                string
    Url                   string
    ReqBodySize           uint
    ReqBodyType           httpBenchReqBodyType

    httpBenchCtx
}
