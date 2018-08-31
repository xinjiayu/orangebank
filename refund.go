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

func (rr *RefundResp)parse(key string, values map[string]interface{})(err error){

	return
}



func (c *Client) PayRefund(req RefundRequest) (refundResp RefundResp, err error) {
	m := make(map[string]interface{})
	m["out_no"] = req.OutNo
	m["trade_account"] = req.TradeAccount
	m["trade_no"] = req.TradeNo
	m["remark"] = req.Remark

	aes := NewAES(c.openKey, m)
	data, _ := aes.Encrypt()

	d := make(map[string]interface{})
	d["data"] = data
	d["open_id"] = c.openID
	d["timestamp"] = chars.ToString(time.Now().Unix())

	sign := NewSign(c.openKey, d)
	d["sign"] = sign.ToSign()

	resp := make(map[string]interface{})
	err = httplib.PostForm(fmt.Sprintf("%s%s", c.BaseURL(), "payorder"), c.format(d), &resp,
		map[string]string{httplib.ResponseResultContentType: "application/json"})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = refundResp.parse(c.openKey, resp)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	return
}
