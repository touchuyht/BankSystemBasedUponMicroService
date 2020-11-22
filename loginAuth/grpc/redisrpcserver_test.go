// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grpc

import (
	"fmt"
	"loginAuth/user"
	"loginAuth/userauth"
	"testing"
)

func TestRedisRPCServer_UpdateToken(t *testing.T) {
	var c chan int = make(chan int)
	redisRPCServer := &RedisRPCServer{}
	redisRPCServer.RedisRPCServerInit("localhost:2100","tcp","x509/server_cert.pem", "x509/server_key.pem")
	go redisRPCServer.StartRPCServer()
	redisRPCClient := &RedisRPCClient{}
	redisRPCClient.RedisRPCClientInit()
	err := redisRPCClient.StartRPCClient()
	if err != nil{
		t.Error("fail to start rpc client")
	}
	user := user.User{}
	user.SetUserName("zhangshan")
	user.SetPassword("123456")
	claims := &userauth.Claims{UserName: user.GetUserName()}
	token, success := userauth.CreateToken(claims)
	if !success {
		t.Error("Wrong token")
	}
	user.SetToken(token)
	//redisRPCClient.InsertUsers([]user.user{user})
	redisRPCClient.UpdateToken(user)
	go func(){
		_user, err := redisHelper.GetUser("zhangshan")
		if err != nil{
			t.Error("fail to get user info")
		}
		if _user.GetToken() != user.GetToken(){
			t.Error(fmt.Sprintf("Token mismatch: Got %v, expect %v", _user.GetToken(), user.GetToken()))
		}
		c<-1
	}()
	<-c
}

func TestRedisRPCServer_InsertUsers2(t *testing.T) {
	var c chan int = make(chan int)
	redisRPCServer := &RedisRPCServer{}
	redisRPCServer.RedisRPCServerInit("localhost:2100","tcp","x509/server_cert.pem", "x509/server_key.pem")
	go redisRPCServer.StartRPCServer()
	redisRPCClient := &RedisRPCClient{}
	redisRPCClient.RedisRPCClientInit()
	err := redisRPCClient.StartRPCClient()
	if err != nil{
		t.Error("fail to start rpc client")
	}
	user := make([]user.User, 10)
	for i:=0; i< 10; i++{
		user[i].SetUserName("zhangshan"+fmt.Sprintf("%d", i))
		user[i].SetPassword("123456")
		user[i].SetToken("")
	}
	redisRPCClient.InsertUsers(user)
	go func(){
		for _, v := range user{
			_user, err := redisHelper.GetUser(v.GetUserName())
			if err != nil{
				t.Error("fail to get user info")
			}
			if _user.GetToken() != v.GetToken() || _user.GetPassword()!=v.GetPassword(){
				t.Error(fmt.Sprintf("user mismatch"))
			}
		}
		c<-1
	}()
	<-c
}