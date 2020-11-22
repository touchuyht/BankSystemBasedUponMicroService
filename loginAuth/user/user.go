// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
The user package contains user struct and some functions
*/
package user

type User struct{
	userName string
	password string
	token string
}

func NewUser(userName, password, token string) User{
	return User{
		userName: userName,
		password: password,
		token: token,
	}
}

func (u *User)GetUserName() string {
	return u.userName
}

func (u *User)GetPassword() string {
	return u.password
}

func (u *User)GetToken() string {
	return u.token
}

func (u *User)SetUserName(userName string){
	u.userName = userName
}

func (u *User)SetPassword(password string){
	u.password = password
}

func (u *User)SetToken(token string){
	u.token = token
}