package httproxy

import (
	"net/http"

	"github.com/libs4go/scf4go"
	"github.com/libs4go/smf4go"
)

type httpProxyImpl struct{}

func New(config scf4go.Config) (smf4go.Runnable, error) {
	return &httpProxyImpl{}, nil
}

func (proxy *httpProxyImpl) Start() error {
	go http.ListenAndServe(":8080", proxy)
	return nil
}

func (proxy *httpProxyImpl) ServeHTTP(http.ResponseWriter, *http.Request) {

}
