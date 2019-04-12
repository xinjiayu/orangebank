package chars

import (
	"sync"
)

const (
	//DefaultMultiples 缺省缩放倍数
	DefaultMultiples = 100
)

var (
	single *Conversion
	once   sync.Once
)

//Conversion 放大缩小转换
type Conversion struct {
	base      float64
	result    Result
	multiples int
}

//NewConversion first param is base, second param is multiples
func NewConversion(params ...interface{}) *Conversion {
	c := new(Conversion)

	c.base = 0
	if len(params) != 0 {
		c.base = ToFloat64(params[0])
	}

	c.multiples = parseMultiples(params[:]...)

	return c
}

//Sc single pattern return the same conversion
func Sc() *Conversion {
	once.Do(func() {
		single = NewConversion()
	})
	return single
}

//SetMultiples 设置倍数
func (c *Conversion) SetMultiples(multiples int) {
	c.multiples = multiples
}

//BaseValue 设置基数
func (c Conversion) BaseValue() float64 {
	return c.base
}

//ZoomOut 放大
func (c Conversion) ZoomOut(base ...interface{}) Result {
	if len(base) != 0 {
		c.base = ToFloat64(base[0])
	}

	c.result = Result(c.base * float64(c.multiples))
	return c.result
}

//ZoomIn 缩小
func (c Conversion) ZoomIn(base ...interface{}) Result {
	if len(base) != 0 {
		c.base = ToFloat64(base[0])
	}

	c.result = Result(c.base / float64(c.multiples))
	return c.result
}

func parseMultiples(params ...interface{}) int {
	if len(params) != 2 {
		return DefaultMultiples
	}

	m, ok := params[1].(int)
	if !ok || m%10 != 0 {
		return DefaultMultiples
	}

	return m
}

func precision(multiples int) int {
	var p int
	for multiples != 0 {
		multiples /= 10
		p++
	}

	return p - 1

}

//Result 转换结果
type Result float64

//ToString 字符显示
func (r Result) ToString(multiples ...int) string {
	var m = DefaultMultiples
	if len(multiples) != 0 {
		m = multiples[0]
	}

	return ToString(float64(r), precision(m))
}

//ToInt 整数显示
func (r Result) ToInt() int {
	return ToInt(float64(r))
}

//ToFloat64 浮点数显示
func (r Result) ToFloat64() float64 {
	return ToFloat64(float64(r))
}
