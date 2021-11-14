package minigin

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createContextWithEngine() *Context {
	e := New()
	return e.allocateContext()
}

func TestGetQuery(t *testing.T) {
	c := createContextWithEngine()
	req, _ := http.NewRequest("GET", "/test?foo=bar&name=test&id=", nil)
	c.Request = req
	// c.Request, _ = http.NewRequest("GET", "http://example.com/?foo=bar&page=10&id=", nil)

	value, ok := c.GetQuery("foo")
	assert.Equal(t, "bar", value)
	assert.True(t, ok)

	value, ok = c.GetQuery("name")
	assert.Equal(t, "test", value)
	assert.True(t, ok)

	value, ok = c.GetQuery("id")
	assert.Equal(t, "", value)
	assert.True(t, ok)

	value, ok = c.GetQuery("bad")
	assert.Equal(t, "", value)
	assert.False(t, ok)

	value, ok = c.GetPostForm("foo")
	assert.Equal(t, "", value)
	assert.False(t, ok)

}

func TestGetForm(t *testing.T) {
	body := bytes.NewBufferString("foo=bar&name=test&id=")
	c := createContextWithEngine()
	req, _ := http.NewRequest("POST", "/test?foo=dup", body)
	req.Header.Add("Content-Type", MIMEPOSTForm)
	//req.ParseMultipartFormは適切なContent-Typeを設定しないとParseされない,jsonとかはだめ
	//jsonの時は直接Body使うので必要ないということかな？
	//"application/x-www-form-urlencoded"じゃなきゃだめっぽい
	//net/http.ProtocolError {ErrorString: "request Content-Type isn't multipart/form-data"}

	c.Request = req

	value, ok := c.GetPostForm("foo")
	assert.Equal(t, "bar", value)
	assert.True(t, ok)

	value, ok = c.GetPostForm("name")
	assert.Equal(t, "test", value)
	assert.True(t, ok)

	value, ok = c.GetPostForm("id")
	assert.Equal(t, "", value)
	assert.True(t, ok)

	value, ok = c.GetPostForm("bad")
	assert.Equal(t, "", value)
	assert.False(t, ok)

	value, ok = c.GetQuery("name")
	assert.Equal(t, "", value)
	assert.False(t, ok)

	value, ok = c.GetQuery("foo")
	assert.Equal(t, "dup", value)
	assert.True(t, ok)
}

func TestGetContentType(t *testing.T) {
	body := bytes.NewBufferString("foo=bar&name=test&id=")
	c := createContextWithEngine()
	req, _ := http.NewRequest("POST", "/test?foo=dup", body)
	req.Header.Add("Content-Type", MIMEJSON)
	c.Request = req

	ct := c.ContentType()

	assert.Equal(t, MIMEJSON, ct)

}

func TestBinding(t *testing.T) {
	body := bytes.NewBufferString("{\"foo\":\"bar\", \"name\":\"test\",\"id\":\"\"}")
	//上記みたいな形式のBodyじゃないとJsonでParse出来ないので注意,
	//foo=bar&name=testみたいな形式はダメ
	c := createContextWithEngine()
	req, _ := http.NewRequest("POST", "/test?foo=dup", body)
	req.Header.Add("Content-Type", MIMEJSON)
	c.Request = req

	var u struct {
		Foo  string `json:"foo"`
		Name string `json:"name"`
		Id   string `json:"id"`
	}
	err := c.Bind(&u)
	assert.Nil(t, err)
	assert.Equal(t, "bar", u.Foo)
	assert.Equal(t, "test", u.Name)
	assert.Equal(t, "", u.Id)

}
