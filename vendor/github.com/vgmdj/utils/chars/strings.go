package chars

import (
	"errors"

	"github.com/vgmdj/utils/logger"
)

//Replace 将给定的str，从offset位开始，一直替换limit位的repStr
func Replace(str string, offset int, limit int, repStr string) (result string, err error) {
	if len(str) < 1 || len(repStr) != 1 || offset+limit-1 > len(str) {
		err = errors.New("字符长度不合法或函数使用错误")
		logger.Error(str)
		return
	}

	repByte := repStr[0]
	bts := []byte(str)
	for i := offset - 1; i < offset+limit-1 && i < len(str); i++ {
		bts[i] = repByte
	}
	result = string(bts)

	return
}
