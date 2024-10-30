package utils

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12 /** bcrypt.MinCost 为对字符串进行哈希的次数*/) //使用bcrypt对密码进行加密
	return string(hash), err
}

const (
	TOKEN_NAME   = "Authorization"
	TOKEN_PREFIX = "Bearer "
)

func GenerateJWT(username string) (string, error) {

	jwtCustomClaims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), //过期时间72小时
	}

	//jwt.NewWithClaims签名生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtCustomClaims)
	signedToken, err := token.SignedString([]byte(viper.GetString("jwt.signingKey"))) //viper.GetString("jwt.signingKey") viper可以直接取值
	return TOKEN_PREFIX + signedToken, err
}

func ParseJWT(tokenString string) (string, error) {
	if len(tokenString) > 7 && tokenString[:7] == TOKEN_PREFIX {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected Signing Method")
		}
		return []byte(viper.GetString("jwt.signingKey")), nil //这里的密钥一定要和加密的时候传的一样
	})

	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { //Valid有效的
		username, ok := claims["username"].(string) //解析出来拿到 username
		if !ok {
			return "", errors.New("username claim is not a string")
		}
		return username, nil
	}

	return "", err
}
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
