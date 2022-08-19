package user

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

// 토큰 발급
func TokenBuild(u User) string {
	at := AuthTokenClaims{
		ID: u.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Minute * 360)), // 6 hours
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

func TokenCheck(authToken string) bool {
	key := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			ErrUnexpectedSigningMethod := errors.New("unexpected signing method")
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte("cho"), nil
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
