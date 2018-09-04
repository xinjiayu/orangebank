package orangebank

import (
	"errors"
	"fmt"

	"github.com/vgmdj/utils/logger"
)

type PayCallBack struct {
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
	Timestamp   string `json:"timestamp"`
}

const (
	CallBackNotifySuccess = "notify_success"
)

func (pcb PayCallBack) CheckSign(openKey string) (err error) {
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
