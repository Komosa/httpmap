package httpmap

import (
	"net/http"
	"strings"
)

var Colon = byte(':')
var Slash = byte('/')
var stringSlash = string(Slash)
var slashColon = stringSlash + string(Colon)

// Object implementing Handler interface can be used as
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, Named)
}

type HandlerFunc func(http.ResponseWriter, *http.Request, Named)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, n Named) {
	f(w, r, n)
}

type handlerWrapper struct {
	name    name2idx
	handler Handler
}

// muxMethod stores routes (URLs) and corresponding handlers for given HTTP method (GET, POST, ...)
// should be fully constructed before any access (query)
// later can be accessed concurrently in safe way
type muxMethod map[string]handlerWrapper
type idx2name []string
type name2idx map[string]int

// Mux stores routes (URLs) and corresponding handlers
type Mux struct {
	mget, mpost, mdelete, mput muxMethod
}

func New() Mux {
	return Mux{make(muxMethod), make(muxMethod), make(muxMethod), make(muxMethod)}
}

func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, p := parseRequest(r.RequestURI)
	switch r.Method {
	case "GET":
		b := Named{idx: p, name: m.mget[s].name}
		m.mget[s].handler.ServeHTTP(w, r, b)
	default:
		panic(nil)
	}
}

func (m Mux) Get(url string, handler Handler) {
	s, n := parseHandler(url)
	m.mget[s] = handlerWrapper{
		name:    n,
		handler: handler,
	}
}

type Named struct {
	idx  idx2name
	name name2idx
}

func (b Named) Param(s string) string {
	i, ok := b.name[s]
	if !ok {
		return ""
	}
	return b.idx[i]
}

func parseRequest(url string) (compressed string, idx idx2name) {
	if len(url) > 0 && url[0] == Slash {
		url = url[1:]
	}

	for _, c := range strings.Split(url, stringSlash) {
		if len(c) == 0 {
			panic(nil)
		}
		if c[0] == Colon {
			if len(c) == 1 {
				panic(nil)
			} // TODO: introduce error
			compressed += slashColon
			idx = append(idx, c[1:])
		} else {
			compressed += stringSlash + c
			idx = append(idx, "")
		}
	}
	return
}

func parseHandler(url string) (compressed string, name name2idx) {
	if len(url) > 0 && url[0] == Slash {
		url = url[1:]
	}
	name = make(name2idx)
	for i, c := range strings.Split(url, stringSlash) {
		if len(c) == 0 {
			panic(nil)
		}
		if c[0] == Colon {
			if len(c) == 1 {
				panic(nil)
			} // TODO: introduce error
			compressed += slashColon
			name[c[1:]] = i
		} else {
			compressed += stringSlash + c
		}
	}
	return
}
