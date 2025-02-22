package cookies

import (
	"sync"

	"github.com/gorilla/sessions"
)

var UserStore map[string]string

var Session *session

type session struct {
	subs       map[chan<- []byte]struct{}
	subs_mutex *sync.Mutex
	cookies    *sessions.CookieStore
}
type Active interface {
	Subscribe(c chan<- []byte) (UnSubscribe, error)
}

type UnSubscribe func() error

func New(session_key []byte) *session {
	session := session{
		subs:       map[chan<- []byte]struct{}{},
		subs_mutex: &sync.Mutex{},
		cookies:    sessions.NewCookieStore(session_key),
	}
	session.cookies.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to false for HTTP
		MaxAge:   3600,
	}
	return &session
}

type Notifier interface {
	Notify(b []byte) error
}

func (s *session) Subscribe(c chan<- []byte) (UnSubscribe, error) {
	s.subs_mutex.Lock()
	s.subs[c] = struct{}{}
	s.subs_mutex.Unlock()
	unsubscribeFn := func() error {
		s.subs_mutex.Lock()
		delete(s.subs, c)
		close(c)
		s.subs_mutex.Unlock()

		return nil
	}

	return unsubscribeFn, nil
}

func (s *session) Notify(b []byte) error {
	s.subs_mutex.Lock()
	defer s.subs_mutex.Unlock()
	for c := range s.subs {
		if len(c) == cap(c) {
			continue
		}
		c <- b
	}
	return nil
}

func (s *session) Store() *sessions.CookieStore {
	return s.cookies
}
