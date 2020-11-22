// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
The redisdao package contains functions which can be used for querying or inserting or deleting userinfo,password and token
*/
package redisdao

import (
	"errors"
	"fmt"
	"sync"
	"loginAuth/user"
	"github.com/gomodule/redigo/redis"
	"github.com/jeanphorn/log4go"
)

var lock = &sync.Mutex{}

//RedisHelper is the only database connection that should exist
//var RedisHelperInstance *RedisHelper

//RedisHelper is the type which you can use to manipulate redis database.
//Use it with the following order: GetInstance->Operations->Close
type RedisHelper struct {
	conn redis.Conn
	protocol string //the protocol used for data transmission
	ip string //the redis server address
	port string //the redis server port
}

//GetInstance returns the RedisHelper instance if it exists and when it's nil, it creates a RedisHelper instance and
//returns it
/*func GetInstance() (*RedisHelper, error){
	if RedisHelperInstance == nil{
		lock.Lock()
		defer lock.Lock()
		if RedisHelperInstance == nil{
			RedisHelperInstance = newRedisHelper()
			err := RedisHelperInstance.ConnectRedis()
			if err != nil{
				log4go.Error(fmt.Sprintf("Fail to connect to redis, err: %v",err))
				return nil, err
			}
		}else{
			log4go.Info("RedisHelper Instance already created")
		}
	}else{
		log4go.Info("RedisHelper Instance already created")
	}
	return RedisHelperInstance,nil
}*/

//ConnectRedis can connect to redis with the given ip and port using the specified protocol
func (rd *RedisHelper)ConnectRedis() (err error) {
	if rd==nil{
		log4go.Error("RedisHelper instance is nil")
		return errors.New("Nil point receivcer")
	}
	rd.conn, err = redis.Dial( rd.protocol,rd.ip+":"+rd.port)
	if err != nil{
		log4go.Error(err)
	}
	return
}

//GetUser returns a user.user object with password and token if there was a token
func (rd *RedisHelper)GetUser(userName string) (_user user.User,err error) {
	if rd==nil{
		log4go.Error("RedisHelper instance is nil")
		return user.User{}, errors.New("Nil point receivcer")
	}
	result, err := redis.Strings(rd.conn.Do("LRANGE", userName, 0, -1))
	if err != nil{
		log4go.Error(err)
		return
	}
	if len(result) == 0{
		log4go.Error(fmt.Sprintf("No such user %v in redis", userName))
		return user.User{}, errors.New(fmt.Sprintf("No such user %v in redis", userName))
	}
	_user.SetUserName(userName)
	_user.SetPassword(result[0])
	if len(result)>1 {
		_user.SetToken(result[1])
	}
	return
}

//InsertUser inserts a user into redis. And if the user already exists, it deletes the old user info, and inserts
//the new info
func (rd *RedisHelper)InsertUser(user user.User)(err error){
	if rd==nil{
		log4go.Error("RedisHelper instance is nil")
		return errors.New("Nil point receivcer")
	}
	_, err = rd.conn.Do("DEL", user.GetUserName())
	if err != nil{
		log4go.Error(fmt.Sprintf("Fail to insert user, err: %v", err))
		return
	}
	if user.GetToken()=="" {
		_, err = rd.conn.Do("RPUSH", user.GetUserName(), user.GetPassword())
	}else{
		_, err = rd.conn.Do("RPUSH", user.GetUserName(), user.GetPassword(), user.GetToken())
	}
	if err != nil{
		log4go.Error(fmt.Sprintf("Fail to insert user, err: %v", err))
	}
	return
}

//InsertToken write the user.user objects' token attribute into redis
func (rd *RedisHelper)InsertToken(user user.User) (err error) {
	if rd==nil{
		log4go.Error("RedisHelper instance is nil")
		return errors.New("Nil point receivcer")
	}
	result, err := redis.Strings(rd.conn.Do("LRANGE", user.GetUserName(), 0, -1))
	if err!=nil {
		log4go.Error(err)
		return
	}
	if len(result) == 0{
		log4go.Error(fmt.Sprintf("No such user %v in redis", user.GetUserName()))
		return
	}
	if len(result) >= 1{
		_, err = rd.conn.Do("LTRIM", user.GetUserName(), 0, 0)
	}
	if user.GetPassword() != result[0]{
		log4go.Error(fmt.Sprintf("Fail to insert token due to wrong password," +
			" user %v", user.GetUserName()))
		return
	}
	_, err = rd.conn.Do("RPUSH", user.GetUserName(), user.GetToken())
	if err!=nil {
		log4go.Error(err)
	}
	return
}

//UpdateToken updates the user.user object's token in redis
func (rd *RedisHelper)UpdateToken(user user.User) (err error){
	if rd==nil{
		log4go.Error("RedisHelper instance is nil")
		return errors.New("Nil point receivcer")
	}
	result, err := redis.Strings(rd.conn.Do("LRANGE", user.GetUserName(), 0, -1))
	if err!=nil{
		log4go.Error(err)
		return
	}
	if len(result) == 0{
		log4go.Error(fmt.Sprintf("No such user %v in redis", user.GetUserName()))
		return
	}
	if len(result) >= 2{
		_, err = rd.conn.Do("LTRIM", user.GetUserName(), 0, 1)
	}
	if user.GetPassword() != result[0]{
		log4go.Error(fmt.Sprintf("Fail to insert token due to wrong password," +
			" user %v", user.GetUserName()))
		return
	}
	if len(result) == 1{
		_, err = rd.conn.Do("RPOP", user.GetUserName())
		if err != nil {
			log4go.Error(err)
			return
		}
	}
	return rd.InsertToken(user)
}

//Close disconnect redis
func (rd *RedisHelper)Close() (err error) {
	if rd==nil{
		log4go.Error("RedisHelper instance is nil")
		return errors.New("Nil point receivcer")
	}
	err = rd.conn.Close()
	if err!=nil{
		log4go.Error(err)
	}
	return
}

func (rd *RedisHelper)RedisHelperInit(protocol, ip, port string){
	if rd==nil{
		log4go.Error("RedisHelper instance is nil")
		return
	}
	rd.protocol = protocol
	rd.ip = ip
	rd.port = port
}

/*func newRedisHelper() *RedisHelper{
	return &RedisHelper{
		protocol: "tcp",
		ip: "localhost",
		port: "6379",
	}
}*/

func (rd *RedisHelper)GetProtocol() string {
	return rd.protocol
}

func (rd *RedisHelper)GetIp() string {
	return rd.ip
}

func (rd *RedisHelper)GetPort() string {
	return rd.port
}

func (rd *RedisHelper)SetProtocol(protocol string){
	rd.protocol = protocol
}

func (rd *RedisHelper)SetIp(ip string){
	rd.ip = ip
}

func (rd *RedisHelper)SetPort(port string){
	rd.port = port
}