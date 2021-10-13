package entity

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id   int64
	Name string
	// ユーザが任意に設定可能なID
	AccountId string
	// 好きなタバコの銘柄
	FavoriteBrand string
	Password      string
}

func EncryptPassword(plain string) (string, error) {
	// 2^10 の強度
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPasswordHash(hashed, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err != bcrypt.ErrMismatchedHashAndPassword
}

const accountIdPattern = `[a-zA-Z0-9\_]+`

func VerifyAccountId(accountId string) bool {
	if len(accountId) > 255 {
		return false
	}
	reg := regexp.MustCompile(accountIdPattern)
	return reg.MatchString(accountId)
}

const passwordPattern = `[a-zA-Z0-9]{8,}`

func VerifyPlainPassword(plain string) bool {
	if len(plain) > 255 {
		return false
	}
	reg := regexp.MustCompile(passwordPattern)
	return reg.MatchString(plain)
}

func VerifyAccountName(name string) bool {
	l := len(name)
	return 0 < l && l < 256
}

func VerifyFavoriteBrand(brand string) bool {
	l := len(brand)
	return l < 256
}
