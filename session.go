package cedar_session

import (
	"crypto/sha1"
	"fmt"
	ap "git.yaop.ink/tungyao/awesome-pool"
	"github.com/tungyao/spruce"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/tungyao/cedar"
)

var (
	KV *ap.Pool
	X  *spruce.Hash
	OP int = -1
)

const (
	LOCAL = iota
	SpruceLocal
	SPRUCE
)

func main() {
	//r := cedar.NewRouter()
	//x := NewSession(r, LOCAL)
	//x.Get("/set", func(w http.ResponseWriter, r *http.Request, s Session) {
	//	s.Set("hello", "world"+r.RemoteAddr)
	//	w.Write([]byte("hello world"))
	//}, nil)
	//x.Get("/get", func(w http.ResponseWriter, r *http.Request, s Session) {
	//	fmt.Fprintf(w, "%s", s.Get("hello"))
	//}, nil)
	//x.Group("/a", func(groups *Group) {
	//	groups.Get("/b", func(w http.ResponseWriter, r *http.Request, s Session) {
	//		w.Write([]byte("hello"))
	//	}, nil)
	//})
	//http.ListenAndServe(":80", x.Handler)
}
func NewSession(hp *cedar.Trie, types int, args ...interface{}) *sessionx {
	s := &sessionx{
		Handler: hp,
		Self:    newId(),
		op:      types,
	}
	OP = types
	switch types {
	case LOCAL:
		X = spruce.CreateHash(4096)
	case SPRUCE:
	case SpruceLocal:
		KV, _ = ap.NewPool(args[0].(int), args[1].(string))
	}
	return s
}

// struct
type sessionx struct {
	Handler *cedar.Trie
	sync.Mutex
	Self []byte
	op   int
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
func (sn Session) Set(key string, body []byte) {
	switch OP {
	case LOCAL:
		X.Set([]byte(sn.Cookie+key), body, 3600)
	case SPRUCE:
	case SpruceLocal:
		kvSet([]byte(sn.Cookie+key), body, 3600)
	}
}
func (sn Session) Get(key string) interface{} {
	switch OP {
	case LOCAL:
		return X.Get([]byte(sn.Cookie + key))
	case SPRUCE:
	case SpruceLocal:
		return kvGet([]byte(sn.Cookie + key))
	}
	return []byte("")
}
func (sn Session) Flush(key string) interface{} {
	switch OP {
	case LOCAL:
		return X.Delete([]byte(sn.Cookie + key))
	case SPRUCE:
	case SpruceLocal:
		return kvDelete([]byte(sn.Cookie + key))
	}
	return []byte("")

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
func kvSet(key, body []byte, exp int) []byte {
	KV.Get().Write(spruce.EntrySet(key, body, exp))
	return KV.Get().Read()
}
func kvGet(key []byte) []byte {
	KV.Get().Write(spruce.EntryGet(key))
	return KV.Get().Read()
}
func kvDelete(key []byte) []byte {
	KV.Get().Write(spruce.EntryDelete(key))
	return KV.Get().Read()
}
