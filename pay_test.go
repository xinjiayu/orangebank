package orangebank

import (
	"testing"
)

func TestClient_PayOrder(t *testing.T) {
	client := NewClient("73b24f53ffc64486eb40d606456fb04d", "7386072b1f94fdd7acaae83cd0f0f1c1", EnvDEV,
		nil, nil)

	pr := PayRequest{
		OutNO:          "test0000000001",
		PmtTag:         TagWeiXin,
		OrdName:        "ord支付",
		OriginalAmount: 10000,
		DiscountAmount: 100,
		IgnoreAmount:   0,
		NotifyURL:      "www.baidu.com",
		Remark:         "85sh",
		Tag:            "tag",
		JumpURL:        "www.baidu.com",
	}

	resp, err := client.PayOrder(pr)
	t.Error(resp, err)
}
