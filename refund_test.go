package orangebank

import (
	"testing"
)

func TestPayRefund(t *testing.T) {
	privateKey := []byte(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDoYLKYSup2POAg
mf2iHwGhUQZSeyyxLvslDcFCRekg1D7w+Y8dybhBOyDvlTbuG15jt8vqqH1xlGTI
66yz351uxrUL10UL6Uw+kx1091/BCiRo1AY9frOlIeTnVZy/S2LPWjEefI9o0nOB
YJkoJZsaFffX6qXVFkCfAGNMohyuU0gNHw3SXXMBiJ2n2E7GfFf/peAmtTg4zFJO
+CrikfSTrL0u5Y2dWQp1KQvkfCMrXaECM24fLc6xxwAdiARC6DjMK8F+ZxovEgb0
VObpbZ5pnV2XrKEbiDXQ/DbeNTmELHN8k0S3NUY6ozpu1M3JDqZPhBIDgd/gcZkn
FngCRJ95AgMBAAECggEAXZYkF0WEq93UfgzGozZNl8RkAW/uDeXX65JglOpG+5u/
RZmcU+jbthm0KAk2OCr5lrt8+qKk8stK08hmo4KZivWoEH7AJg3tUP46zNKb08jb
5QQPB1Ex1H2UDL7kA/6+arfuNFMCBrtLHX3j8NFEZ/sU9/ZelzUBDYhAdaqMVoAb
fOX7MIZMRIZRhrZCuQlZI0Wz1oGzOf4WB40HEvgRfkGkdC2C+co9tCK7PRcJCiAn
Sro175E8NLmYIs+0yUtztUUPSUcVmZ4NfQzM/tD7oyeHREkmBEM6NT62U2PAHyhd
FJH23T5Rq80u2Qbzpqv6Fk9aNlBM5uN24umVU7C2oQKBgQD0Hyw7SKOGExbBOOyJ
FLpeISsiACokPz7rGwgQJfNfGrbnBLuo7VyxRAX9y01FQrsZL0ZqZ09dSnY6Kqjm
MzUQdtF1mteiF9wC6xHcIN9oXzDtDc0geB/hb3zLRIHs5yrCgNci8FjGL0OwBaNN
z0vLa02C8Vm5P6rjlclyjK7ZtQKBgQDzrz0i+mCgklD6+cbEA3Z+0b75s5Lg9cu1
3eBntZ/+11d2Kf1JCrgrCiRcX1ggOG2ZHmPmi5cs93J6rd+HWInPWJ+j4GUYN5wP
S8GxiQun42b3FloT5Zoe4Cj3TmonIHvhFQojYAniKPx3jb3mCLOUyHWjIc5r9zE6
Tovth195NQKBgE7k9CqEozRlXuk7OFZk+IYLOiFW5EeqmO7qYYS2fxyxSYMHqI5D
h71SOo128pX7pvPQr3Ubxi5kLilGOCeNTQzxGWhkjmO4SkY3KiJ2DT1x5iH2X+Cq
ccMtgKtAjKy/WLZbZSvJeSczhzCP4eL3p4sqNnanAVQ5G0VJ1zzJ8ogxAoGAA+4i
nUrOfih9995Jb2Xi5l65pstXphswwukmMmYCg5izh2tb826h08fhGEBNao+ebObJ
k7FSqd3/0ay2OzeZWWfDg2AeIUrcUH7XS+a68mU/huKsZz+/wZm572srWSAz/0hY
loN5BVXF5KO7mVcwlki5ZP0pmCIvgBI+PYF+b7UCgYEA3WpyaYgDSVrE2pcF4+dv
qFNIDGcFBGqCyeCaGx4E1WuP0R7y1SFM2miGnbl+te89ZLSzSudoeFq/qS8S3W+M
+PuQJlmlFnQ+B/DaEJmKGiiNAqWT9cdohLeYNaJBvk55MdnRJl4dAfkf5PFIf64C
WvdDEvhehuV6kub3i8KEvH8=
-----END PRIVATE KEY-----`)

	client := NewClient("73b24f53ffc64486eb40d606456fb04d", "7386072b1f94fdd7acaae83cd0f0f1c1", EnvDEV,
		nil, privateKey)

	req := RefundRequest{
		OutNo:        "2018090596025907",
		RefundOutNo:  "0000000000000001",
		RefundAmount: 50,
		ShopPass:     "123456",
	}

	resp, err := client.PayRefund(req)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(resp)

}

func TestRefundQuery(t *testing.T) {

	client := NewClient("3fcba3858214ed9e74056bf395a3519e",
		"252dec9c5c3279d83feb3653c3db545b", EnvPRO, nil, nil)

	resp, err := client.PayRefundQuery("2018090596025907")
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log(resp)

}
