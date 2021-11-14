package minigin

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
)

var expectedValue = ""

//trieからは handler - func(c Context)を返すようにする
func fakeHandler(val string) HandlerFunc {
	return func(c *Context) {
		expectedValue = val
	}
}

func TestSearch(t *testing.T) {
	trie := NewTrie()
	trie.Insert(http.MethodGet, "/", HandlerChain{fakeHandler("root")})
	trie.Insert(http.MethodGet, "/foo", HandlerChain{fakeHandler("foo")})
	trie.Insert(http.MethodGet, "/bar", HandlerChain{fakeHandler("bar")})
	trie.Insert(http.MethodPost, "/foo/bar", HandlerChain{fakeHandler("fooBar")})

	cases := []struct {
		name     string
		path     string
		method   string
		expected string
		err      error
	}{
		{name: "rootOK", method: http.MethodGet, path: "/", expected: "root", err: nil},
		{name: "fooOK", method: http.MethodGet, path: "/foo", expected: "foo", err: nil},
		{name: "fooNG", method: http.MethodPost, path: "/foo", expected: "", err: methodNotRegisteredError(http.MethodPost)},
		{name: "barOK", method: http.MethodGet, path: "/bar", expected: "bar", err: nil},
		{name: "bazNG", method: http.MethodGet, path: "/baz", expected: "", err: pathNotRegisteredError("/baz")},
		{name: "barNG", method: http.MethodGet, path: "/bar/dummy", expected: "", err: pathNotRegisteredError("/bar/dummy")},
		{name: "fooBarOK", method: http.MethodPost, path: "/foo/bar", expected: "fooBar", err: nil},
	}

	for _, c := range cases {

		t.Run(c.name, func(t *testing.T) {
			t.Helper()
			res, err := trie.Search(c.method, c.path)

			if err == nil {

				assert.Equal(t, 1, len(res.handlers))
				res.handlers[0](nil) //hを起動して、expectedValueに値を入れる
				assert.Equal(t, expectedValue, c.expected)

			} else {
				assert.Equal(t, err.Error(), c.err.Error())
			}
		})

	}
}

func TestSearchWild(t *testing.T) {
	trie := NewTrie()
	trie.Insert(http.MethodGet, "/", HandlerChain{fakeHandler("root")})
	trie.Insert(http.MethodGet, "/foo/:name", HandlerChain{fakeHandler("foo")})
	trie.Insert(http.MethodGet, "/bar/:id/:meta", HandlerChain{fakeHandler("bar")})

	cases := []struct {
		name     string
		path     string
		method   string
		expected string
		params   Params
		err      error
	}{
		{name: "rootOK", method: http.MethodGet, path: "/", expected: "root", params: []Param{}, err: nil},
		{name: "fooOK", method: http.MethodGet, path: "/foo/myname", expected: "foo", params: []Param{{Key: "name", Value: "myname"}}, err: nil},
		{name: "fooNG", method: http.MethodGet, path: "/foo", expected: "", err: pathNotRegisteredError("/foo")},
		{name: "barOK", method: http.MethodGet, path: "/bar/56/someMeta", expected: "bar", params: []Param{{Key: "id", Value: "56"}, {Key: "meta", Value: "someMeta"}}, err: nil},
		{name: "bazNG", method: http.MethodGet, path: "/baz", expected: "", err: pathNotRegisteredError("/baz")},
		{name: "barNG", method: http.MethodGet, path: "/bar/dummy", expected: "", params: []Param{}, err: pathNotRegisteredError("/bar/dummy")},
		{name: "barNG2", method: http.MethodGet, path: "/bar/dummy/aaa/bbb", expected: "", params: []Param{}, err: pathNotRegisteredError("/bar/dummy/aaa/bbb")},
	}

	for _, c := range cases {

		t.Run(c.name, func(t *testing.T) {
			t.Helper()
			res, err := trie.Search(c.method, c.path)

			if err == nil {
				assert.Equal(t, 1, len(res.handlers))
				res.handlers[0](nil) //hを起動して、expectedValueに値を入れる
				assert.Equal(t, expectedValue, c.expected)

				if diff := cmp.Diff(res.params, c.params); diff != "" {
					t.Errorf("diff is %v\n", diff)
				}

			} else {
				assert.Equal(t, err.Error(), c.err.Error())
			}
		})

	}
}
