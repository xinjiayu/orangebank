package orangebank

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/vgmdj/utils/encrypt"
	"github.com/vgmdj/utils/logger"
)

//Sign 签名
type Sign struct {
	OpenKey string
	Values  map[string]interface{}
}

//NewSign 初始化
func NewSign(openKey string, values map[string]interface{}) *Sign {
	return &Sign{
		OpenKey: openKey,
		Values:  values,
	}
}

//ToSign 签名
func (s Sign) ToSign() string {
	s.Values["open_key"] = s.OpenKey

	var keys []string
	for k := range s.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var str string
	for k, v := range keys {
		if k == 0 {
			str = fmt.Sprintf("%s=%v", v, s.Values[v])
			continue
		}
		str = fmt.Sprintf("%s&%s=%v", str, v, s.Values[v])
	}

	s1 := sha1.New()
	io.WriteString(s1, str)
	s1Sign := s1.Sum(nil)

	m5 := md5.New()
	io.WriteString(m5, hex.EncodeToString(s1Sign))

	delete(s.Values, "open_key")

	logger.Info(str, hex.EncodeToString(s1Sign), hex.EncodeToString(m5.Sum(nil)))

	return fmt.Sprintf("%x", m5.Sum(nil))

}

type AES struct {
	OpenKey string
	Values  map[string]interface{}
}

func NewAES(openKey string, values interface{}) *AES {
	v, ok := values.(map[string]interface{})
	if !ok {
		v = make(map[string]interface{})
		if text, ok := values.(string); ok {
			v["cipherText"] = text
		}
	}

	return &AES{
		OpenKey: openKey,
		Values:  v,
	}
}

func (aes *AES) Encrypt() (data string, err error) {
	logger.Info(aes.Values)

	text, err := json.Marshal(aes.Values)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	var tmp []byte
	tmp, err = encrypt.AesECBEncrypt(string(text), aes.OpenKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	data = hex.EncodeToString(tmp)

	logger.Info(string(text), data)

	return
}

func (aes *AES) Decrypt() map[string]interface{} {
	text, ok := aes.Values["cipherText"].(string)
	if !ok || text == "" {
		return nil
	}

	bts, err := hex.DecodeString(text)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	data, err := encrypt.AesECBDecrypt(string(bts), aes.OpenKey)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	logger.Info(string(data))

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return m
}
