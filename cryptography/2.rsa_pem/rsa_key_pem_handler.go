package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"os"

	"crypto/rsa"

	"encoding/pem"
)

// RsaKey 是公钥和私钥两个组成一组的密钥对的二进制格式。可以通过二进制转换为文件
type RsaKey struct {
	bytePrivateKey []byte
	bytePublicKey  []byte
}

// NewRsaKey 生成密钥对
func NewRsaKey(bits int) *RsaKey {
	// 随机生成一个给定大小的 RSA 密钥对。可以使用 crypto 包中的 rand.Reader 来随机。
	rsaPrivateKey, _ := rsa.GenerateKey(rand.Reader, bits)
	// 从私钥中，获取公钥
	rsaPublicKey := rsaPrivateKey.PublicKey

	// ======================================================
	// ======== 比创建基本 RSA 密钥对多出来的行为，开始 ========
	// ======================================================
	// 将密钥转为二进制流，以便使用 PEM 包将其编码。
	// bytePrivateKey := rsaPrivateKey.D.Bytes()
	// bytePublicKey := rsaPublicKey.N.Bytes()
	// 为什么不能直接转换而必须使用 X.509 呢，因为 go 的标准库中，无法直接将 []byte 格式的数据转换为 *rsa.PrivateKey 和 *rsa.PublicKey 类型
	// 这样会导致后面在使用 rsa 加密/解密，签名/验签时，无法正确获取使用 *rsa.PrivateKey 和 *rsa.PublicKey 类型的密钥。
	// 因为 go 没有函数或方法，可以将密钥从 []byte 类型转换为 *rsa.PrivateKey 或 *rsa.PublicKey 类型
	// =================================================================================================================
	// ！！！！！注意！！！！！在正常情况下，无需将密钥先编码为 DER 格式。这是 go 语言的强制要求。*****也是 PKCS 这个标准搞的*****
	// go 在处理密钥时，需要先将密钥先转换成 DER 格式再使用 PEM 编码的，是 go 的需求！！！！！！！！！！！！！！！！注意！！！！！！！！！！！！！！！！！！！！！！
	// ==================================================================================================================
	// 因此，只能先使用 x509 包中的方法，将密钥对转换为 PKCS#1,ASN.1 DER 的形式，并以 []byte 的数据类型保存，以供 PEM 包将其编码。
	bytePrivateKey := x509.MarshalPKCS1PrivateKey(rsaPrivateKey)
	bytePublicKey := x509.MarshalPKCS1PublicKey(&rsaPublicKey)
	// 将密钥对转换为 PEM 格式，并在密钥对中添加用于标识类型的页眉与页脚
	// 1.首先，声明两个 bytes.Buffer 类型的变量，用于存放编码后的 PEM 格式密钥。
	var (
		bufPemPrivateKey bytes.Buffer
		bufPemPublicKey  bytes.Buffer
	)
	// 2.然后，使用 pem 包将密钥对编码为 PEM 格式的数据
	// 因为 bytes.Buffer 这个类型的结构体实现了 pem.Encode 第一个参数的 io.Writer 接口，所以可以通过 pem 包，将 PEM 的标签和二进制类型的编码内容，再编码为 PEM 格式的数据
	pem.Encode(&bufPemPrivateKey, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: bytePrivateKey,
	})
	pem.Encode(&bufPemPublicKey, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: bytePublicKey,
	})
	// 这里就可以看到平时看到的带页眉页脚的 PEM 格式的编码后的密钥内容了。
	fmt.Printf("======== PEM 格式私钥内容：========\n%s", bufPemPrivateKey.String())
	fmt.Printf("======== PEM 格式公钥内容：========\n%s", bufPemPublicKey.String())
	// =====================================================
	// ======== 比创建基本 RSA 密钥对多出来的行为，结束 =======
	// =====================================================

	return &RsaKey{
		bytePrivateKey: bufPemPrivateKey.Bytes(), // 返回 PEM 格式的二进制私钥
		bytePublicKey:  bufPemPublicKey.Bytes(),  // 返回 PEM 格式的二进制公钥
	}
}

// RsaPemEncrypt 使用 RSA 算法，加密指定明文，其中私钥是 PEM 编码后的格式
func (r *RsaKey) RsaPemEncrypt(plaintext []byte) []byte {
	// 由于这次要通过 PEM 格式编码的公钥进行加密，所以需要先解码 PEM 格式，再将解码后的数据转换为 *rsa.PublicKey 类型
	block, _ := pem.Decode(r.bytePublicKey)
	// 之前在编码时，使用了 x509 进行了编码，所以同样，需要使用 x509 解码以获得 *rsa.PublicKey 类型的公钥
	rsaPublicKey, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	// 使用公钥加密 plaintext(明文，也就是准备加密的消息)。并返回 ciphertext(密文)
	// 其中 []byte("DesistDaydream") 是加密中的标签，解密时标签需与加密时的标签相同，否则解密失败
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPublicKey, plaintext, []byte("DesistDaydream"))
	if err != nil {
		panic(err)
	}
	return ciphertext
}

// RsaPemDecrypt 使用 RSA 算法，解密指定密文，其中公钥是 PEM 编码后的格式
func (r *RsaKey) RsaPemDecrypt(ciphertext []byte) []byte {
	// 使用私钥解密 ciphertext(密文，也就是加过密的消息)。并返回 plaintext(明文)
	// 其中 []byte("DesistDaydream") 是加密中的标签，解密时标签需与加密时的标签相同，否则解密失败
	block, _ := pem.Decode(r.bytePrivateKey)
	rsaPrivateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, ciphertext, []byte("DesistDaydream"))
	if err != nil {
		panic(err)
	}
	return plaintext
}

// RsaPemSign RSA 签名
func (r *RsaKey) RsaPemSign(plaintext []byte) []byte {
	// 只有小消息可以直接签名； 因此，对消息的哈希进行签名，而不能对消息本身进行签名。
	// 这要求哈希函数必须具有抗冲突性。 SHA-256是编写本文时(2016年)应使用的最低强度的哈希函数。
	hashed := sha256.Sum256(plaintext)
	// 由于这次要通过 PEM 格式编码的公钥进行加密，所以需要先解码 PEM 格式，再将解码后的数据转换为 *rsa.PublicKey 类型
	block, _ := pem.Decode(r.bytePrivateKey)
	// 之前在编码时，使用了 x509 进行了编码，所以同样，需要使用 x509 解码以获得 *rsa.PublicKey 类型的公钥
	rsaPrivateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	// 使用私钥签名，必须要将明文hash后才可以签名，当验证时，同样需要对明文进行hash运算。签名与验签并不用于加密消息或消息传递，仅仅作为验证传递消息方的真实性。
	signed, err := rsa.SignPKCS1v15(rand.Reader, rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from signing: %s\n", err)
		return nil
	}
	fmt.Printf("已签名的消息为: %x\n", signed)
	return signed
}

// RsaPemVerify RSA 验签
func (r *RsaKey) RsaPemVerify(plaintext []byte, signed []byte) bool {
	// 与签名一样，只可以对 hash 后的消息进行验证。
	hashed := sha256.Sum256(plaintext)
	block, _ := pem.Decode(r.bytePublicKey)
	rsaPublicKey, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	// 使用公钥、已签名的信息，验证签名的真实性
	err := rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hashed[:], signed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from verification: %s\n", err)
		return false
	}
	return true
}
