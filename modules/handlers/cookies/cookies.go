package cookies

import (
	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("very-secret-key-here"))

func init() {

	Store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to false for HTTP
		MaxAge:   3600,
	}
}

var UserStore map[string]string
