package http

import (
    "fmt"
    "sync"
    "time"
    "io"
    "strings"

    "io/ioutil"
    "net/http"
    "crypto/tls"
    "github.com/golang/glog"
)

func NewHttpBench( )( *HttpBench ) {
    return &HttpBench {
        ConcurrentConns     :   1,
        Duration            :   10 * time.Second,
        Secure              :   false,
        InsecureSkipVerify  :   false,
        Method              :   http.MethodGet,
        Url                 :   "/",
        ReqBodySize         :   0,
        ReqBodyType         :   httpBenchReqBodyTypePlain,

        httpBenchCtx        :   httpBenchCtx {
            wg              :   &sync.WaitGroup{ }, 
            duration_chan   :   make( chan bool ),
        },
    }
}

func ( hBench *HttpBench )SetupHttpBenchReqBodyType( bodyType string ) {
    switch httpBenchReqBodyType( bodyType ) {
        case httpBenchReqBodyTypePlain:
            hBench.ReqBodyType = httpBenchReqBodyTypePlain

        case httpBenchReqBodyTypeJson:
            hBench.ReqBodyType = httpBenchReqBodyTypeJson

        default:
            hBench.ReqBodyType = httpBenchReqBodyTypePlain
    }
}

func ( hBench *HttpBench )Start( ) {
    err := hBench.setupRequest( )
    if err != nil {
        glog.Errorf( "Request setup failed with error %v", err )
        return
    }

    hBench.conns_chan = make( chan bool, hBench.ConcurrentConns  )

    hBench.startTest( )
    hBench.continueTest( )
    hBench.cleanup( )
}

func ( hBench *HttpBench )setupRequest( )( err error ) {
    if 0 == len( hBench.Host ) {
        return fmt.Errorf( "empty test host" )
    }

    url := hBench.Host
    if hBench.Secure {
        url  = "https://" + url

        if len( hBench.Port ) > 0 && "443" != hBench.Port {
            url += ":" + hBench.Port
        }
    } else {
        url  = "http://" + url

        if len( hBench.Port ) > 0 && "80" != hBench.Port {
            url += ":" + hBench.Port
        }
    }

    if strings.HasPrefix( hBench.Url, "/" ) {
        url += hBench.Url
    } else {
        url += "/" + hBench.Url
    }

    req, err := http.NewRequest( hBench.Method, url, nil )
    if err != nil {
        return err
    }

    req.Header.Add( "Accept", "*/*" )
    req.Header.Add( "Cache-Control", "no-cache" )

    hBench.request = req
    return nil
}

func ( hBench *HttpBench )startTest( ) {
    go func( ) {
        time.Sleep( hBench.Duration )
        hBench.duration_chan <- true
    }( )

    for i := uint( 0 ); i < hBench.ConcurrentConns; i++ {
        hBench.wg.Add( 1 )
        go hBench.runSingleTest( )
    }
}

func ( hBench *HttpBench )continueTest( ) {
    for {
        select {
            case <-hBench.duration_chan:
                return

            default:
        }

        select {
            case <-hBench.conns_chan:
                hBench.wg.Add( 1 )
                go hBench.runSingleTest( )
        }
    }
}

func ( hBench *HttpBench )cleanup( ) {
        hBench.wg.Wait( )
}

func ( hBench *HttpBench )runSingleTest( ) {
    httpClient := &http.Client {
        Transport: &http.Transport {
            MaxIdleConns:       1,
            IdleConnTimeout:    hBench.IdleTimeout,
            DisableCompression: true,
            TLSClientConfig:    &tls.Config {
                InsecureSkipVerify: hBench.InsecureSkipVerify,
            },
        },
    }

    defer func( ) {
        httpClient.CloseIdleConnections( )
        hBench.conns_chan <- true
        hBench.wg.Done( )
    }( )

    for i := uint( 0 ); i < hBench.ConnReqs; i++ {
        err := hBench.sendSingleRequest( httpClient )
        if err != nil {
            glog.Errorf( "Test failed with error %v", err )
            break
        }

        time.Sleep( hBench.ConnReqInterval )
    }
}

func ( hBench *HttpBench )sendSingleRequest( httpClient *http.Client )( err error ) {
    resp, err := httpClient.Do( hBench.request )
    if err != nil {
        return err
    }

    defer func( ) {
        io.Copy( ioutil.Discard, resp.Body )
        resp.Body.Close( )
    }( )

    if resp.StatusCode >= http.StatusBadRequest {
        return fmt.Errorf( "Request failed with response code = %v", resp.StatusCode )
    }

    return nil
}
