package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// GenerateRsaPrivateKey 使用私钥生成 *rsa.PrivateKey 实例
func GenerateRsaPrivateKey(fileName string) (rsaPrivateKey *rsa.PrivateKey, err error) {
	// 解码私钥
	block, _ := pem.Decode(GetKeyByte(fileName))
	rsaPrivateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("格式化私钥错误：", err)
		return nil, err
	}
	return
}

// GenerateRsaPublicKey 使用公钥生成 *rsa.PublicKey 实例
func GenerateRsaPublicKey(fileName string) (rsaPublicKey *rsa.PublicKey, err error) {
	// 解码公钥
	block, _ := pem.Decode(GetKeyByte(fileName))
	rsaPublicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		fmt.Println("格式化公钥错误：", err)
		return nil, err
	}
	return
}

// Signature 使用私钥签名
func (r *RSA) Signature(message []byte, fileName string) ([]byte, error) {
	// 只有小消息可以直接签名； 因此，对消息的哈希进行签名，而不能对消息本身进行签名。
	// 这要求哈希函数必须具有抗冲突性。 SHA-256是编写本文时(2016年)应使用的最低强度的哈希函数。
	hashed := sha256.Sum256(message)
	rsaPrivateKey, _ := GenerateRsaPrivateKey(fileName)

	// ######################
	// ######## 签名 ########
	// ######################
	signature, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return nil, err
	}
	return signature, nil

}

// Verifying 使用公钥验证签名
func (r *RSA) Verifying(signedMessage []byte, message []byte, fileName string) bool {
	// 只有小消息可以直接签名； 因此，对消息的哈希而不是消息本身进行签名。
	// 这要求哈希函数必须具有抗冲突性。 SHA-256是编写本文时(2016年)应使用的最低强度的哈希函数。
	hashed := sha256.Sum256(message)
	rsaPublicKey, _ := GenerateRsaPublicKey(fileName)

	// #########################
	// ######## 验证签名 ########
	// #########################
	err := rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed[:], signedMessage)
	if err != nil {
		fmt.Printf("验证失败：%v\n", err)
		return false
	}
	return true
}

// SignatureAndVerifying 使用私钥签名，公钥验签
func SignatureAndVerifying(message []byte, r *RSA) {
	// 使用指定的私钥进行签名
	signedMessage, _ := r.Signature(message, "./practice/cryptography/private.pem")
	fmt.Printf("已签名的消息为: %x\n", signedMessage)
	// 验证签名的数据
	ok := r.Verifying(message, signedMessage, "./practice/cryptography/public.pem")
	if ok == true {
		fmt.Println("验证成功")
	}
}
