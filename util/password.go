package util

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func MD5(str, salt string) string {
	m := md5.New()
	m.Write([]byte(str))
	if salt != "" {
		m.Write([]byte(salt))
	}
	return hex.EncodeToString(m.Sum(nil))
}

func MD5Short(str, salt string) string {
	return MD5(str, salt)[8:24]
}

func Sha256(str, salt string) string {
	s := sha256.New()
	s.Write([]byte(str))
	if salt != "" {
		s.Write([]byte(salt))
	}
	return hex.EncodeToString(s.Sum(nil))
}

func Sha512(str, salt string) string {
	s := sha512.New()
	s.Write([]byte(str))
	if salt != "" {
		s.Write([]byte(salt))
	}
	return hex.EncodeToString(s.Sum(nil))
}
