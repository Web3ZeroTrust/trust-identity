package httproxy

import (
	"net/http"

	"github.com/libs4go/scf4go"
	"github.com/libs4go/slf4go"
	"github.com/libs4go/smf4go"
)

type httpProxyImpl struct {
	slf4go.Logger
	config         scf4go.Config
	metaMethodPath string
}

func New(config scf4go.Config) (smf4go.Runnable, error) {
	return &httpProxyImpl{
		Logger:         slf4go.Get("trust-identity-httproxy"),
		config:         config,
		metaMethodPath: config.Get("trust-identity-path").String("/trust-identity"),
	}, nil
}

func (proxy *httpProxyImpl) Start() error {
	addr := proxy.config.Get("listen").String(":8080")
	proxy.I("httproxy listen on {@addr}", addr)

	http.HandleFunc(proxy.metaMethodPath, proxy.dispatchMethod)

	go http.ListenAndServe(addr, nil)

	return nil
}

func (proxy *httpProxyImpl) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	proxy.I("{@path}", req.RequestURI)

	proxy.proxyDispatch(resp, req)
}

func (proxy *httpProxyImpl) dispatchMethod(resp http.ResponseWriter, req *http.Request) {
	proxy.I("{@path}", req.RequestURI)

	if req.Header.Get("Content-Type") != "application/json" {
		resp.WriteHeader(http.StatusBadRequest)
		resp.Write([]byte("jsonrpc expect !!!!!"))
		return
	}

}

func (proxy *httpProxyImpl) proxyDispatch(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusForbidden)
}
