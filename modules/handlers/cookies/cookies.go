package cookies

import (
	"GoRecordurbate/modules/file"
	"net/http"
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

func (s *session) IsLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	session, err := s.Store().Get(r, "session")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return false
	}
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {

		apiKey := r.URL.Query().Get("api_key")

		var secrets file.API_secrets
		file.ReadJson(file.API_keys_file, &secrets)
		for _, key := range secrets.Keys {
			if VerifyAPIKey(key.Key, apiKey) {
				return false
			}
		}
		return false
		//http.HandleFunc("/login", handlers.GetLogin)

	}
	return true
}
