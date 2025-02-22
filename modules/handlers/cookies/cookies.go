package cookies

import (
	"sync"

	"github.com/gorilla/sessions"
)

var UserStore map[string]string

var Session *session

type session struct {
	subs_mutex *sync.Mutex
	cookies    *sessions.CookieStore
}

func New(session_key []byte) *session {
	session := session{
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

func (s *session) Store() *sessions.CookieStore {
	return s.cookies
}
