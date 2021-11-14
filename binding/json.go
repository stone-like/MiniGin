package binding

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonBinding struct{}

func (jsonBinding) Name() string {
	return "json"
}

func (jsonBinding) Bind(req *http.Request, obj interface{}) error {
	//ここでrequestのクエリと、formを取ってきて、candidateにセット、paramsと合わせてゼロだったらBindは絶対できないのでエラー
	//最初はreq.BodyとParam,Queryすべてを合わせてJsonにパースしようと思ったんだけど、
	//結構面倒くさそうなので、それはユーザーに任せた方がいい気もする
	if req == nil || req.Body == nil {
		return errors.New("invalid request")
	}
	return decodeJSON(req.Body, obj)
}

//req.Bodyはio.ReadCloserだけど受けるときはio.Reader、常に使う最小のインターフェースで受ける
func decodeJSON(r io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}

	//後々ここにvalidation
	return nil

}
