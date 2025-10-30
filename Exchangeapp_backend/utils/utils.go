package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)
const (
	TOKEN_NAME   = "Authorization"
	TOKEN_PREFIX = "Bearer "
)
//对用户传过来的密码进行加密 
func HashPassword(pwd string) (string, error) {
	//用于对密码做不可逆哈希（带随机盐，不能解密）
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12 /** bcrypt.MinCost 为对字符串进行哈希的次数*/) //使用bcrypt对密码进行加密
	return string(hash), err
}
//对用户传过来的密码进行校验
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//JWT：用于生成带签名的访问令牌，客户端携带，服务端验签与校验过期，无需存会话。
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

// ParseJWT 解析并校验 JWT：
// 1) 去掉 Bearer 前缀；2) 验签；3) 校验标准声明（exp/nbf/iat）；4) 返回 username
func ParseJWT(tokenString string) (string, error) {
    if len(tokenString) > 7 && tokenString[:7] == TOKEN_PREFIX {
        tokenString = tokenString[7:]
    }

    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected Signing Method")
        }
        return []byte(viper.GetString("jwt.signingKey")), nil // 签名密钥需与签发时一致
    })
    if err != nil {
        return "", err
    }
	//token.Valid：只代表“签名方式正确、格式OK”。不一定检查到期时间。
    if !token.Valid {
        return "", errors.New("invalid token")
    }
	fmt.Printf("全部claims内容：%#v\n", claims) 
    // 打印过期时间（如存在）
    if expVal, ok := claims["exp"]; ok {
        var expUnix int64
        switch v := expVal.(type) {
        case float64:
            expUnix = int64(v)
        case int64:
            expUnix = v
        }
        if expUnix > 0 {
			fmt.Println("token 过期时间:", time.Unix(expUnix, 0).Format(time.RFC3339))
        }
    }
    // claims.Valid()：专门检查“时间相关”的有效性，比如是否过期(exp)、是否未到生效时间(nbf)、签发时间是否合理(iat)。
    if err := claims.Valid(); err != nil {
        return "", err
    }
    username, _ := claims["username"].(string)
    if username == "" {
        return "", errors.New("username claim is not a string")
    }
    return username, nil
}

// 去掉转义字符
func RemoveEscapeChars(s string) string {
	return regexp.MustCompile(`\\(.)`).ReplaceAllStringFunc(s, func(m string) string {
		return string([]byte(m)[1:])
	})
}

// 去掉转义字符
func RemoveEscapeChars1(s string) string {
	Logger.Errorf("测试 %t", strings.Contains(s, "\\"))
	return strings.ReplaceAll(s, "\\", "") // \"ID\":9 变成 "ID":9
}
func IsNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func GetDuration(Hour int, Min int, Sec int) time.Duration {
	// 获取当前时间
	now := time.Now()
	Logger.Println("现在时间", now)
	// 计算今天20点的时间
	target := time.Date(now.Year(), now.Month(), now.Day(), Hour, Min, Sec, 0, now.Location())
	Logger.Println("计算今天6点的时间", target)
	// 如果当前时间已经是下午6点之后，则计算明天6点
	if now.After(target) {
		target = target.AddDate(0, 0, 1)
	}

	// 计算倒计时
	duration := target.Sub(now) //time.Sub方法用于计算两个时间点之间的时间差
	Logger.Println("时间差是：", duration)
	return duration
}

// IntAbs 自定义函数，用于求整数的绝对值
func IntAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// removeChineseCharacters 移除字符串中的中文字符
func RemoveChineseCharacters(s string) string {
	s = strings.TrimSpace(s)
	re := regexp.MustCompile(`[\p{Han}]+`)
	return re.ReplaceAllString(s, "")
}

// GenerateRandomValue 生成指定范围内的随机数
// min: 最小值（包含）
// max: 最大值（包含）
func GenerateRandomValue(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}
