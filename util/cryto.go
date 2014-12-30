package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

func MD5(str, salt string) string {
	m := md5.New()
	m.Write([]byte(str))
	if len(salt) > 0 {
		m.Write([]byte(salt))
	}
	return hex.EncodeToString(m.Sum(nil))
}

func MD5Short(str, salt string) string {
	return MD5(str, salt)[8:24]
}

func Sha1(str, salt string) string {
	s := sha1.New()
	s.Write([]byte(str))
	if len(salt) > 0 {
		s.Write([]byte(salt))
	}
	return hex.EncodeToString(s.Sum(nil))
}

func Sha256(str, salt string) string {
	s := sha256.New()
	s.Write([]byte(str))
	if len(salt) > 0 {
		s.Write([]byte(salt))
	}
	return hex.EncodeToString(s.Sum(nil))
}
