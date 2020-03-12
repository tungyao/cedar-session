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
	x := cedar_session.NewSession(r)
	x.Get("/set", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
		s.Set("hello", "world"+r.RemoteAddr)
		w.Write([]byte("hello world"))
	}, nil)
	x.Get("/get", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
		fmt.Fprintf(w, "%s", s.Get("hello"))
	}, nil)
	x.Group("/a", func(groups *cedar_session.TheGroup) {
		groups.Get("/b", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
			w.Write([]byte("hello"))
		}, nil)
		groups.Group("/c", func(groups *cedar_session.TheGroup) {
			groups.Get("/d", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
				w.Write([]byte("bye bye"))
			}, nil)
		})
	})
	http.ListenAndServe(":80", x.Handler)
}
