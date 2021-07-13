package identity

import "github.com/libs4go/errors"

// ScopeOfAPIError .
const errVendor = "trust-identity"

// errors
var (
	ErrSession   = errors.New("Session not found", errors.WithVendor(errVendor))
	ErrTypedData = errors.New("TypedData sign verify failed", errors.WithVendor(errVendor))
)

type Authenticator interface {
	GetID(session string) (did string, err error)
	EthSignPrepare(from string, chainId uint) (typedData []byte, err error)
	EthSignVerify(session string, chainId uint, signature []byte) (err error)
}

type Token struct {
	DID     string `json:"did"`
	Session string `json:"session"`
}

type SessionManager interface {
	Get(session string) (token *Token, committed bool)
	Prepare(token *Token) error
	Commit(session string) error
}
