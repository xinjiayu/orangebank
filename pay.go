package orangebank

import (
	"fmt"
	"time"

	"github.com/vgmdj/utils/chars"
	"github.com/vgmdj/utils/httplib"
	"github.com/vgmdj/utils/logger"
)

//PayRequest 金额统一以分为单位
type PayRequest struct {
	OutNO          string //商户订单流水号
	PmtTag         PmtTag //支付模式
	PmtName        string //diy模式下，支付名称
	OrdName        string //订单名称
	OriginalAmount int    //原始金额-订单金额
	DiscountAmount int    //优惠金额
	IgnoreAmount   int    //抹零金额
	NotifyURL      string //异步回调通知地址
	TradeAccount   string //交易账号，可为空
	TradeNO        string //交易号，可为空
	Remark         string //订单备注，可为空
	Tag            string //订单附加数据，可为空

	JsAPI       string //使用jsapi调用方式值为1
	SubAppID    string //商户appid
	SubOpenID   string //商户下用户id
	JumpURL     string //支付后跳转页面
	LimitCredit bool   //是否限制使用信用卡
}

//PayResponse 下单回复
type PayResponse struct {
	ErrCode   int
	Msg       string
	TimeStamp int64
	Data      PayData
	Sign      string
}

//PayData 下单数据
type PayData struct {
	PmtName        string
	PmtTag         string
	OrdMctID       int
	OrdShopID      int
	OrdNO          string
	OrdType        int
	OrdCurrency    string
	OrdMctRatio    string
	OriginalAmount int
	DiscountAmount int
	IgnoreAmount   int
	TradeAccount   string
	TradeAmount    int
	TradeNo        string
	TradeQrcode    string
	TradePayTime   string
	Status         int
	TradeResult    string
	OutNO          string
	JsAPIPayURL    string
	AppID          string
	NonceStr       string
	SignType       string
	Package        string
	PaySign        string
}

func (pr *PayResponse) parse(key string, values map[string]interface{}) (err error) {
	values["timestamp"] = chars.ToInt(values["timestamp"])
	logger.Info(values)

	vs, ok := values["sign"]
	if !ok {
		ec, ok := values["errcode"]
		if !ok || chars.ToInt(ec) == 0 {
			return fmt.Errorf("数据返回异常")
		}
	}

	pr.Sign = chars.ToString(vs)
	sign := NewSign(key)
	delete(values, "sign")
	if pr.Sign != sign.ToSign(values) {
		logger.Error(values, pr.Sign)
		logger.Info(fmt.Errorf("sign 验证不通过"))
	}

	pr.ErrCode = chars.ToInt(values["errcode"])
	pr.Msg = chars.ToString(values["msg"])
	pr.TimeStamp = int64(chars.ToInt(values["timestamp"]))

	data := chars.ToString(values["data"])
	logger.Info(data)

	aes := NewAES(key)
	m, err := aes.Decrypt(data)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info(m)

	pr.Data.PmtTag = chars.ToString(m["pmt_tag"])
	pr.Data.PmtName = chars.ToString(m["pmt_name"])
	pr.Data.OriginalAmount = chars.ToInt(m["original_amount"])
	pr.Data.DiscountAmount = chars.ToInt(m["discount_amount"])
	pr.Data.IgnoreAmount = chars.ToInt(m["ignore_amount"])
	pr.Data.OrdMctID = chars.ToInt(m["ord_mch_id"])
	pr.Data.OrdShopID = chars.ToInt(m["ord_shop_id"])
	pr.Data.OrdNO = chars.ToString(m["ord_no"])
	pr.Data.OrdType = chars.ToInt(m["ord_type"])
	pr.Data.OrdMctRatio = chars.ToString(m["ord_mct_ratio"])
	pr.Data.OrdCurrency = chars.ToString(m["ord_currency"])
	pr.Data.TradeAccount = chars.ToString(m["trade_account"])
	pr.Data.TradeAmount = chars.ToInt(m["trade_amount"])
	pr.Data.TradePayTime = chars.ToString(m["trade_pay_time"])
	pr.Data.TradeQrcode = chars.ToString(m["trade_qrcode"])
	pr.Data.TradeResult = chars.ToString(m["trade_result"])
	pr.Data.JsAPIPayURL = chars.ToString(m["jsapi_pay_url"])
	pr.Data.Status = chars.ToInt(m["status"])
	pr.Data.AppID = chars.ToString(m["appId"])
	pr.Data.Package = chars.ToString(m["package"])
	pr.Data.PaySign = chars.ToString(m["paySign"])
	pr.Data.SignType = chars.ToString(m["signType"])
	pr.Data.NonceStr = chars.ToString(m["nonceStr"])

	return
}

//PayOrder 统一下单
func (c *Client) PayOrder(req PayRequest) (payResp PayResponse, err error) {
	m := make(map[string]interface{})
	m["out_no"] = req.OutNO
	m["pmt_tag"] = req.PmtTag.ToString(c.env)
	m["pmt_name"] = req.PmtName
	m["ord_name"] = req.OrdName
	m["original_amount"] = req.OriginalAmount
	m["discount_amount"] = req.DiscountAmount
	m["ignore_amount"] = req.IgnoreAmount
	m["trade_amount"] = req.OriginalAmount - req.DiscountAmount - req.IgnoreAmount
	m["trade_account"] = req.TradeAccount
	m["trade_no"] = req.TradeNO
	m["remark"] = req.Remark
	m["tag"] = req.Tag
	m["notify_url"] = req.NotifyURL
	m["jump_url"] = req.JumpURL

	if !req.LimitCredit {
		m["limit_pay"] = "no_credit"
	}

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
	err = httplib.PostForm(fmt.Sprintf("%s%s", c.BaseURL(), "payorder"), c.format(d), &resp,
		map[string]string{httplib.ResponseResultContentType: "application/json"})
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = payResp.parse(c.openKey, resp)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	return
}
