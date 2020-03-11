package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/tungyao/cedar"
	"github.com/tungyao/spruce"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	r := cedar.NewRouter()
	x := NewSession(r)
	x.Get("/", func(writer http.ResponseWriter, request *http.Request, session *Session) {
		session.SetSession()
	}, nil)
	http.ListenAndServe(":80", x)

}
func NewSession(hp *cedar.Trie) *Session {
	log.Println("Session : starting")
	h := spruce.CreateHash(1024)
	s := &Session{
		RWMutex: &sync.RWMutex{},
		Mutex:   &sync.Mutex{},
		tree:    hp,
		hash:    h,
		Self:    newId(),
	}
	return s
}

type SessionSelf struct {
	Id   string
	Time int64
}
type Session struct {
	*sync.RWMutex
	tree *cedar.Trie
	hash *spruce.Hash
	*sync.Mutex
	Self []byte
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
func (si *Session) GetSession(uuid string) {

}
func (si *Session) Get(path string, fn func(w http.ResponseWriter, r *http.Request, s *Session), hal http.Handler) {
	//x := Sha1(si.CreateUUID())
	//c1 := http.Cookie{
	//	Name:     "session",
	//	Value:    string(x),
	//	HttpOnly: true,
	//}
	//w.Header().Set("Set-Cookie", c1.String())
	//w.Header().Add("Set-Cookie", c1.String())
	si.tree.Get(path, func(writer http.ResponseWriter, request *http.Request) {

	}, hal)
}

//  you don't need to input a value of token when put a session key in first time
// it's second times , you must input it
// output sha-1
func (si *Session) SetSession(w http.ResponseWriter, key []byte, body interface{}) []byte {
	if key == nil {

		return x
	}
	return nil
}
func Sha1(b []byte) []byte {
	h := sha1.New()
	h.Write(b)
	return []byte(fmt.Sprintf("%x", h.Sum(nil)))
}

// UUID 64 bit
// 8-4-4-12 16hex string
func (si *Session) CreateUUID(xtr []byte) []byte {
	str := fmt.Sprintf("%x", xtr)
	strLow := ComplementHex(str[:(len(str)-1)/3], 8)
	strMid := ComplementHex(str[(len(str)-1)/3:(len(str)-1)*2/3], 4)
	si.Mutex.Lock()
	defer si.Mutex.Unlock()
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
