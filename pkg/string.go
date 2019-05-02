package pkg

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandString 产生定长的随机字符串
// refrence https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func CombineOSSURL(endpoint, bucket string) (string, error) {
	URL, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}
	if URL.Host == "" || URL.Scheme == "" {
		return "", fmt.Errorf("Invalid endpoint: %s", endpoint)
	}
	return fmt.Sprintf("%s://%s.%s", URL.Scheme, bucket, URL.Host), nil
}

func PrimaryKey() string {
	return xid.New().String()
}

// EncryptPassword encrypt password
func EncryptPassword(password string) (string, error) {

	var hashedPassword []byte
	var err error

	// Hashing the password with the default cost of 10
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compare password
func ComparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
