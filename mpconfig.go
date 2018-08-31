package orangebank

import (
	"fmt"
	"time"

	"github.com/vgmdj/utils/chars"
	"github.com/vgmdj/utils/httplib"
	"github.com/vgmdj/utils/logger"
)

type ConfigRequest struct {
	PmtTag   PmtTag
	SubAppID string
}

func (c *Client) ConfigAdd(req ConfigRequest) (err error) {
	m := make(map[string]interface{})
	m["pmt_tag"] = req.PmtTag.ToString(c.env)
	m["sub_appid"] = req.SubAppID

	aes := NewAES(c.openKey, m)
	data, err := aes.Encrypt()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	d := make(map[string]interface{})
	d["data"] = data
	d["timestamp"] = chars.ToString(time.Now().Unix())
	d["open_id"] = c.openID

	sign := NewSign(c.openKey, d)
	d["sign"] = sign.ToSign()

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
