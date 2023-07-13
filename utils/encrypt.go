package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
	"log"
	"math/rand"
	"time"
)

type _encrypt struct{}

var Encryptor = new(_encrypt)

// ScryptHash 使用 scrypt 对密码进行加密生成一个哈希值
// 有关具体的加密操作还是要看官方文档
func (*_encrypt) ScryptHash(password string) string {
	const KeyLen = 10
	salt := []byte{12, 32, 4, 6, 66, 22, 222, 11} // 这里是八个数字

	hashPwd, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, KeyLen)
	if err != nil {
		log.Fatal("加密失败：", err)
	}
	return base64.StdEncoding.EncodeToString(hashPwd)
}

// ScryptCheck 使用 scrypt 对比 明文密码 和 数据库中的哈希值
func (c *_encrypt) ScryptCheck(password, hash string) bool {
	return c.ScryptHash(password) == hash // 在这里直接判断并返回，md，我以后也要这样装逼
}

// BcryptHash 使用 bcrypt 对密码进行加密生成一个哈希值
func (*_encrypt) BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 使用 bcrypt 对比 明文密码 和 数据库中哈希值
// 这里函数名之前的接受者的变量名是否必要取决于是否要使用指针的方法，所以指针名不是必须的
func (*_encrypt) BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // 这里 bcrypt 比较出来的函数结果是一个err，判断是否是空值
}

// MD5 加密
func (*_encrypt) MD5(str string, b ...byte) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(b))
}

// ValidateCode 验证码
func (*_encrypt) ValidateCode() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

// UUID TODO
func UUID() string {
	return ""
}
