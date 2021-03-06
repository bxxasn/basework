package rsa

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	cfg "github.com/bsir2020/basework/configs"
	"io/ioutil"
)

var (
	publicKey  []byte
	privateKey []byte
	pub        *rsa.PublicKey
)

func init() {
	var err error
	privateKey, err = ioutil.ReadFile(cfg.EnvConfig.Authkey.PrivateKey)
	if err != nil {
		panic("read private pem fail")
	}

	publicKey, err = ioutil.ReadFile(cfg.EnvConfig.Authkey.Publickey)
	if err != nil {
		panic("read public pem fail")
	}
}

// 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub = pubInterface.(*rsa.PublicKey)

	partLen := pub.N.BitLen()/8 - 11
	chunks := split([]byte(origData), partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, pub, chunk)
		if err != nil {
			return nil, err
		}

		buffer.Write(bytes)
	}

	return buffer.Bytes(), nil

	//return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext string) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	partLen := pub.N.BitLen() / 8
	chunks := split([]byte([]byte(ciphertext)), partLen)
	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, priv, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}

	return buffer.Bytes(), nil
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
