package localsession

import (
	"sync"
	"time"

	"github.com/libs4go/scf4go"
	"github.com/libs4go/slf4go"
	identity "github.com/web3zerotrust/trust-identity"
)

type sessionImpl struct {
	token       identity.Token
	committed   bool
	expiredTime time.Time
}

type sessionManagerImpl struct {
	sync.RWMutex
	slf4go.Logger
	prepareDuration time.Duration
	commitDuration  time.Duration
	sessions        map[string]*sessionImpl
}

// New create new authenticator service
func New(config scf4go.Config) (identity.SessionManager, error) {

	return &sessionManagerImpl{
		Logger:          slf4go.Get("trust-identity-localsession"),
		prepareDuration: config.Get("timeout", "prepare").Duration(10 * time.Second),
		commitDuration:  config.Get("timeout", "commit").Duration(7 * 24 * time.Hour),
		sessions:        make(map[string]*sessionImpl),
	}, nil
}

func (s *sessionManagerImpl) Get(session string) (*identity.Token, bool) {
	s.RLock()
	impl, ok := s.sessions[session]
	s.RUnlock()

	if !ok {
		return nil, false
	}

	if impl.expiredTime.Before(time.Now()) {
		s.Lock()
		delete(s.sessions, session)
		s.Unlock()

		return nil, false
	}

	token := impl.token

	return &token, impl.committed
}
func (s *sessionManagerImpl) Prepare(token *identity.Token) error {
	session := &sessionImpl{
		token:       *token,
		expiredTime: time.Now().Add(s.prepareDuration),
	}

	s.Lock()
	defer s.Unlock()

	s.sessions[token.Session] = session

	return nil
}
func (s *sessionManagerImpl) Commit(session string) error {
	s.Lock()
	defer s.Unlock()

	impl, ok := s.sessions[session]

	if !ok {
		return identity.ErrSession
	}

	impl.committed = true

	return nil
}
