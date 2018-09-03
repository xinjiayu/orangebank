package orangebank

import (
	"github.com/vgmdj/utils/chars"
)

type Client struct {
	openID  string
	openKey string
	env     Env

	publicKey  []byte
	privateKey []byte
}

func NewClient(openID, openKey string, env Env, publicKey, privateKey []byte) *Client {
	return &Client{
		openID:     openID,
		openKey:    openKey,
		env:        env,
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

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

func (c *Client) format(values map[string]interface{}) map[string]string {
	kv := make(map[string]string)
	for k, v := range values {
		kv[k] = chars.ToString(v)
	}

	return kv
}

type Env int

const (
	EnvDEV Env = iota
	EnvPRO
)

type PmtTag int

const (
	TagWeiXin PmtTag = iota
	TagAliPayCS
	TagJdPay
	TagJdOL
	TagQCS
	TagCash
	TagDiy
)

var (
	devTagList = []string{"WeixinBERL", "AlipayCS", "Jdpay", "JdOL", "QpayCS", "Cash", "Diy"}
	proTagList = []string{"Weixin", "AlipayPAZH", "Jdpay", "JdOL", "Qpay", "Cash", "Diy"}
)

func (pt PmtTag) ToString(env Env) string {
	if env == EnvDEV {
		return devTagList[int(pt)]
	}

	return proTagList[int(pt)]
}
