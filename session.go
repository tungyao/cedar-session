package cedar_session

import (
	"crypto/sha1"
	"fmt"
	"github.com/tungyao/cedar"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var X map[string][]*SX

func init() {
	X = make(map[string][]*SX)
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
	x.Group("/a", func(groups *TheGroup) {
		groups.Get("/b", func(w http.ResponseWriter, r *http.Request, s Session) {
			w.Write([]byte("hello"))
		}, nil)
	})
	http.ListenAndServe(":80", x.Handler)
}
func NewSession(hp *cedar.Trie) *SessionX {
	log.Println("Session : starting")
	s := &SessionX{
		RWMutex: &sync.RWMutex{},
		Mutex:   &sync.Mutex{},
		Handler: hp,
		Self:    newId(),
	}
	return s
}

// struct
type SessionX struct {
	*sync.RWMutex
	Handler *cedar.Trie
	*sync.Mutex
	Self []byte
}
type TheGroup struct {
	cedar.Groups
	S *SessionX
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
func (si *SessionX) Get(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), hal http.Handler) {
	si.Handler.Get(path, func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(si.CreateUUID([]byte(request.RemoteAddr)))
			http.SetCookie(writer, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true,
				Expires:  time.Now().Add(1 * time.Hour),
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
	}, hal)
}
func (si *SessionX) Post(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), hal http.Handler) {
	si.Handler.Post(path, func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(si.CreateUUID([]byte(request.RemoteAddr)))
			http.SetCookie(writer, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true,
				Expires:  time.Now().Add(1 * time.Hour),
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
	}, hal)
}
func (si *SessionX) Put(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), hal http.Handler) {
	si.Handler.Put(path, func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(si.CreateUUID([]byte(request.RemoteAddr)))
			http.SetCookie(writer, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true,
				Expires:  time.Now().Add(1 * time.Hour),
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
	}, hal)
}
func (si *SessionX) Delete(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), hal http.Handler) {
	si.Handler.Delete(path, func(writer http.ResponseWriter, request *http.Request) {
		c, err := request.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(si.CreateUUID([]byte(request.RemoteAddr)))
			http.SetCookie(writer, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true,
				Expires:  time.Now().Add(1 * time.Hour),
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
	}, hal)
}
func (si *SessionX) Group(path string, fn func(groups *TheGroup)) {
	g := new(TheGroup)
	g.Tree = si.Handler
	g.Path = path
	g.S = si
	fn(g)
}

// group function
func (t *TheGroup) Get(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), handler http.Handler) {
	t.Groups.Get(path, func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(t.S.CreateUUID([]byte(r.RemoteAddr)))
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true,
				Expires:  time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(w, r, Session{
				RWMutex: sync.RWMutex{},
				Cookie:  c.Value,
			})
		} else {
			fn(w, r, Session{
				RWMutex: sync.RWMutex{},
			})
		}
	}, handler)
}
func (t *TheGroup) Post(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), handler http.Handler) {
	t.Groups.Post(path, func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(t.S.CreateUUID([]byte(r.RemoteAddr)))
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true,
				Expires:  time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(w, r, Session{
				RWMutex: sync.RWMutex{},
				Cookie:  c.Value,
			})
		} else {
			fn(w, r, Session{
				RWMutex: sync.RWMutex{},
			})
		}
	}, handler)
}
func (t *TheGroup) Put(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), handler http.Handler) {
	t.Groups.Put(path, func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(t.S.CreateUUID([]byte(r.RemoteAddr)))
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true,
				Expires:  time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(w, r, Session{
				RWMutex: sync.RWMutex{},
				Cookie:  c.Value,
			})
		} else {
			fn(w, r, Session{
				RWMutex: sync.RWMutex{},
			})
		}
	}, handler)
}
func (t *TheGroup) Delete(path string, fn func(w http.ResponseWriter, r *http.Request, s Session), handler http.Handler) {
	t.Groups.Get(path, func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err == http.ErrNoCookie {
			x := Sha1(t.S.CreateUUID([]byte(r.RemoteAddr)))
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    string(x),
				HttpOnly: true,
				Expires:  time.Now().Add(1 * time.Hour),
			})
		}
		if c != nil {
			fn(w, r, Session{
				RWMutex: sync.RWMutex{},
				Cookie:  c.Value,
			})
		} else {
			fn(w, r, Session{
				RWMutex: sync.RWMutex{},
			})
		}
	}, handler)
}
func (t *TheGroup) Group(path string, fn func(groups *TheGroup)) {
	g := new(TheGroup)
	g.Path = t.Path + path
	g.Tree = t.Tree
	fn(g)
}

//func (mux *SessionX) Delete(path string, handlerFunc http.HandlerFunc, handler http.Handler) {
//	mux.tree.Delete(mux.path+path, handlerFunc, handler)
//}
// UUID 64 bit
// 8-4-4-12 16hex string
func (si *SessionX) CreateUUID(xtr []byte) []byte {
	str := fmt.Sprintf("%x", xtr)
	strLow := ComplementHex(str[:(len(str)-1)/3], 8)
	strMid := ComplementHex(str[(len(str)-1)/3:(len(str)-1)*2/3], 4)
	si.Mutex.Lock()
	defer si.Mutex.Unlock()
	<-time.After(1 * time.Nanosecond)
	ti := time.Now().UnixNano()
	return []byte(fmt.Sprintf("%s-%x-%s-%s", strLow, ti, strMid, si.Self))
}

// session function
func (sn Session) Set(key string, body interface{}) {
	X[sn.Cookie] = append(X[sn.Cookie], &SX{
		Key:  key,
		Body: body,
	})
}
func (sn Session) Get(key string) interface{} {
	x := X[sn.Cookie]
	for _, v := range x {
		if v.Key == key {
			return v.Body
		}
	}
	return nil
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
	log.Printf("Session : Create Session Id => %s", da)
	return da
}
