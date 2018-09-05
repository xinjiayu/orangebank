package orangebank

import (
	"fmt"
	"time"

	"github.com/vgmdj/utils/chars"
	"github.com/vgmdj/utils/httplib"
	"github.com/vgmdj/utils/logger"
)

type ConfigRequest struct {
	PmtTag         PmtTag
	SubAppID       string
	SubscribeAppID string
}

func (c *Client) ConfigAdd(req ConfigRequest) (err error) {
	m := make(map[string]interface{})
	m["pmt_tag"] = req.PmtTag.ToString(c.env)
	m["sub_appid"] = req.SubAppID
	m["subscribe_appid"] = req.SubscribeAppID

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
