package ai

import (
    "net"
    "net/http"
    "time"
)

func DefaultHTTPClient(timeout time.Duration) *http.Client {
    tr := &http.Transport{
        Proxy: http.ProxyFromEnvironment,
        DialContext: (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
        TLSHandshakeTimeout: 10 * time.Second,
        MaxIdleConns:        100,
        IdleConnTimeout:     90 * time.Second,
    }
    return &http.Client{Transport: tr, Timeout: timeout}
}
