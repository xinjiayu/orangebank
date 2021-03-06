package orangebank

import (
	"fmt"
	"net/url"

	"github.com/vgmdj/utils/logger"
)

//PayCallBack 支付回复
type PayCallBack struct {
	Timestamp   string
	OrdNo       string
	OutNo       string
	Status      string
	Amount      string
	PayTime     string
	RandStr     string
	Sign        string
	TradeResult string
	OpenID      string
}

const (
	//CallBackNotifySuccess 支付回复确认
	CallBackNotifySuccess = "notify_success"
)

//ParseCallBack 支付回复数据解析
func (c *Client) ParseCallBack(body []byte) (pcb PayCallBack, err error) {
	values, err := url.ParseQuery(string(body))

	m := make(map[string]string)
	for k, v := range values {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}

	pcb.OrdNo = m["ord_no"]
	pcb.OutNo = m["out_no"]
	pcb.Status = m["status"]
	pcb.Amount = m["amount"]
	pcb.PayTime = m["pay_time"]
	pcb.RandStr = m["rand_str"]
	pcb.Sign = m["sign"]
	pcb.TradeResult = m["trade_result"]
	pcb.OpenID = m["open_id"]

	delete(m, "sign")

	err = pcb.CheckSign(c.openKey, m)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	return

}

//CheckSign 支付回复验签
func (pcb *PayCallBack) CheckSign(openKey string, m map[string]string) (err error) {
	cm := make(map[string]interface{})
	for k, v := range m {
		cm[k] = v
	}

	sign := NewSign(openKey)
	check := sign.ToSign(cm)
	if pcb.Sign != check {
		logger.Error(fmt.Sprintf("check sign error , %s not match %s", pcb.Sign, check))
		return fmt.Errorf("check sign error , %s not match %s", pcb.Sign, check)
	}

	return
}
