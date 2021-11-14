package minigin

// type header struct {
// 	Key   string
// 	Value string
// }

// func performRequest(r http.Handler, method, path string, headers ...header) *httptest.ResponseRecorder {
// 	req := httptest.NewRequest(method, path, nil)
// 	for _, h := range headers {
// 		req.Header.Add(h.Key, h.Value)
// 	}
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)
// 	return w
// }

// func TestEngineHandleContext(t *testing.T) {
// 	e := New()
// 	e.GET("/", func(c *Context) {
// 		c.Writer.WriteHeader(203)
// 		c.Writer.Write([]byte("hello"))
// 	})

// 	w := performRequest(e, "GET", "/")

// 	assert.Equal(t, "hello", w.Body.String())
// 	assert.Equal(t, 203, w.Code)

// }
