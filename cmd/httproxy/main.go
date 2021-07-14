package main

import (
	"github.com/libs4go/scf4go"
	_ "github.com/libs4go/scf4go/codec"
	_ "github.com/libs4go/slf4go/backend/console"
	"github.com/libs4go/smf4go"
	"github.com/libs4go/smf4go/app"
	"github.com/libs4go/smf4go/service/localservice"
	"github.com/web3zerotrust/trust-identity/app/httproxy"
	"github.com/web3zerotrust/trust-identity/auth"
	"github.com/web3zerotrust/trust-identity/session/localsession"
)

func main() {

	localservice.Register("identity.authenticator", func(config scf4go.Config) (smf4go.Service, error) {
		return auth.New(config)
	})

	localservice.Register("identity.sessionManager", func(config scf4go.Config) (smf4go.Service, error) {
		return localsession.New(config)
	})

	localservice.Register("identity.httproxy", func(config scf4go.Config) (smf4go.Service, error) {
		return httproxy.New(config)
	})

	app.Run("httproxy")
}
