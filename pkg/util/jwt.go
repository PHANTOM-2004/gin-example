package util

import (
	"gin-example/pkg/setting"
	"time"

	jwt "github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	log.Debugf("generate token with [%s]", string(jwtSecret))

	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(
		token,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

	if tokenClaims == nil {
		return nil, err
	}

	claims, ok := tokenClaims.Claims.(*Claims)
	if ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, err
}
