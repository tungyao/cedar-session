package cedar_session

import (
	"crypto/sha1"
	"fmt"
	"github.com/tungyao/spruce"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/tungyao/cedar"
)

var X *spruce.Hash

func init() {
	X = spruce.CreateHash(1024)
}
func main() {
	r := cedar.NewRouter()
	x := NewSession(r)
	x.Get("/set", func(w http.ResponseWriter, r *http.Request, s Session) {
		s.Set("hello", "world"+r.RemoteAddr)
		w.Write([]byte("hello world"))
	}, nil)
	x.Get("/get", func(w http.ResponseWriter, r *http.Request, s Session) {
		fmt.Fprintf(w, "%s", s.Get("hello"))
	}, nil)
	x.Group("/a", func(groups *Group) {
		groups.Get("/b", func(w http.ResponseWriter, r *http.Request, s Session) {
			w.Write([]byte("hello"))
		}, nil)
	})
	http.ListenAndServe(":80", x.Handler)
}
func NewSession(hp *cedar.Trie) *sessionx {
	s := &sessionx{
		Handler: hp,
		Self:    newId(),
	}
	return s
}

// struct
type sessionx struct {
	Handler *cedar.Trie
	sync.Mutex
	Self []byte
}
type Group struct {
	gG cedar.Groups
	S  *sessionx
}
type Session struct {
	sync.RWMutex
	Cookie string
}
type SX struct {
	Key  string
	Body interface{}
}

//
func (si *sessionx) Get(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), hal http.Handler, middleware ...string) {
	si.Handler.Get(path, func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(si.CreateUUID([]byte(request.RemoteAddr)))
			http.SetCookie(writer, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true, Secure: true,
				Expires: time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(writer, request, Session{
				RWMutex: sync.RWMutex{},
				Cookie:  c.Value,
			})
		} else {
			fn(writer, request, Session{
				RWMutex: sync.RWMutex{},
			})
		}
	}, hal, middleware...)
}
func (si *sessionx) Post(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), hal http.Handler, middleware ...string) {
	si.Handler.Post(path, func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(si.CreateUUID([]byte(request.RemoteAddr)))
			http.SetCookie(writer, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true, Secure: true,
				Expires: time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(writer, request, Session{
				Cookie: c.Value,
			})
		} else {
			fn(writer, request, Session{})
		}
	}, hal, middleware...)
}
func (si *sessionx) Put(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), hal http.Handler, middleware ...string) {
	si.Handler.Put(path, func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(si.CreateUUID([]byte(request.RemoteAddr)))
			http.SetCookie(writer, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true, Secure: true,
				Expires: time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(writer, request, Session{
				Cookie: c.Value,
			})
		} else {
			fn(writer, request, Session{})
		}
	}, hal, middleware...)
}
func (si *sessionx) Delete(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), hal http.Handler, middleware ...string) {
	si.Handler.Delete(path, func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(si.CreateUUID([]byte(request.RemoteAddr)))
			http.SetCookie(writer, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true, Secure: true,
				Expires: time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(writer, request, Session{
				Cookie: c.Value,
			})
		} else {
			fn(writer, request, Session{})
		}
	}, hal, middleware...)
}
func (si *sessionx) Group(path string, fn func(groups *Group)) {
	g := new(Group)
	g.gG.Tree = si.Handler
	g.gG.Path = path
	g.S = si
	fn(g)
}
func (si *sessionx) Dynamic(ymlPath string) {
	si.Handler.Dynamic(ymlPath)
}
func (si *sessionx) Middleware(name string, fn func(w http.ResponseWriter, r *http.Request, s Session) bool) {
	si.Handler.Middle(name, func(w http.ResponseWriter, r *http.Request) bool {
		return fn(w, r, Session{
			RWMutex: sync.RWMutex{},
		})
	})
}

