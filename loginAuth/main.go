// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/jeanphorn/log4go"
	"io"
	"loginAuth/config"
	"loginAuth/grpc"
	"loginAuth/ipaddress"
	"loginAuth/messagebroke"
	"loginAuth/redisdao"
	"loginAuth/serviceregister"
	"loginAuth/userauth"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

var (
	conf *config.Config
	nacosService *serviceregister.NacosService
	messageBroke *messagebroke.MessageBroke
	redisRPCServer *grpc.RedisRPCServer
	redisHelper *redisdao.RedisHelper
	message string
	ipAddress *ipaddress.IPAddress
	err error
	messageId uint64                   //messageId increases every time a message was received or was send
	groupId string                     //groupId should be unique for each service
	dec *json.Decoder
	v map[string]interface{}           //v is used for saving results from message transfer into json object
	userAuth *userauth.UserAuth        //userAuth is used for authentication
)

func init(){
	//variable declaration
	var success bool
	var bytes []byte

	//Read config from config.json
	conf, err = config.ConfigInstance()
	if err!=nil{
		log4go.Error("Fail to initialize from config.json")
		return
	}

	//get ipaddress and set it to DataPush and UpdateToken config's ip
	ipAddress = &ipaddress.IPAddress{}
	ipAddress.GetIpAddress()
	conf.DataPush.IpAddr=ipAddress.IpAddress
	conf.UpdateToken.IpAddr=ipAddress.IpAddress
	groupId="group-"+ipAddress.IpAddress

	//service register
	nacosService = &serviceregister.NacosService{}
	nacosService.NacosServiceInit()
	nacosService.SetServerConfig(conf.NacosServerConfig.IpAddr,conf.NacosServerConfig.ContextPath,conf.NacosServerConfig.Port)
	err = nacosService.CreateClient()
	if err!=nil {
		log4go.Error("Fail to create NacosService client")
		return
	}
	success, err = nacosService.RegisterService(conf.NacosServerConfig.ServicePort)
	if !success||err!=nil{
		log4go.Error("Fail to register service")
		return
	}

	//start redisgrpcserver
	redisRPCServer = &grpc.RedisRPCServer{}
	redisRPCServer.RedisRPCServerInit(conf.RedisRPCServerConfig.ServerAddr,conf.RedisRPCServerConfig.NetProtocol,
		conf.RedisRPCServerConfig.CertFile,conf.RedisRPCServerConfig.KeyFile)
	go redisRPCServer.StartRPCServer()

	//make redishelper object
	redisHelper = &redisdao.RedisHelper{}
	redisHelper.RedisHelperInit(conf.RedisHelperConfig.Protocol, conf.RedisHelperConfig.Ip, conf.RedisHelperConfig.Port)
	err = redisHelper.ConnectRedis()
	if err!=nil {
		log4go.Error("Fail to connect to redis")
		return
	}

	//send message to kafka to get redis data from data push and consistency layer
	conf.DataPush.MessageId=1
	bytes, err = json.Marshal(conf.DataPush)
	if err!=nil {
		log4go.Error("fail to read datapush message format")
		return
	}
	message = string(bytes)
	messageBroke = &messagebroke.MessageBroke{}
	messageBroke.MessageBrokeInit(conf.MessageBrokeConfig.Address, conf.MessageBrokeConfig.BufferSize,conf.MessageBrokeConfig.TopicName)
	go messageBroke.MessageProducer("dataPush")
	go messageBroke.MessageConsumer([]string{"confirm"}, groupId)
	messageBroke.SendMessage(message)

	//before data push and consistency service had sent data into redis, the init process has to be blocked
	//the question is how to make the gourpid of a consumer unique
	for{
		message=<- messageBroke.Value
		dec = json.NewDecoder(strings.NewReader(message))
		for dec.More(){
			if err = dec.Decode(&v); err != nil{
				log4go.Error("fail to decode from confirm message from data push and consistency layer")
			}
		}
		if v["result"]=="success"&&v["ip"]==ipAddress.IpAddress&&uint64(reflect.ValueOf(v["messageId"]).Float())>messageId{
			messageBroke.StopConsumer()
			log4go.Info("data push and consistency layer succeed in pushing data into redis")
			break
		}else{
			log4go.Error("data push and consistency layer fail to push data into redis")
		}
	}
	//set new messageId
	messageId=uint64(reflect.ValueOf(v["messageId"]).Float())+1

	//It's not until now, the data has been successfully pushed into redis
	//And it's time for server to deal requests
	//But there's still some tasks which should be running in the background
	userAuth = &userauth.UserAuth{}
	userAuth.UserAuthInit(redisHelper,messageBroke,conf,&messageId)
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	if r.Method=="PUT"{
		//extract username and password from request and check if there is really a username and a password
		s, _ := url.ParseQuery(r.URL.RawQuery)
		username := s.Get("username")
		if username==""{
			io.WriteString(w, "requires username")
			return
		}
		password := s.Get("password")
		if password==""{
			io.WriteString(w, "requires password")
			return
		}
		result, err := userAuth.CheckPassword(username, password)
		if err!=nil{
			log4go.Info(fmt.Sprintf("Wrong password for user %v"), username)
			io.WriteString(w, result)
			return
		}
		result, err = userAuth.CheckToken(username)
		if err!=nil{
			log4go.Error("fail to generate token")
			io.WriteString(w, result)
			return
		}
		io.WriteString(w, result)
		return
	}
}

func main(){
	http.HandleFunc("/login", loginHandler)
	err := http.ListenAndServe(string(conf.NacosServerConfig.ServicePort), nil)
	if err!=nil {
		log4go.Error(fmt.Sprintf("ListenAndServe: %v", err))
	}
}