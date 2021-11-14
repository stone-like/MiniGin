package binding

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type user struct {
	Name string
	Id   string
}

func requestWithBody(method, path, body string) (req *http.Request) {
	req, _ = http.NewRequest(method, path, bytes.NewBufferString(body))
	return
}

func testBodyBinding(t *testing.T, b Binding, name, path, badPath, body, badBody string) {
	assert.Equal(t, name, b.Name())

	obj := user{}
	req := requestWithBody("POST", path, body)
	err := b.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "testOK", obj.Name)
	assert.Equal(t, "okok", obj.Id)

	obj = user{}
	req = requestWithBody("POST", badPath, badBody)
	err = JSON.Bind(req, &obj)
	assert.Error(t, err)
}

//jsonDecode自体は別にjsonに使うparamが含まれていなくてもエラーははかない
//ex. user{Name:aaa}でbodyが{BadName:bbb}だったとしても特にエラーは出ずDecodeされないだけ
//なのでvalidationをどこかに適用する必要がある
//requestに対してするか、それともBindingの最後にするか、それともEntityの時にするか

func TestJson(t *testing.T) {

	obj := user{}
	req := requestWithBody("POST", "/", `{"Name":"testOK","Id":"okok"}`)
	err := jsonBinding{}.Bind(req, &obj)
	assert.NoError(t, err)
	assert.Equal(t, "testOK", obj.Name)
	assert.Equal(t, "okok", obj.Id)
}
