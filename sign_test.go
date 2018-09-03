package orangebank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	ast := assert.New(t)

	values := map[string]interface{}{
		"data":      "DATA",
		"open_id":   "OPEN_ID",
		"timestamp": "1234567890",
	}
	sign := NewSign("OPEN_KEY")

	ast.Equal("6cc408dc9b34d1bcdc3d6111f4e27a81", sign.ToSign(values))

}
