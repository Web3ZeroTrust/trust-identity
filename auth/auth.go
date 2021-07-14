package auth

import (
	"encoding/json"

	"github.com/bwmarrin/snowflake"
	"github.com/libs4go/errors"
	"github.com/libs4go/scf4go"
	"github.com/libs4go/slf4go"
	identity "github.com/web3zerotrust/trust-identity"
)

type authenticatorImpl struct {
	slf4go.Logger
	SessionManager identity.SessionManager `inject:"identity.sessionManager"`
	node           *snowflake.Node
}

// New create new authenticator service
func New(config scf4go.Config) (identity.Authenticator, error) {

	node, err := snowflake.NewNode(int64(config.Get("cluster-id").Int(1)))

	if err != nil {
		return nil, errors.Wrap(err, "create snowflake node error")
	}

	return &authenticatorImpl{
		Logger: slf4go.Get("trust-identity-authenticator"),
		node:   node,
	}, nil
}

func (auth *authenticatorImpl) GetID(session string) (string, error) {
	token, ok := auth.SessionManager.Get(session)

	if !ok {
		return "", identity.ErrSession
	}

	return token.DID, nil

}

func (auth *authenticatorImpl) EthSignPrepare(from string, chainId uint) (typedData []byte, err error) {
	session := auth.node.Generate().String()

	token := &identity.Token{
		DID:     etherAddressToDID(from),
		Session: session,
	}

	data := newTypedData(token, chainId)

	buff, err := json.Marshal(data)

	if err != nil {
		return nil, errors.Wrap(err, "marshal typedData error")
	}

	return buff, nil

}

func (auth *authenticatorImpl) EthSignVerify(session string, chainId uint, signature []byte) (err error) {

	token, committed := auth.SessionManager.Get(session)

	if token == nil {
		return errors.Wrap(identity.ErrSession, "not found session: %s", session)
	}

	if committed {
		return nil
	}

	ok, err := verifyTypedData(token, chainId, signature)

	if err != nil {
		return err
	}

	if !ok {
		return errors.Wrap(identity.ErrTypedData, "session %s chainId %d typedData verify failed", session, chainId)
	}

	return auth.SessionManager.Commit(session)

}
