package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ //组合了head和payload
		"username": username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	signedToken, err := token.SignedString([]byte("secret")) //secret是密钥，SignedString通过特定公式对token加密
	return "Bearer " + signedToken, err
}

func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ParseJWT(tokenString string) (string, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//执行parse时base64uri解密head发现alg是hs256，于是通过new或者取地址符，Method得到指针实例*jwt.SigningMethodHMAC，为什么能加，见以下注释
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { //Method是SigningMethod 接口类型变量，接口里有verify sign alg三个方法
			return nil, errors.New("Unexpected Signing Method") //jwt包里的指向SigningMethodHMAC结构体的指针类型 实现了verify sign alg三个方法，可添加，可断言
			//补充：Go是静态强类型语言，只有接口变量是特殊的，当一个类型实现了接口的所有方法，他的实例可以被赋值到接口变量
		}
		return []byte("secret"), nil
		//有了key之后，结合head和payload，哈希出signature，对比前端传来的signature
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { //Claims静态类型是接口，动态类型是map
		usernaem, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("Username claim is not a string")
		}
		return usernaem, nil
	}
	return "", err
}
