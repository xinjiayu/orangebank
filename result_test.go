package orangebank

import (
	"testing"
)

func TestResult(t *testing.T) {
	result := []byte("amount=50&open_id=3fcba3858214ed9e74056bf395a3519e&ord_no=9381536134830923197662797&out_no=2018090596025907&pay_time=20180905160948&rand_str=hDm60NqlETwvgn0RVZJtVuimCHjRfRmwvIDvw3RaXnGdKH5GGO9CjsYVahMq99WFSAboE2zBqfPaXURDJ1f2ueYEwL4FU1lMBwbfzLRZ0GaYB1Bk3Rnx5lcB6gh1iCOT&status=1&timestamp=1536134990&trade_result=%7B%22transaction_id%22%3A%224200000179201809053662483815%22%2C%22nonce_str%22%3A%22aNdL5XOArQYsjtqoT0DLSSKY9I8EMJL9%22%2C%22trade_state%22%3A%22SUCCESS%22%2C%22bank_type%22%3A%22CFT%22%2C%22openid%22%3A%22oRfOu1YlPQtwsQwK9ZIoxeLCCETw%22%2C%22sign%22%3A%22dW5CkAjLOvxF%2B80xkpPw0OHOEFtQSdoan69b3%2BfI8MkKneWB%5C%2Fr5j3OVgDX%2B2UMVs9EAoVB3Cvjo2%2BYr14VlstTmtHvUB84SLkSn1%2B5Evr%2BEZbYdOaFpGhoiYPbnxZaMwkWdixxWkHD7VDbk0BtnVlK3hCRKH2JnDoTCmHHzlMXBogj4nHm9DmimgLC8Lk8FZs5XDXmr3ZdsZZnhPFkUpnsZuOldRt5wzmwC33x0HZ92fbg1oejZQF0QQ5PlD%2BWrOdc%5C%2FP9j4JhVIBBS48dELbZsUDioOIri5Z2CPSkzZXnjO1hvqFl4mb3t9N%2B1ZkCy6twvFeE73Tj88RaxWI2%2B2yLg%3D%3D%22%2C%22fee_type%22%3A%22CNY%22%2C%22mch_id%22%3A%221456738702%22%2C%22sub_mch_id%22%3A%22238469967%22%2C%22cash_fee%22%3A%2250%22%2C%22out_trade_no%22%3A%229381536134830923197662797%22%2C%22cash_fee_type%22%3A%22CNY%22%2C%22coupon_fee%22%3A%220%22%2C%22total_fee%22%3A%2250%22%2C%22appid%22%3A%22wxd758a38cf06aa4f4%22%2C%22trade_state_desc%22%3A%22%5Cu652f%5Cu4ed8%5Cu6210%5Cu529f%22%2C%22settlement_total_fee%22%3A%2250%22%2C%22trade_type%22%3A%22JSAPI%22%2C%22result_code%22%3A%22SUCCESS%22%2C%22attach%22%3A%22bank_mch_name%3D%26bank_mch_id%3D887980225%22%2C%22time_end%22%3A%2220180905160948%22%2C%22return_code%22%3A%22SUCCESS%22%7D&sign=5705605448d3664b522c1f3d973c082a")

	client := NewClient("73b24f53ffc64486eb40d606456fb04d", "252dec9c5c3279d83feb3653c3db545b", EnvDEV,
		nil, nil)

	pcb, err := client.ParseCallBack(result)
	t.Error(pcb, err)

}
