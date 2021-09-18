package entity

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id   int64
	Name string
	// ユーザが任意に設定可能なID
	AccountId string
	Password  string
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
