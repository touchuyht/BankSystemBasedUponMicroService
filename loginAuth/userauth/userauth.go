// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package userauth provides functions used for user login authentication.
Use it in the following order:
	1. Creates a UserAuth type variable;
	2. Call UserAuthInit to initialize the variable with the parameters, note that the message should contain
	the user's information and token;
	3. Call CheckPassword to check password's correctness;
	4.Call CheckToken to check token's validation if token exist's or it will generates a new token and insert it into
	redis and send it to data push and consistency layer.
*/
package userauth

import (
	"encoding/json"
	"fmt"
	"github.com/jeanphorn/log4go"
	"loginAuth/config"
	"loginAuth/redisdao"
	"loginAuth/messagebroke"
	"errors"
	"loginAuth/user"
)

//UserAuth is the struct used for check password and password
type UserAuth struct {
	redisHelper *redisdao.RedisHelper
	messageBroke *messagebroke.MessageBroke
	user user.User
	conf *config.Config
	messageId *uint64
}

func (userAuth *UserAuth)UserAuthInit(redisHelper *redisdao.RedisHelper, messageBroke *messagebroke.MessageBroke,
	conf *config.Config, messageId *uint64){
	userAuth.redisHelper = redisHelper
	userAuth.messageBroke = messageBroke
	userAuth.conf = conf
	userAuth.messageId = messageId
}

func (userAuth *UserAuth)Check() (err error){
	if userAuth.redisHelper==nil||userAuth.messageBroke==nil{
		log4go.Error("UserAuth attributes are nil")
		err = errors.New("UserAuth attributes are nil")
	}
	return
}

func (userAuth *UserAuth)CheckPassword(username, password string)(result string, err error){
	if err=userAuth.Check();err!=nil{
		log4go.Info("No username or password found")
		return "No username or password found", err
	}
	userAuth.user, err = userAuth.redisHelper.GetUser(username)
	if err!=nil{
		log4go.Info("No such user")
		return "user unregistered", err
	}
	if userAuth.user.GetPassword()!=password{
		log4go.Info("Password is incorrect")
		return "password incorrect", errors.New("Password incorrect")
	}
	log4go.Info(fmt.Printf("%v password check password", username))
	return "password check passed", nil
}

func (userAuth *UserAuth)CheckToken(username string)(result string, err error){
	if err=userAuth.Check();err!=nil{
		log4go.Info("No username or password found")
		return "No username or password found", err
	}
	if userAuth.user.GetToken()!=""{
		//There is already a token of the user in redis
		token := userAuth.user.GetToken()
		if err=CheckToken(token);err!=nil{
			return userAuth.dealToken(username)
		}else{
			return "token="+token, nil
		}
	}else{
		//generates a new token
		return userAuth.dealToken(username)
	}
}

func (userAuth *UserAuth)dealToken(username string)(result string, err error){
	//generates a new token and sends a message to messaageBroke
	//inserts token into redis, need to distinguish update and insert
	claims := Claims{UserName: username}
	token, success := CreateToken(&claims)
	if success{
		userAuth.user.SetToken(token)
		err = userAuth.redisHelper.InsertToken(userAuth.user)
		if err!=nil{
			s := "Fail to insert token into redis"
			log4go.Error(s)
			return s, err
		}
		(*userAuth.messageId)+=1
		userAuth.conf.UpdateToken.MessageId=(*userAuth.messageId)
		userAuth.conf.UpdateToken.UserName=userAuth.user.GetUserName()
		userAuth.conf.UpdateToken.Password=userAuth.user.GetPassword()
		userAuth.conf.UpdateToken.Token=userAuth.user.GetPassword()
		bytes, err := json.Marshal(userAuth.conf.UpdateToken)
		if err!=nil{
			s := "Fail to marshal UpdateToken into json"
			log4go.Error(s)
			return s, err
		}
		userAuth.messageBroke.SendMessage(string(bytes))
		return "token="+token, nil
	}else{
		log4go.Error("fail to generate token")
		return "fail to generate token", errors.New("fail to generate token")
	}
}