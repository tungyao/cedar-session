# session component for cedar
## usage
* step1 `r:=cedar.NewRouter()`
* step2 `x := cedar_session.NewSession(r)`
* step3
```go
// you can use all router func
x.Get("/",http.HandlerFunc,http.Handler)
// and session func
x.SetSession()
x.GetSession()
```

## one more thing
* x.CreateUUID([]byte) ,transfer ip or other parameter to the first
### get byte => `********-****************-****-****`
### get byte =>  8-16-4-4
