package urlsigner

import (
	"fmt"
	"strings"
	"time"

	goalone "github.com/bwmarrin/go-alone"
)

type Signer struct {
	Secrect []byte `json:"secrect"`
}

func (s *Signer) GenerateTokenFromString(data string) string {
	var urlToSign string

	crypt := goalone.New(s.Secrect, goalone.Timestamp)
	if strings.Contains(data, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSign = fmt.Sprintf("%shash=", data)
	}

	tokenBytes := crypt.Sign([]byte(urlToSign))
	token := string(tokenBytes)
	return token
}

func (s *Signer) VerifyToken(token string) bool {
	crypt := goalone.New(s.Secrect, goalone.Timestamp)
	_, err := crypt.Unsign([]byte(token))
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (s *Signer) Expired(token string, minutesUntilExpire int) bool {
	crypt := goalone.New(s.Secrect, goalone.Timestamp)
	ts := crypt.Parse([]byte(token))
	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire)*time.Minute
}
