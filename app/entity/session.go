package entity

import (
	"crypto/rand"
	"fmt"
	"time"
)

type Session struct {
	Id        int64
	Uuid      string
	UserId    int64
	CreatedAt time.Time
}

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func CreateUuid() (uuid string, err error) {
	u := new([16]byte)
	_, err = rand.Read(u[:])
	if err != nil {
		return
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7f
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xf) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}
