package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type ResultToken struct {
	Valid    bool
	Claims   jwt.StandardClaims
	ErrorMsg string
}

func CheckToken(secret []byte, token string) ResultToken {
	res := ResultToken{}
	t, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		res.Valid = false
		res.ErrorMsg = err.Error()
		return res
	}
	if t.Valid {
		res.Valid = true
		if claims, ok := t.Claims.(jwt.StandardClaims); ok {
			res.Claims = claims
		} else {
			fmt.Println(err)
		}

	} else if ve, ok := err.(*jwt.ValidationError); ok {
		res.Valid = false
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			res.Valid = false
			res.ErrorMsg = fmt.Sprintf("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			res.Valid = false
			res.ErrorMsg = fmt.Sprintf("Timing is everything")
		} else {
			res.Valid = false
			res.ErrorMsg = fmt.Sprintf("Couldn't handle this token:%v", err)
		}
	} else {
		res.Valid = false
		res.ErrorMsg = fmt.Sprintf("Couldn't handle this token:%v", err)
	}
	return res
}

func GetTokenInfo(token string) (string, error) {
	mySigningKey := []byte("hzwy23")
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(*jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		return "", errors.New("wrong")
	}
	claims := t.Claims.(*jwt.StandardClaims)
	return claims.Id, nil
}
