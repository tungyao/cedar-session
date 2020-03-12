# session component for cedar
## usage
* step1 `r:=cedar.NewRouter()`
* step2 `x := cedar_session.NewSession(r)`
* step3
### Session method ,it must use in the http method
**Set(key string ,body interface{})**

**Get(key string) interface{}**

```go
x.Get("/set", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
	s.Set("hello", "world"+r.RemoteAddr) // set session
}, nil)
x.Get("/get", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session) {
	s.Get("hello") // get session
}, nil)
```
### Http method
**it like cedar router , i've only changed a few code**
```
// before
r.Get("/",func(w http.ResponseWriter, r *http.Request,nil)
// after
x.Get("/", func(w http.ResponseWriter, r *http.Request, s cedar_session.Session), nil)
```
### group
```go
// before
r.Group("/a", func(groups *cedar.Groups) {}
// after
x.Group("/a", func(groups *cedar_session.TheGroup) {}
```
## one more thing
* x.CreateUUID([]byte) []byte ,transfer ip or other parameter to the first
### get byte => `********-****************-****-****`
### get byte =>  8-16-4-4
