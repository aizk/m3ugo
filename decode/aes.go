package decode

import (
	"crypto/aes"
	"crypto/cipher"
)

func AESCBCDecrypt(data, key, iv []byte) (r []byte, err error) {
	secret, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	cipher.NewCBCDecrypter(secret, iv).CryptBlocks(data, data)
	r = data
	// un padding?
	return
}