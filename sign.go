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
}

//NewSign 初始化
func NewSign(openKey string) *Sign {
	return &Sign{
		OpenKey: openKey,
	}
}

//ToSign 签名
func (s Sign) ToSign(Values map[string]interface{}) string {
	Values["open_key"] = s.OpenKey

	var keys []string
	for k := range Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var str string
	for k, v := range keys {
		if k == 0 {
			str = fmt.Sprintf("%s=%v", v, Values[v])
			continue
		}
		str = fmt.Sprintf("%s&%s=%v", str, v, Values[v])
	}

	s1 := sha1.New()
	io.WriteString(s1, str)
	s1Sign := s1.Sum(nil)

	m5 := md5.New()
	io.WriteString(m5, hex.EncodeToString(s1Sign))

	delete(Values, "open_key")

	logger.Info(str, hex.EncodeToString(s1Sign), hex.EncodeToString(m5.Sum(nil)))

	return fmt.Sprintf("%x", m5.Sum(nil))

}

//AES AES加密
type AES struct {
	OpenKey string
}

//NewAES 初始化
func NewAES(openKey string) *AES {
	return &AES{
		OpenKey: openKey,
	}
}

//Encrypt 加密
func (aes *AES) Encrypt(Values map[string]interface{}) (data string, err error) {
	logger.Info(Values)

	text, err := json.Marshal(Values)
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

//Decrypt 解密
func (aes *AES) Decrypt(cipherText string) (map[string]interface{}, error) {
	bts, err := hex.DecodeString(cipherText)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	data, err := encrypt.AesECBDecrypt(string(bts), aes.OpenKey)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return m, nil
}

//RSA RSA
type RSA struct {
	publicKey  []byte
	privateKey []byte
}

//NewRSA RSA类初始化
func NewRSA(publicKey, privateKey []byte) *RSA {
	return &RSA{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

func (rsa *RSA) Encrypt(Values map[string]interface{}) (data string, err error) {
	var keys []string
	for k := range Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var str string
	for k, v := range keys {
		if k == 0 {
			str = fmt.Sprintf("%s=%v", v, Values[v])
			continue
		}
		str = fmt.Sprintf("%s&%s=%v", str, v, Values[v])
	}

	cipherText, err := encrypt.RsaEncrypt([]byte(str), rsa.publicKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	return hex.EncodeToString(cipherText), nil
}

func (rsa *RSA) Decrypt(cipherText string) (data map[string]interface{}, err error) {
	plainText, err := encrypt.RsaDecrypt([]byte(cipherText), rsa.privateKey)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info(string(plainText))

	err = json.Unmarshal(plainText, &data)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	logger.Info(data)

	return
}
