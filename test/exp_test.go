package test

import (
	"fmt"
	"github.com/tungyao/cedar"
	"net/http"
	"testing"
)
import "../../cedar-session"

func TestMainX(t *testing.T) {
	r := cedar.NewRouter()
	x := cedar_session.NewSession(r, cedar_session.SpruceLocal)
	x.Middleware("test", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) bool {
		fmt.Println("12312")
		http.Redirect(w, r, "/a/b", 302)
		return false
	})
	x.Get("/", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
		w.Write([]byte("hellox"))
	}, nil, "test")
	x.Get("/set", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
		s.Set("hello", []byte("world"+r.RemoteAddr))
		w.Write([]byte("hello world"))
	}, nil)
	x.Get("/get", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
		fmt.Fprintf(w, "%s", s.Get("hello"))
	}, nil)
	x.Group("/a", func(groups *cedar_session.Group) {
		groups.Get("/b", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
			w.Write([]byte("hello"))
		}, nil)
		groups.Group("/c", func(groups *cedar_session.Group) {
			groups.Get("/d", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
				w.Write([]byte("bye bye"))
			}, nil)
		})
	})
	http.ListenAndServe(":82", x.Handler)
}
