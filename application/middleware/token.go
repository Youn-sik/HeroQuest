package middleware

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

func TokenAuthenticate(c *gin.Context) {
	authToken := c.Request.Header.Get("Authorization")

	if authToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"result": false, "errStr": "No Token"})
		c.Abort()
		return
	}

	// token 검증 로직
	isValid := TokenCheck(authToken)
	if isValid {
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"result": false, "errStr": "Expired Token"})
		c.Abort()
		return
	}
}

// 토큰 발급
func TokenBuild(id string) string {
	at := AuthTokenClaims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Minute * 36000)), // 6 hours
		},
	}

	atoken := jwt.NewWithClaims(jwt.SigningMethodHS256, &at)
	signedAuthToken, err := atoken.SignedString([]byte("heroquest"))

	if err != nil {
		log.Println(err)
		return "false"
	}

	return signedAuthToken
}

/*
func Auth(c *gin.Context) {
	var send_data struct {
		result bool
		msg    string
	}
	authToken := c.Request.Header.Get("authorization")
	authToken = strings.Replace(authToken, "Bearer ", "", 1)

	send_data.result = false
	if authToken == "" {
		send_data.msg = "Token is required."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "msg": send_data.msg})
		return
	}

	isValid := TokenCheck(authToken)
	if isValid {
		send_data.result = true
		send_data.msg = "Token is verified."
		c.JSON(http.StatusOK, gin.H{"result": send_data.result, "msg": send_data.msg})
		return
	} else {
		send_data.msg = "Token is ."
		c.JSON(http.StatusUnauthorized, gin.H{"result": "false", "msg": "Expired Token"})
		c.Abort()
		return
	}
}
*/

func TokenCheck(authToken string) bool {
	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			ErrUnexpectedSigningMethod := errors.New("unexpected signing method")
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte("heroquest"), nil
	}

	user := AuthTokenClaims{}
	token, err := jwt.ParseWithClaims(authToken, &user, key)

	if err != nil {
		// token is expired by ...
		log.Println(err)
		return false
	}

	return token.Valid
}

func ExtractClaims(tokenStr string) (bool, jwt.MapClaims) {
	hmacSecret := []byte("heroquest")
	authToken := tokenStr
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		log.Println(err)
		return false, nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, claims
	} else {
		log.Printf("Invalid JWT Token")
		log.Println(err)
		return false, nil
	}
}

func GetIdFromToken(authToken string) (bool, string, string) {
	var token struct {
		Id string
	}

	flag, value := ExtractClaims(authToken)
	if !flag {
		return false, "Token Parsing Error1", ""
	}

	mapstructure.Decode(value, &token)
	if token.Id == "" {
		return false, "Token Parsing Error2", ""
	}
	return true, "", token.Id
}
