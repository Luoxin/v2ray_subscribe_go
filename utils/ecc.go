package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go/types"
	"math/big"

	"github.com/Luoxin/Eutamias/utils/json"

	"github.com/ethereum/go-ethereum/crypto/ecies"
	log "github.com/sirupsen/logrus"
)

type Ecc struct {
	prv2    *ecies.PrivateKey
	private *ecdsa.PrivateKey
}

func NewEcc() *Ecc {
	return &Ecc{}
}

func (p *Ecc) Init(key string) error {
	// 初始化椭圆曲线
	pubkeyCurve := elliptic.P256()

	reader := bytes.NewBufferString(key + "0000000000000000000000000000000000000000")
	// reader.WriteString()

	// 随机挑选基点,生成私钥
	private, err := ecdsa.GenerateKey(pubkeyCurve, reader)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	p.private = private
	p.prv2 = ecies.ImportECDSA(private)
	return nil
}

func (p *Ecc) GetPrivate() *ecdsa.PrivateKey {
	return p.private
}

// ess 加密
func (p *Ecc) ECCEncrypt4Str(pt string) (string, error) {
	return p.ECCEncrypt4Bytes([]byte(pt))
}

// ess 加密
func (p *Ecc) ECCEncrypt4Bytes(pt []byte) (string, error) {
	ct, err := ecies.Encrypt(rand.Reader, &p.prv2.PublicKey, pt, nil, nil)
	if err != nil {
		return "", err
	}

	raw := base64.RawStdEncoding.EncodeToString(ct)
	return raw, nil
}

// ess 加密
func (p *Ecc) ECCEncrypt(pt interface{}) (string, error) {
	switch pt := pt.(type) {
	case string:
		return p.ECCEncrypt4Str(pt)
	case []byte:
		return p.ECCEncrypt4Bytes(pt)
	case types.Slice, types.Map, types.Struct:
		b, err := json.Marshal(pt)
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}

		return p.ECCEncrypt4Bytes(b)
	default:
		b, err := json.Marshal(pt)
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}

		return p.ECCEncrypt4Bytes(b)
	}
}

// ecc(ecies)解密
func (p *Ecc) ECCDecrypt(ct string) (string, error) {
	raw, err := base64.RawStdEncoding.DecodeString(ct)
	if err != nil {
		return "", err
	}

	pt, err := p.prv2.Decrypt(raw, nil, nil)
	return string(pt), err
}

// ecc签名
func (p *Ecc) EccSign(pt []byte) (sign []byte, err error) {
	// 根据明文plaintext和私钥，生成两个big.Ing
	r, s, err := ecdsa.Sign(rand.Reader, p.private, pt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	rs, err := r.MarshalText()
	if err != nil {
		return nil, err
	}
	ss, err := s.MarshalText()
	if err != nil {
		return nil, err
	}
	// 将r，s合并（以“+”分割），作为签名返回
	var b bytes.Buffer
	b.Write(rs)
	b.Write([]byte(`+`))
	b.Write(ss)
	return b.Bytes(), nil
}

// ecc验签
func (p *Ecc) EccSignVer(pt, sign []byte) bool {
	var rint, sint big.Int
	// 根据sign，解析出r，s
	rs := bytes.Split(sign, []byte("+"))

	err := rint.UnmarshalText(rs[0])
	if err != nil {
		log.Errorf("err:%v", err)
		return false
	}

	err = sint.UnmarshalText(rs[1])
	if err != nil {
		log.Errorf("err:%v", err)
		return false
	}

	// 根据公钥，明文，r，s验证签名
	v := ecdsa.Verify(&p.private.PublicKey, pt, &rint, &sint)
	return v
}
