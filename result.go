package orangebank

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vgmdj/utils/chars"
	"github.com/vgmdj/utils/logger"
)

type PayCallBack struct {
	ErrCode     int    `json:"errcode"`
	Msg         string `json:"msg"`
	Data        string `json:"data"`
	Timestamp   string `json:"timestamp"`
	TradeType   string `json:"trade_type"`
	OrdNo       string `json:"ord_no"`
	OutNo       string `json:"out_no"`
	Status      string `json:"status"`
	Amount      string `json:"amount"`
	PayTime     string `json:"pay_time"`
	RandStr     string `json:"rand_str"`
	Sign        string `json:"sign"`
	TradeResult string `json:"trade_result"`
	OpenID      string `json:"open_id"`
}

const (
	CallBackNotifySuccess = "notify_success"
)

func (pcb *PayCallBack) Parse(body []byte, openKey string) (err error) {
	err = json.Unmarshal(body, pcb)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	aes := NewAES(openKey)
	m, err := aes.Decrypt(pcb.Data)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	pcb.TradeType = chars.ToString(m["trade_type"])
	pcb.OrdNo = chars.ToString(m["ord_no"])
	pcb.OutNo = chars.ToString(m["out_no"])
	pcb.Status = chars.ToString(m["status"])
	pcb.Amount = chars.ToString(m["amount"])
	pcb.PayTime = chars.ToString(m["pay_time"])
	pcb.RandStr = chars.ToString(m["rand_str"])
	pcb.Sign = chars.ToString(m["sign"])
	pcb.TradeResult = chars.ToString(m["trade_result"])
	pcb.OpenID = chars.ToString(m["open_id"])

	return pcb.CheckSign(openKey)

}

func (pcb *PayCallBack) CheckSign(openKey string) (err error) {
	m := make(map[string]interface{})
	m["trade_type"] = pcb.TradeType
	m["ord_no"] = pcb.OrdNo
	m["out_no"] = pcb.OutNo
	m["status"] = pcb.Status
	m["amount"] = pcb.Amount
	m["pay_time"] = pcb.PayTime
	m["rand_str"] = pcb.RandStr
	m["trade_result"] = pcb.TradeResult
	m["open_id"] = pcb.OpenID
	m["timestamp"] = pcb.Timestamp

	sign := NewSign(openKey)
	check := sign.ToSign(m)
	if pcb.Sign != check {
		logger.Error(fmt.Sprintf("check sign error , %s not match %s", pcb.Sign, check))
		return errors.New(fmt.Sprintf("check sign error , %s not match %s", pcb.Sign, check))
	}

	return
}
