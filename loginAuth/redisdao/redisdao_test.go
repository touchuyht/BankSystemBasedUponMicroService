package redisdao

import (
	"loginAuth/user"
	"loginAuth/userauth"
	"testing"
)

func TestRedisHelper_ConnectRedis(t *testing.T) {
	rd := RedisHelper{}
	rd.SetProtocol("tcp")
	rd.SetIp("localhost")
	rd.SetPort("6379")
	err := rd.ConnectRedis()
	if err != nil{
		t.Error("Failed to connect to redis")
	}
	rd.Close()
}

func TestRedisHelper_GetUser(t *testing.T) {
	rd := RedisHelper{}
	rd.SetProtocol("tcp")
	rd.SetIp("localhost")
	rd.SetPort("6379")
	err := rd.ConnectRedis()
	if err != nil{
		t.Error("Failed to connect to redis")
	}
	user, err := rd.GetUser("zhangshan")
	if user.GetPassword() != "123456"{
		t.Error("Failed to fetch user info from redis")
	}
	rd.Close()
}

func TestRedisHelper_InsertToken(t *testing.T) {
	rd := RedisHelper{}
	rd.SetProtocol("tcp")
	rd.SetIp("localhost")
	rd.SetPort("6379")
	err := rd.ConnectRedis()
	if err != nil{
		t.Error("Failed to connect to redis")
	}
	claims := userauth.Claims{UserName: "zhangshan"}
	token, success := userauth.CreateToken(&claims)
	if !success {
		t.Error("Failed to create token")
	}
	user := user.User{}
	user.SetUserName("zhangshan")
	user.SetPassword("123456")
	user.SetToken(token)
	err = rd.InsertToken(user)
	if err!=nil{
		t.Error("Failed to insert token into redis")
	}
	rd.Close()
}

func TestRedisHelper_UpdateToken(t *testing.T) {
	rd := RedisHelper{}
	rd.SetProtocol("tcp")
	rd.SetIp("localhost")
	rd.SetPort("6379")
	err := rd.ConnectRedis()
	if err != nil{
		t.Error("Failed to connect to redis")
	}
	claims := userauth.Claims{UserName: "zhangshan"}
	token, success := userauth.CreateToken(&claims)
	if !success {
		t.Error("Failed to create token")
	}
	user := user.User{}
	user.SetUserName("zhangshan")
	user.SetPassword("123456")
	user.SetToken(token)
	err = rd.UpdateToken(user)
	if err!=nil{
		t.Error("Failed to insert token into redis")
	}
	rd.Close()
}