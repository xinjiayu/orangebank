package orangebank

import (
	"fmt"
	"time"

	"github.com/vgmdj/utils/chars"
	"github.com/vgmdj/utils/httplib"
	"github.com/vgmdj/utils/logger"
)

//ConfigRequest 配置请求
type ConfigRequest struct {
	PmtTag         PmtTag //标签
	SubAppID       string //相关联appid，三选一配置
	SubscribeAppID string //推荐关注公众号
	JsAPIPath      string //安全支付域名
}

//ConfigAdd 配置
func (c *Client) ConfigAdd(req ConfigRequest) (err error) {
	m := make(map[string]interface{})
	m["pmt_tag"] = req.PmtTag.ToString(c.env)
	m["sub_appid"] = req.SubAppID
	m["subscribe_appid"] = req.SubscribeAppID
	m["jsapi_path"] = req.JsAPIPath

	aes := NewAES(c.openKey)
	data, err := aes.Encrypt(m)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	d := make(map[string]interface{})
	d["data"] = data
	d["timestamp"] = chars.ToString(time.Now().Unix())
	d["open_id"] = c.openID

	sign := NewSign(c.openKey)
	d["sign"] = sign.ToSign(d)

	resp := make(map[string]interface{})
	err = httplib.PostForm(fmt.Sprintf("%s%s", c.BaseURL(), "mpconfig/add"), c.format(d),
		&resp, map[string]string{httplib.ResponseResultContentType: "application/json"})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info(resp)

	return
}

//ConfigQuery 配置查询
func (c *Client) ConfigQuery(tag PmtTag) (err error) {
	m := make(map[string]interface{})
	m["pmt_tag"] = tag.ToString(c.env)

	aes := NewAES(c.openKey)
	data, err := aes.Encrypt(m)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	d := make(map[string]interface{})
	d["data"] = data
	d["timestamp"] = chars.ToString(time.Now().Unix())
	d["open_id"] = c.openID

	sign := NewSign(c.openKey)
	d["sign"] = sign.ToSign(d)

	resp := make(map[string]interface{})
	err = httplib.PostForm(fmt.Sprintf("%s%s", c.BaseURL(), "mpconfig/query"), c.format(d),
		&resp, map[string]string{httplib.ResponseResultContentType: "application/json"})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info(resp)

	rm, _ := aes.Decrypt(chars.ToString(resp["data"]))
	logger.Info(rm)

	return
}
