package main

import (
    "flag"
    "os"
    "time"
    "strconv"

    "github.com/golang/glog"
    l7bhttp "github.com/l7bench/internal/http"
)

var (
    version     string

    host                = flag.String( "test-host", "", "Hostname or IPv4 or IPv6 address being tested. This is a mandatory parameter" )
    port                = flag.String( "test-port", "", "Port serving traffic. Optional and defaults to 80 or 443" )
    method              = flag.String( "test-method", "GET", "HTTP request method. Optional and defaults to GET" )
    url                 = flag.String( "test-url", "/", "HTTP request url. Optional and defaults to /" )
    reqBodySize         = flag.Uint( "test-body-size", uint( 0 ), "HTTP request body size. Optional and defaults to 0 (no request body)" )
    reqBodyType         = flag.String( "test-body-type", "plain", "HTTP request body type \"plain\" or \"json\". Optional and defaults to plain text" )
    concurrentConns     = flag.Uint( "concurrent-conns", uint( 50 ), "Concurrent connections to maintain with the server. Defaults to 50" )
    duration            = flag.Duration( "duration", 5 * time.Minute, "Total duration of the test. Defaults to 5 minutes" )
    conn_reqs           = flag.Uint( "conn-requests", uint( 2 ), "Number of requests to send in a TCP connection. Defaults to 2" )
    connReqIntvl        = flag.Duration( "conn-req-interval", 1 * time.Second, "Interval in between requests on the same TCP connection. Defaults to 1 second" )
    idleTimeout         = flag.Duration( "idle-timeout", 30 * time.Second, "Timeout after which test will give up on an unresponsive server. Defaults to 30 seconds" )
    secure              = flag.Bool( "secure", false, "Send requests to TLS enabled server. Defaults to false" )
    disableServerAuth   = flag.Bool( "disable-server-auth", false, "Disable TLS server authentication. Defaults to false" )
)

func main( ) {
    flag.Parse( )

    err := flag.Lookup( "logtostderr" ).Value.Set( "true" )
    if err != nil {
        glog.Fatalf( "Error setting logtostderr to true: %v", err )
    }

    glog.Infof( "Starting l7bench %v", version )

    httpBench := l7bhttp.NewHttpBench( )

    setupString( &httpBench.Host, host, "L7BENCH_TEST_HOST" )
    if 0 == len( httpBench.Host ) {
        glog.Fatalf( "Test host cannot be empty" )
    }

    setupString( &httpBench.Method, method, "L7BENCH_TEST_METHOD" )
    setupString( &httpBench.Url, url, "L7BENCH_TEST_URL" )

    setupUint( &httpBench.ReqBodySize, reqBodySize, "L7BENCH_TEST_REQ_BODY_SIZE" )

    var finalReqBodyType string
    setupString( &finalReqBodyType, reqBodyType, "L7BENCH_TEST_REQ_BODY_TYPE" )
    httpBench.SetupHttpBenchReqBodyType( finalReqBodyType )

    setupBool( &httpBench.Secure, secure, "L7BENCH_TEST_SECURE" )
    setupBool( &httpBench.InsecureSkipVerify, disableServerAuth, "L7BENCH_TEST_DISABLE_SERVER_AUTH" )

    setupString( &httpBench.Port, port, "L7BENCH_TEST_PORT" )
    if httpBench.Secure {
        if "443" == httpBench.Port {
            httpBench.Port = ""
        }
    } else {
        if "80" == httpBench.Port {
            httpBench.Port = ""
        }
    }

    setupUint( &httpBench.ConcurrentConns, concurrentConns, "L7BENCH_TEST_CONCURRENT_CONNS" )
    setupUint( &httpBench.ConnReqs, conn_reqs, "L7BENCH_TEST_CONN_REQS" )

    setupDuration( &httpBench.Duration, duration, "L7BENCH_TEST_DURATION" )
    setupDuration( &httpBench.ConnReqInterval, connReqIntvl, "L7BENCH_TEST_CONN_REQ_INTERVAL" )
    setupDuration( &httpBench.IdleTimeout, idleTimeout, "L7BENCH_TEST_IDLE_TIMEOUT" )

    glog.Infof( "Starting l7bench test %+v", httpBench )
    httpBench.Start( )
}

func setupString( field, arg *string, envVar string ) {
    envVal := os.Getenv( envVar )
    if len( envVal ) > 0 {
        *field = envVal
        return
    }

    if arg != nil && len( *arg ) > 0 {
        *field = *arg
    }
}

func setupBool( field, arg *bool, envVar string ) {
    envVal := os.Getenv( envVar )
    if len( envVal ) > 0 {
        if boolVal, err := strconv.ParseBool( envVal ); nil == err {
            *field = boolVal
            return
        }
    }

    if arg != nil {
        *field = *arg
    }
}

func setupUint( field, arg *uint, envVar string ) {
    envVal := os.Getenv( envVar )
    if len( envVal ) > 0 {
        if uintVal, err := strconv.ParseUint( envVal, 10, 64 ); nil == err {
            *field = uint( uintVal )
            return
        }
    }

    if arg != nil {
        *field = *arg
    }
}

func setupDuration( field, arg *time.Duration, envVar string ) {
    envVal := os.Getenv( envVar )
    if len( envVal ) > 0 {
        if durVal, err := time.ParseDuration( envVal ); nil == err {
            *field = durVal
            return
        }
    }

    if arg != nil {
        *field = *arg
    }
} 
