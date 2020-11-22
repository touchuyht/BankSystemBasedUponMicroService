// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Package grpc use grpc to receive data from another program
package grpc

import (
	"context"
	"fmt"
	"github.com/jeanphorn/log4go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"loginAuth/config"
	"loginAuth/data"
	pb "loginAuth/grpc/redisrpc"
	"loginAuth/redisdao"
	"loginAuth/user"
	"net"
)

var redisHelper *redisdao.RedisHelper

//RedisRPCServer act as a rpc server and have some attributes which describes
//the server's information. Use it in the following order:
//RedisRPCServerInit->go startRPCServer()
type RedisRPCServer struct{
	pb.UnimplementedIRedisRPCServer
	netProtocol string
	serverAddr string
	certFile string
	keyFile string
}

func (redisRPCServer *RedisRPCServer)SetNetProtocol(netProtocol string){
	redisRPCServer.netProtocol = netProtocol
}

func (redisRPCServer *RedisRPCServer)SetServerAddr(serverAddr string){
	redisRPCServer.serverAddr = serverAddr
}

func (redisRPCServer *RedisRPCServer)SetCertFile(certFile string){
	redisRPCServer.certFile = certFile
}

func (redisRPCServer *RedisRPCServer)SetKeyFile(keyFile string){
	redisRPCServer.keyFile = keyFile
}

func (redisRPCServer *RedisRPCServer)GetNetProtocol() string {
	return redisRPCServer.netProtocol
}

func (redisRPCServer *RedisRPCServer)GetServerAddr() string {
	return redisRPCServer.serverAddr
}

func (redisRPCServer *RedisRPCServer)GetCertFile() string {
	return redisRPCServer.certFile
}

func (redisRPCServer *RedisRPCServer)GetKeyFile() string {
	return redisRPCServer.keyFile
}

func (redisRPCServer *RedisRPCServer)RedisRPCServerInit(serverAddr, netProtocol, certFile, keyFile string) error {
	redisRPCServer.serverAddr = serverAddr
	redisRPCServer.netProtocol = netProtocol
	redisRPCServer.keyFile = data.Path(keyFile)
	redisRPCServer.certFile = data.Path(certFile)
	c, err := config.ConfigInstance()
	if err!=nil {
		log4go.Error("Fail to init config instance")
		return err
	}
	redisHelper = &redisdao.RedisHelper{}
	redisHelper.RedisHelperInit(c.RedisHelperConfig.Protocol, c.RedisHelperConfig.Ip, c.RedisHelperConfig.Port)
	err = redisHelper.ConnectRedis()
	if err != nil{
		log4go.Error(fmt.Sprintf("Fail to connect redis, err: %v", err))
		return err
	}
	return err
}

//UpdateToken updates the token of a user
func (redisRPCServer *RedisRPCServer)UpdateToken(ctx context.Context, pbUser *pb.User) (*pb.Result, error){
	_user := user.User{}
	_user.SetUserName(pbUser.UserName)
	_user.SetPassword(pbUser.Password)
	_user.SetToken(pbUser.Token)
	err := redisHelper.UpdateToken(_user)
	if err != nil{
		log4go.Error(fmt.Sprintf("Fail to update token, error: %v", err))
		return &pb.Result{Result: "failure"}, err
	}
	return &pb.Result{Result: "success"}, nil
}


//InsertUsers inserts an array of users into redis
func (redisRPCServer *RedisRPCServer)InsertUsers(stream pb.IRedisRPC_InsertUsersServer) error {
	for {
		pbUser, err := stream.Recv()
		if err == io.EOF{
			return stream.SendAndClose(&pb.Result{Result: "success"})
		}
		if err!=nil{
			log4go.Error(fmt.Sprintf("Fail to send in a stream, err: %v", err))
			return err
		}
		redisHelper.InsertUser(user.NewUser(pbUser.UserName,pbUser.Password,pbUser.Token))
	}
}

//StartRPCServer creates a rpc server and makes it serve on a specified port
func (redisRPCServer *RedisRPCServer)StartRPCServer() error {
	lis, err := net.Listen(redisRPCServer.netProtocol, redisRPCServer.serverAddr)
	if err!=nil{
		log4go.Error(fmt.Sprintf("Failed to listen, err: %v", err))
		return err
	}
	var opts []grpc.ServerOption
	creds, err := credentials.NewServerTLSFromFile(redisRPCServer.certFile, redisRPCServer.keyFile)
	if err != nil {
		log4go.Error(fmt.Sprintf("Failed to create TLS " +
			"credentials, err %v", err))
		return err
	}
	opts = []grpc.ServerOption{grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterIRedisRPCServer(grpcServer, redisRPCServer)
	log4go.Info("rpc server start to listen")
	grpcServer.Serve(lis)
	defer redisHelper.Close()
	return nil
}