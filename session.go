package main

import (
	"fmt"
	"github.com/tungyao/cedar"
	"github.com/tungyao/spruce"
	"net/http"
	"sync"
	"time"
)

func main() {
	r := cedar.NewRouter()
	x := NewSession(r)
	x.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write(x.SetSession(nil, ""))
	}, nil)
	http.ListenAndServe(":80", x)
}
func NewSession(hp *cedar.Trie) *Session {
	h := spruce.CreateHash(1024)
	s := &Session{
		RWMutex: &sync.RWMutex{},
		Mutex:   &sync.Mutex{},
		Trie:    hp,
		hash:    h,
	}
	return s
}

type Session struct {
	*sync.RWMutex
	*cedar.Trie
	hash *spruce.Hash
	*sync.Mutex
}

func (si *Session) GetSession(uuid string) {

}

//  you don't need to input a value of token when put a session key in first time
// it's second times , you must input it
func (si *Session) SetSession(token []byte, stt interface{}) []byte {
	if token == nil {
		<-time.After(1 * time.Nanosecond)
		return si.CreateUUID(token, time.Now().UnixNano())
	}
	return nil
}

// UUID 64 bit
// 8-4-4-12 16hex string
func (si *Session) CreateUUID(name []byte, tms int64) []byte {
	fmt.Printf("%x\n", tms)
	return []byte(fmt.Sprintf("%x\n", tms))
}
