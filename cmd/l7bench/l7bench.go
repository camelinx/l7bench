package main

import (
    "fmt"
    "flag"
    "os"
    "time"

    l7bhttp "github.com/l7bench/internal/http"
)

var (
    version     string

    host                = flag.String( "test-host", "", "Hostname or IPv4 or IPv6 address being tested. This is a mandatory parameter" )
    port                = flag.String( "test-port", "", "Port serving traffic. Optional and defaults to 80 or 443" )
    method              = flag.String( "test-method", "GET", "HTTP request method. Optional and defaults to GET" )
    url                 = flag.String( "test-url", "/", "HTTP request url. Optional and defaults to /" )
    concurrent_conns    = flag.Uint( "concurrent-conns", uint( 50 ), "Concurrent connections to maintain with the server. Defaults to 50" )
    duration            = flag.Duration( "duration", 5 * time.Minute, "Total duration of the test. Defaults to 5 minutes" )
    conn_reqs           = flag.Uint( "conn-requests", uint( 2 ), "Number of requests to send in a TCP connection. Defaults to 2" )
    conn_req_intvl      = flag.Duration( "conn-req-interval", 1 * time.Second, "Interval in between requests on the same TCP connection. Defaults to 1 second" )
    idle_timeout        = flag.Duration( "idle-timeout", 30 * time.Second, "Timeout after which test will give up on an unresponsive server. Defaults to 30 seconds" )
    secure              = flag.Bool( "secure", false, "Send requests to TLS enabled server. Defaults to false" )
    disable_server_auth = flag.Bool( "disable-server-auth", false, "Disable TLS server authentication. Defaults to false" )
)

func main( ) {
    if len( *host ) == 0 {
        fmt.Println( "Test host cannot be empty" )
        os.Exit( 1 ) 
    }

    httpBench := l7bhttp.NewHttpBench( )

    httpBench.Host   = *host
    httpBench.Method = *method
    httpBench.Url    = *url

    httpBench.Secure             = *secure
    httpBench.InsecureSkipVerify = *disable_server_auth

    if len( *port ) > 0 {
        if *secure {
            if "443" != *port {
                httpBench.Port = *port
            }
        } else {
            if "80" != *port {
                httpBench.Port = *port
            }
        }
    }

    httpBench.ConcurrentConns = *concurrent_conns
    httpBench.Duration        = *duration
    httpBench.ConnReqs        = *conn_reqs
    httpBench.ConnReqInterval = *conn_req_intvl
    httpBench.IdleTimeout     = *idle_timeout

    httpBench.Start( )
}
