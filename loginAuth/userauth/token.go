// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package userauth provides functions used for user login authentication.
*/
package userauth

import (
	"fmt"
	"time"

	"github.com/jeanphorn/log4go"
	"github.com/dgrijalva/jwt-go"
)

var Secret string = "hdfkhfkajhfkjadshkl"

//Claims custome token
type Claims struct{
	UserName string //the name of the login user
	jwt.StandardClaims
}

//CreateToken create token
func CreateToken(claims *Claims) (signedToken string, success bool) {
	claims.ExpiresAt = time.Now().Add(time.Minute*30).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	signedToken, err := token.SignedString([]byte(Secret))
	if err != nil{
		log4go.Error(err)
	}
	success=true
	return
}

//CheckToken checks token's validation
func CheckToken(tokenString string) (err error){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		// validate alg is what I expected
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// Secret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return Secret, nil
	})
	if token.Valid {
		log4go.Info("Succeeded in parsing the token")
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			log4go.Error("This is not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			log4go.Error("Token has expired or is inactive")
			return err
		} else {
			log4go.Error("Couldn't handle this token")
			return err
		}
	} else {
		log4go.Error("Couldn't handle this token")
		return err
	}
	return nil
}
