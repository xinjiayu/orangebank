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
	sign := NewSign("OPEN_KEY", values)

	ast.Equal("d51a7e83378b1a5324fe4b06692188d7", sign.ToSign())

}