// group function
func (t *Group) Get(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), handler http.Handler, middleware ...string) {
	t.gG.Get(path, func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(t.S.CreateUUID([]byte(r.RemoteAddr)))
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true, Secure: true,
				Expires: time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(w, r, Session{
				Cookie: c.Value,
			})
		} else {
			fn(w, r, Session{})
		}
	}, handler, middleware...)
}
func (t *Group) Post(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), handler http.Handler, middleware ...string) {
	t.gG.Post(path, func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(t.S.CreateUUID([]byte(r.RemoteAddr)))
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true, Secure: true,
				Expires: time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(w, r, Session{
				Cookie: c.Value,
			})
		} else {
			fn(w, r, Session{})
		}
	}, handler, middleware...)
}
func (t *Group) Put(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), handler http.Handler, middleware ...string) {
	t.gG.Put(path, func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(t.S.CreateUUID([]byte(r.RemoteAddr)))
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true, Secure: true,
				Expires: time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(w, r, Session{
				Cookie: c.Value,
			})
		} else {
			fn(w, r, Session{})
		}
	}, handler, middleware...)
}
func (t *Group) Delete(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), handler http.Handler, middleware ...string) {
	t.gG.Delete(path, func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(t.S.CreateUUID([]byte(r.RemoteAddr)))
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true,
				Secure:   true,
				Expires:  time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(w, r, Session{
				Cookie: c.Value,
			})
		} else {
			fn(w, r, Session{})
		}
	}, handler, middleware...)
}
func (t *Group) Group(path string, fn func(groups *Group)) {
	g := new(Group)
	g.gG.Path = t.gG.Path + path
	g.gG.Tree = t.gG.Tree
	g.S = t.S
	fn(g)
}
func (t *Group) Middleware(name string, fn func(w http.ResponseWriter, r *http.Request, s Session) bool) {
	t.gG.Middleware(name, func(w http.ResponseWriter, r *http.Request) bool {
		return fn(w, r, Session{
			RWMutex: sync.RWMutex{},
		})
	})
}

// func (mux *SessionX) Delete(path string, handlerFunc http.HandlerFunc, handler http.Handler) {
//	mux.tree.Delete(mux.path+path, handlerFunc, handler)
// }
// UUID 64 bit
// 8-4-4-12 16hex string
func (si *sessionx) CreateUUID(xtr []byte) []byte {
	str := fmt.Sprintf("%x", xtr)
	strLow := ComplementHex(str[:(len(str)-1)/3], 8)
	strMid := ComplementHex(str[(len(str)-1)/3:(len(str)-1)*2/3], 4)
	si.Lock()
	defer si.Unlock()
	<-time.After(1 * time.Nanosecond)
	ti := time.Now().UnixNano()
	return []byte(fmt.Sprintf("%s-%x-%s-%s", strLow, ti, strMid, si.Self))
}

func ComplementHex(s string, x int) string {
	if len(s) == x {
		return s
	}
	if len(s) < x {
		for i := 0; i < x-len(s); i++ {
			s += "0"
		}
	}
	if len(s) > x {
		return s[:x]
	}
	return s
}

// session function
func (sn Session) Set(key string, body interface{}) {
	X.Set([]byte(key), body, 3600)
	//X[sn.Cookie] = append(X[sn.Cookie], &SX{
	//	Key:  key,
	//	Body: body,
	//})
}
func (sn Session) Get(key string) interface{} {
	//x := X[sn.Cookie]
	//for _, v := range x {
	//	if v.Key == key {
	//		return v.Body
	//	}
	//}
	//return nil
	return X.Get([]byte(key))
}

// other function
func Sha1(b []byte) []byte {
	h := sha1.New()
	h.Write(b)
	return []byte(fmt.Sprintf("%x", h.Sum(nil)))
}
func newId() []byte {
	d := "abcdef012345689"
	da := make([]byte, 4)
	for i := 0; i < 4; i++ {
		<-time.After(time.Nanosecond * 10)
		da[i] = d[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(15)]
	}
	return da
}
