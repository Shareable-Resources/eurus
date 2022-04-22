package auth

import (
	"errors"
	"eurus-backend/foundation"
	"eurus-backend/foundation/log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var mySigningKey []byte
var tokenTimeoutSecond int64 = 3600

func SetupJwt(signKey []byte, timeoutSecond int64) {
	mySigningKey = signKey
	tokenTimeoutSecond = timeoutSecond
}

type claims struct {
	jwt.StandardClaims
	ClientInfo string `json:"clientInfo"`
}

func GenerateJwtToken(acc, clientInfo string, hasTimeout bool, second int64) (string, int64, error) {
	if mySigningKey == nil {
		return "", 0, errors.New("JWT not setup")
	}

	var expireTime int64
	if hasTimeout {
		if second > 0 {
			expireTime = time.Now().Unix() + second
		} else {
			expireTime = time.Now().Unix() + tokenTimeoutSecond
		}
	} else {
		expireTime = 0
	}

	// Create the Claims
	jwtClaims := &claims{
		jwt.StandardClaims{
			NotBefore: int64(time.Now().Unix() - tokenTimeoutSecond), //TIMEOUT
			ExpiresAt: int64(expireTime),                             //TIMEOUT
			Issuer:    acc,
		},
		clientInfo,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		log.GetLogger(log.Name.Root).Error("GenerateJwtToken error: ", err)
		return "", 0, err
	}
	return ss, expireTime, nil
}

/// \\return return code, error message, client information JSON
func ValidateJwtoken(token string) (foundation.ServerReturnCode, string, string) {

	if mySigningKey == nil {
		return foundation.UnauthorizedAccess, "JWT not setup", ""
	}

	t, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		var byPass bool = false
		validateErr, ok := err.(*jwt.ValidationError)
		if ok {
			if (validateErr.Errors & jwt.ValidationErrorExpired) > 0 {
				byPass = true
			}
		}
		if !byPass {
			log.GetLogger(log.Name.Root).Error("parase with claims failed.", err)
			var code foundation.ServerReturnCode = foundation.LoginTokenInvalid
			return code, code.String(), ""
		}
	}

	var mapClaims jwt.MapClaims = t.Claims.(jwt.MapClaims)
	// validateErr := mapClaims.Valid()
	// if validateErr != nil {
	// 	actualErr := validateErr.(*jwt.ValidationError)
	// 	if (actualErr.Errors & jwt.ValidationErrorExpired) > 0 {
	// 		return foundation.LoginTokenExired, actualErr.Inner.Error(), ""
	// 	}
	// }

	return foundation.Success, foundation.ServerReturnCode(foundation.Success).String(), mapClaims["clientInfo"].(string)
}
