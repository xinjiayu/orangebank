package orangebank

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/vgmdj/utils/chars"
	"github.com/vgmdj/utils/httplib"
	"github.com/vgmdj/utils/logger"
)

//RefundRequest 退款请求
type RefundRequest struct {
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

//RefundResp 退款回复
type RefundResp struct {
	ErrCode   int
	Timestamp int
	Msg       string
	Data      RefundData
	Sign      string
}

//RefundData 退款回复关键数据
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

//parse 退款回复关键数据解析
func (rr *RefundResp) parse(openKey string, privateKey []byte, values map[string]interface{}) (err error) {
	values["timestamp"] = chars.ToInt(values["timestamp"])

	rr.ErrCode = chars.ToInt(values["errcode"])
	rr.Msg = chars.ToString(values["msg"])
	rr.Timestamp = chars.ToInt(values["timestamp"])
	rr.Sign = chars.ToString(values["sign"])

	data := chars.ToString(values["data"])

	if data == "" {
		logger.Info(values)
		return
	}

	sign := NewSign(openKey)
	if rr.Sign != sign.ToSign(values) {
		logger.Error(values, rr.Sign, sign.ToSign(values))
		return fmt.Errorf("退款签名认证不通过")
	}

	aes := NewAES(openKey)
	m, err := aes.Decrypt(data)
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

//PayRefund 退款
func (c *Client) PayRefund(req RefundRequest) (refundResp RefundResp, err error) {
	m := make(map[string]interface{})
	m["out_no"] = req.OutNo
	m["refund_out_no"] = req.RefundOutNo
	m["refund_ord_name"] = req.RefundOutName
	m["refund_amount"] = req.RefundAmount
	m["remark"] = req.Remark

	s1 := sha1.New()
	s1.Write([]byte(req.ShopPass))
	m["shop_pass"] = hex.EncodeToString(s1.Sum(nil))

	aes := NewAES(c.openKey)
	data, err := aes.Encrypt(m)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	d := make(map[string]interface{})
	d["data"] = data
	d["open_id"] = c.openID
	d["timestamp"] = chars.ToString(time.Now().Unix())
	d["sign_type"] = "RSA"

	rsa := NewRSA(c.openKey, c.publicKey, c.privateKey)
	sign, err := rsa.Sign(d)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	d["sign"] = sign

	resp := make(map[string]interface{})
	err = httplib.PostForm(fmt.Sprintf("%s%s", c.BaseURL(), "payrefund"), c.format(d), &resp,
		map[string]string{httplib.ResponseResultContentType: "application/json"})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = refundResp.parse(c.openKey, c.privateKey, resp)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	return
}

//RefundQueryResp 退款订单查询
type RefundQueryResp struct {
	ErrCode   int
	Msg       string
	Sign      string
	Timestamp int64
	Data      RefundQueryRespData
}

//RefundQueryRespData 退款查询回复
type RefundQueryRespData struct {
	RefundOrdNo string
	RefundOutNo int
	TradeAmount int
	Status      string
	TradeResult string
}

//parse 退款查询回复解析
func (rqr *RefundQueryResp) parse(key string, values map[string]interface{}) (err error) {
	values["timestamp"] = chars.ToInt(values["timestamp"])
	logger.Info(values)

	rqr.ErrCode = chars.ToInt(values["errcode"])
	rqr.Msg = chars.ToString(values["msg"])
	rqr.Timestamp = int64(chars.ToInt(values["timestamp"]))

	vs, ok := values["sign"]
	if !ok {
		ec, ok := values["errcode"]
		if !ok || chars.ToInt(ec) == 0 {
			return fmt.Errorf("数据返回异常")
		}

		return

	}

	rqr.Sign = chars.ToString(vs)
	sign := NewSign(key)
	delete(values, "sign")
	if rqr.Sign != sign.ToSign(values) {
		logger.Error(values, rqr.Sign)
		logger.Info(fmt.Errorf("sign 验证不通过"))
	}

	data := chars.ToString(values["data"])

	aes := NewAES(key)
	m, err := aes.Decrypt(data)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	rqr.Data.RefundOutNo = chars.ToInt(m["refund_out_no"])
	rqr.Data.RefundOrdNo = chars.ToString(m["refund_ord_no"])
	rqr.Data.TradeAmount = chars.ToInt(m["trade_amount"])
	rqr.Data.Status = chars.ToString(m["status"])
	rqr.Data.TradeResult = chars.ToString(m["trade_result"])

	return
}

//PayRefundQuery 退款查询
func (c *Client) PayRefundQuery(outNo string) (rqResp RefundQueryResp, err error) {
	m := make(map[string]interface{})
	m["refund_out_no"] = outNo

	aes := NewAES(c.openKey)
	data, err := aes.Encrypt(m)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	d := make(map[string]interface{})
	d["data"] = data
	d["open_id"] = c.openID
	d["timestamp"] = chars.ToString(time.Now().Unix())

	sign := NewSign(c.openKey)
	d["sign"] = sign.ToSign(d)

	resp := make(map[string]interface{})
	err = httplib.PostForm(fmt.Sprintf("%s%s", c.BaseURL(), "payrefundquery"), c.format(d), &resp,
		map[string]string{httplib.ResponseResultContentType: "application/json"})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = rqResp.parse(c.openKey, resp)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	return
}
