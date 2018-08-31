package orangebank

import (
	"testing"
)

func TestConifg(t *testing.T) {
	client := NewClient("73b24f53ffc64486eb40d606456fb04d", "7386072b1f94fdd7acaae83cd0f0f1c1", EnvDEV)

	conf := ConfigRequest{
		PmtTag:   TagWeiXin,
		SubAppID: "wxbe4c8b8be110dc3b",
	}

	err := client.ConfigAdd(conf)
	t.Error(err)
}
