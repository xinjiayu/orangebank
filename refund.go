package orangebank

import (
	"fmt"
	"time"

	"github.com/vgmdj/utils/chars"
	"github.com/vgmdj/utils/httplib"
	"github.com/vgmdj/utils/logger"
)

type RefundRequest struct {
	SignType      string //加密类型RSA或RSA2
	OutNo         string //原始订单号
	RefundOutNo   string //退款订单号
	RefundOutName string //退款订单名称
	RefundAmount  int    //退款金额
	TradeAccount  string //交易账号
	TradeNo       string //交易号
	TradeResult   string //收单机构原始交易信息
	TmlToken      string //终端令牌，终端上线后获得的令牌
	Remark        string //退款备注
	ShopPass      string //主管密码
}

type RefundResp struct {
	ErrCode   int
	TimeStamp int64
	Msg       string
	Data      RefundData
}

type RefundData struct {
	OrdNo         string
	OrdShopID     int
	OrdMctID      int
	TradeAmount   int
	TradeNo       string
	TradeResult   string
	OriginalOrdNo int
	Status        int
	OrdCurrency   string
	OutNo         string
	TradeTime     string
}

func (rr *RefundResp) parse(privateKey []byte, values map[string]interface{}) (err error) {
	logger.Info(values)

	rr.ErrCode = chars.ToInt(values["errcode"])
	rr.Msg = chars.ToString(values["msg"])
	rr.TimeStamp = int64(chars.ToInt(values["timestamp"]))

	data := chars.ToString(values["data"])
	logger.Info(data)

	rsa := NewRSA(nil, privateKey)
	m, err := rsa.Decrypt(data)
	if err != nil {
		logger.Error(m)
		return
	}

	logger.Info(m)

	rr.Data.OrdMctID = chars.ToInt(m["ord_mch_id"])
	rr.Data.OrdShopID = chars.ToInt(m["ord_shop_id"])
	rr.Data.OrdNo = chars.ToString(m["ord_no"])
	rr.Data.OrdCurrency = chars.ToString(m["ord_currency"])
	rr.Data.TradeAmount = chars.ToInt(m["trade_amount"])
	rr.Data.TradeResult = chars.ToString(m["trade_result"])
	rr.Data.Status = chars.ToInt(m["status"])

	return
}

func (c *Client) PayRefund(req RefundRequest) (refundResp RefundResp, err error) {
	m := make(map[string]interface{})
	m["out_no"] = req.OutNo
	m["refund_out_no"] = req.RefundOutNo
	m["refund_ord_name"] = req.RefundOutName
	m["refund_amount"] = req.RefundAmount
	m["trade_account"] = req.TradeAccount
	m["trade_result"] = req.TradeResult
	m["tml_token"] = req.TmlToken
	m["remark"] = req.Remark
	m["shop_pass"] = req.ShopPass

	rsa := NewRSA(c.publicKey, c.privateKey)
	data, err := rsa.Encrypt(m)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	d := make(map[string]interface{})
	d["data"] = data
	d["open_id"] = c.openID
	d["timestamp"] = chars.ToString(time.Now().Unix())
	d["sign_type"] = "RSA"

	resp := make(map[string]interface{})
	err = httplib.PostForm(fmt.Sprintf("%s%s", c.BaseURL(), "payorder"), c.format(d), &resp,
		map[string]string{httplib.ResponseResultContentType: "application/json"})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = refundResp.parse(c.privateKey, resp)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	return
}
