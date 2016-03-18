package gorest

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {

	mwOut := make([]int, 0, 20)
	mws := make([]Middleware, 0, 10)
	for i := 0; i < 20; i++ {
		mws = append(mws, func() Middleware {
			x := i
			return Middleware(func(h ContextHandlerFunc) ContextHandlerFunc {
				return ContextHandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
					mwOut = append(mwOut, x)
					h(c, w, r)
				})
			})
		}())
	}

	r := NewRouter()
	mwOut = mwOut[:0]
	r.Use(mws[0]).Use(mws[1]).Use(mws[2]).Use(mws[3])
	r.Get("/t1", func(c context.Context, w http.ResponseWriter, r *http.Request) { mwOut = append(mwOut, 101) })
	e := r.builtEntries[0]
	e.handlerFunc(nil, nil, nil)
	if !assert.Equal(t, e.method, "GET") {
		return
	}
	if !assert.Equal(t, e.path, "/t1") {
		return
	}
	if !assert.ObjectsAreEqual(mwOut, []int{0, 1, 2, 3, 101}) {
		log.Fatal("mwOut err %+v", mwOut)
	}
	//fmt.Printf("out %+v\n", mwOut)

	mwOut = mwOut[:0]
	g1 := r.NewGroup("/g1/g1")
	g1.Use(mws[5]).Use(mws[6])
	g1.Post("/t2", func(c context.Context, w http.ResponseWriter, r *http.Request) { mwOut = append(mwOut, 102) })
	e = r.builtEntries[1]
	e.handlerFunc(nil, nil, nil)
	if !assert.Equal(t, e.method, "POST") {
		return
	}
	if !assert.Equal(t, e.path, "/g1/g1/t2") {
		return
	}
	if !assert.ObjectsAreEqual(mwOut, []int{0, 1, 2, 3, 5, 6, 102}) {
		log.Fatal("mwOut err %+v", mwOut)
	}

	mwOut = mwOut[:0]
	g2 := r.NewGroup("/g2")
	g2.Use(mws[8])
	g2.Post("/t3/t3", func(c context.Context, w http.ResponseWriter, r *http.Request) { mwOut = append(mwOut, 103) })
	e = r.builtEntries[2]
	e.handlerFunc(nil, nil, nil)
	if !assert.Equal(t, e.method, "POST") {
		return
	}
	if !assert.Equal(t, e.path, "/g2/t3/t3") {
		return
	}
	if !assert.ObjectsAreEqual(mwOut, []int{0, 1, 2, 3, 8, 103}) {
		log.Fatal("mwOut err %+v", mwOut)
	}

	mwOut = mwOut[:0]
	g3 := g1.NewGroup("/g3/g3")
	g3.Use(mws[10]).Use(mws[11])
	g3.Delete("/t4", func(c context.Context, w http.ResponseWriter, r *http.Request) { mwOut = append(mwOut, 104) })
	e = r.builtEntries[3]
	e.handlerFunc(nil, nil, nil)
	if !assert.Equal(t, e.method, "DELETE") {
		return
	}
	if !assert.Equal(t, e.path, "/g1/g1/g3/g3/t4") {
		return
	}
	if !assert.ObjectsAreEqual(mwOut, []int{0, 1, 2, 3, 5, 6, 10, 11, 104}) {
		log.Fatal("mwOut err %+v", mwOut)
	}

	mwOut = mwOut[:0]
	g3.Put("/t5", func(c context.Context, w http.ResponseWriter, r *http.Request) { mwOut = append(mwOut, 107) })
	e = r.builtEntries[4]
	e.handlerFunc(nil, nil, nil)
	if !assert.Equal(t, e.method, "PUT") {
		return
	}
	if !assert.Equal(t, e.path, "/g1/g1/g3/g3/t5") {
		return
	}
	if !assert.ObjectsAreEqual(mwOut, []int{0, 1, 2, 3, 5, 6, 10, 11, 107}) {
		log.Fatal("mwOut err %+v", mwOut)
	}

	mwOut = mwOut[:0]
	r.Get("", func(c context.Context, w http.ResponseWriter, r *http.Request) { mwOut = append(mwOut, 106) })
	e = r.builtEntries[5]
	e.handlerFunc(nil, nil, nil)
	if !assert.Equal(t, e.method, "GET") {
		return
	}
	if !assert.Equal(t, e.path, "/") {
		return
	}
	if !assert.ObjectsAreEqual(mwOut, []int{0, 1, 2, 3, 106}) {
		log.Fatal("mwOut err %+v", mwOut)
	}

	mwOut = mwOut[:0]
	g2.Get("", func(c context.Context, w http.ResponseWriter, r *http.Request) { mwOut = append(mwOut, 107) })
	e = r.builtEntries[6]
	e.handlerFunc(nil, nil, nil)
	if !assert.Equal(t, e.method, "GET") {
		return
	}
	if !assert.Equal(t, e.path, "/g2") {
		return
	}
	if !assert.ObjectsAreEqual(mwOut, []int{0, 1, 2, 3, 8, 107}) {
		log.Fatal("mwOut err %+v", mwOut)
	}

	assert.Panics(t, func() {
		g3.Delete("/t6", nil)
	})
	assert.Panics(t, func() {
		g3.Delete("t7", func(c context.Context, w http.ResponseWriter, r *http.Request) {})
	})
	assert.Panics(t, func() {
		g3.Delete("/t4", func(c context.Context, w http.ResponseWriter, r *http.Request) {})
	})
	assert.Panics(t, func() { r.NewGroup("/g1/g1") })

}
