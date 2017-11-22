package network

import (
    "crypto/cipher"
    "crypto/aes"
    "errors"
    "io"
    "crypto/rand"
    "fmt"
    "encoding/hex"
    "log"
)

type AesEncrypt struct {
    key []byte
    blockSize int
}

func (c *AesEncrypt) SetKey(key []byte) error {

    size := len(key)
    if(size != 16 && size != 24 && size != 32) {
        log.Println(size)
        return errors.New("error key length")
    }
    c.key = key
    c.blockSize = size
    return nil
}

//加密字符串
func (c *AesEncrypt) Encrypt(src []byte) ([]byte, error) {
    encrypted := make([]byte, c.blockSize + len(src))
    iv := encrypted[:c.blockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, err
    }
    block, err := aes.NewCipher(c.key)
    if err != nil {
        log.Println(err, c.key)
        return nil, err
    }
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(encrypted[c.blockSize:], src)
    return encrypted, nil
}

func (c *AesEncrypt) Decrypt(src []byte) ([]byte, error) {
    // hex
    decryptText, err := hex.DecodeString(fmt.Sprintf("%x", string(src)))
    if err != nil {
        log.Println(err)
        return nil, err
    }

    if len(decryptText) < c.blockSize {
        log.Println("crypto/cipher: ciphertext too short")
        return nil, errors.New("crypto/cipher: ciphertext too short")
    }

    iv := decryptText[:c.blockSize]
    decryptText = decryptText[c.blockSize:]

    decrypted := src[c.blockSize:]
    var block cipher.Block
    block, err = aes.NewCipher(c.key)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(decrypted, decrypted)
    return decrypted, nil
}