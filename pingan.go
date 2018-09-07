package orangebank

import (
	"github.com/vgmdj/utils/chars"
)

//Client 平安银行client
type Client struct {
	openID  string
	openKey string
	env     Env

	publicKey  []byte
	privateKey []byte
}

//NewClient 初始化
func NewClient(openID, openKey string, env Env, publicKey, privateKey []byte) *Client {
	return &Client{
		openID:     openID,
		openKey:    openKey,
		env:        env,
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

//BaseURL 基础url，默认为生产环境
func (c *Client) BaseURL() string {
	switch c.env {
	default:
		return "https://api.orangebank.com.cn/mct1/"

	case EnvDEV:
		return "https://mixpayuat4.orangebank.com.cn/mct1/"

	case EnvPRO:
		return "https://api.orangebank.com.cn/mct1/"
	}
}

//format 格式化数据
func (c *Client) format(values map[string]interface{}) map[string]string {
	kv := make(map[string]string)
	for k, v := range values {
		kv[k] = chars.ToString(v)
	}

	return kv
}

//Env 环境选择
type Env int

const (
	//EnvDEV 测试环境
	EnvDEV Env = iota
	//EnvPRO 生产环境
	EnvPRO
)

//PmtTag 支付标签选择
type PmtTag int

const (
	//TagWeiXin 微信支付
	TagWeiXin PmtTag = iota
	//TagAliPayCS 支付宝支付
	TagAliPayCS
	//TagJdPay 京东支付
	TagJdPay
	//TagJdOL 京东在线
	TagJdOL
	//TagQCS QQ支付
	TagQCS
	//TagCash 现金支付，银行支付
	TagCash
	//TagDiy 自选支付方式
	TagDiy
)

var (
	devTagList = []string{"WeixinBERL", "AlipayCS", "Jdpay", "JdOL", "QpayCS", "Cash", "Diy"}
	proTagList = []string{"Weixin", "AlipayPAZH", "Jdpay", "JdOL", "Qpay", "Cash", "Diy"}
)

//ToString 标签转换
func (pt PmtTag) ToString(env Env) string {
	if env == EnvDEV {
		return devTagList[int(pt)]
	}

	return proTagList[int(pt)]
}
