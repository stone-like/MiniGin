package minigin

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedGroupValue = ""

func fakeGroupHandler(val string) HandlerFunc {
	return func(c *Context) {
		expectedGroupValue += val + " "
	}
}

func TestRouterGroup(t *testing.T) {

	e := New()

	rg1 := e.Group("/test")

	rg1.Use(fakeGroupHandler("middleware1"), fakeGroupHandler("middleware2"))

	rg1.GET("/foo", fakeGroupHandler("foo"))

	rg2 := rg1.Group("/test2")

	rg2.Use(fakeGroupHandler("middleware3"))

	rg2.POST("/bar", fakeGroupHandler("bar"))

	cases := []struct {
		name       string
		path       string
		method     string
		expected   string
		err        error
		handlerLen int
	}{
		{name: "fooOK", method: http.MethodGet, path: "/test/foo", expected: "middleware1 middleware2 foo ", err: nil, handlerLen: 3},
		{name: "fooNG", method: http.MethodPost, path: "/test/foo", expected: "", err: methodNotRegisteredError(http.MethodPost), handlerLen: 0},
		{name: "barOK", method: http.MethodPost, path: "/test/test2/bar", expected: "middleware1 middleware2 middleware3 bar ", err: nil, handlerLen: 4},
		{name: "barNG", method: http.MethodPost, path: "/test2/bar", expected: "", err: pathNotRegisteredError("/test2/bar"), handlerLen: 0},
		{name: "barNG2", method: http.MethodPost, path: "/baz", expected: "", err: pathNotRegisteredError("/baz"), handlerLen: 0},
	}

	for _, c := range cases {

		t.Run(c.name, func(t *testing.T) {
			t.Helper()
			res, err := e.router.Search(c.method, c.path)

			if err == nil {

				assert.Equal(t, c.handlerLen, len(res.handlers))

				for _, h := range res.handlers {
					h(nil)
				}
				assert.Equal(t, expectedGroupValue, c.expected)
				//リセット
				expectedGroupValue = ""
			} else {
				assert.Equal(t, err.Error(), c.err.Error())
			}
		})

	}

}
