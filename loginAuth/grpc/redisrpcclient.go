// Copyright 2020 yehaotian. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//Package grpc use grpc to receive data from another program
package grpc

import (
	"fmt"
	"time"
	"context"
	"loginAuth/data"
	"loginAuth/user"

	"github.com/jeanphorn/log4go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "loginAuth/grpc/redisrpc"
)

//RedisRPCClient act as a rpc client and it uses https encrypted by tls to transport messages
type RedisRPCClient struct{
	caFile string
	serverNameOverride string
	serverAddr string
	client pb.IRedisRPCClient
	conn *grpc.ClientConn
}

func (redisRPCClient *RedisRPCClient)RedisRPCClientInit(){
	redisRPCClient.caFile = data.Path("x509/ca_cert.pem")
	redisRPCClient.serverNameOverride = "x.test.example.com"
	redisRPCClient.serverAddr = "localhost:2100"
}

func (redisRPCClient *RedisRPCClient)GetCaFile() string {
	return redisRPCClient.caFile
}

func (redisRPCClient *RedisRPCClient)GetServerNameOverride() string {
	return redisRPCClient.serverNameOverride
}

func (redisRPCClient *RedisRPCClient)GetServerAddr() string {
	return redisRPCClient.serverAddr
}

func (redisRPCClient *RedisRPCClient)SetCaFile(caFile string) {
	redisRPCClient.caFile = caFile
}

func (redisRPCClient *RedisRPCClient)SetServerNameOverride(serverNameOverride string) {
	redisRPCClient.serverNameOverride = serverNameOverride
}

func (redisRPCClient *RedisRPCClient)SetServerAddr(serverAddr string) {
	redisRPCClient.serverAddr = serverAddr
}

//UpdateToken sends a user to rpc server to update it's token in redis
func (redisRPCClient *RedisRPCClient)UpdateToken(user user.User){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := redisRPCClient.client.UpdateToken(ctx,&pb.User{
		UserName: user.GetUserName(),
		Password: user.GetPassword(),
		Token: user.GetToken(),
	})
	if err!=nil {
		log4go.Error(fmt.Sprintf("Fail to update token, err: %v", err))
	}
}

//InsertUsers sends an array of users to the server
func (redisRPCClient *RedisRPCClient)InsertUsers(user []user.User){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := redisRPCClient.client.InsertUsers(ctx)
	if err != nil{
		log4go.Error(fmt.Sprintf("Fail to get stream, err: %v", err))
		return
	}
	for _, _user := range user{
		if err := stream.Send(&pb.User{
			UserName: _user.GetUserName(),
			Password: _user.GetPassword(),
			Token: _user.GetToken(),
		}); err != nil{
			log4go.Error(fmt.Sprintf("Fail to send user, err: %v", err))
			return
		}
	}
	_, err = stream.CloseAndRecv()
	if err != nil{
		log4go.Error(fmt.Sprintf("Send users failed, err: %v", err))
		return
	}
	log4go.Info(fmt.Sprintf("Succeeded in sending users"))
}

//StartRPCClient creates a rpc client and use it to send messages and when it finishes, close the connection
func (redisRPCClient *RedisRPCClient)StartRPCClient() error {
	var opts []grpc.DialOption
	creds, err := credentials.NewClientTLSFromFile(redisRPCClient.caFile,redisRPCClient.serverNameOverride)
	if err != nil {
		log4go.Error(fmt.Sprintf("Failed to create TLS credentials %v", err))
		return err
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	redisRPCClient.conn, err = grpc.Dial(redisRPCClient.serverAddr,opts...)
	if err!=nil{
		log4go.Error(fmt.Sprintf("fail to dail: %v", err))
		return err
	}
	redisRPCClient.client = pb.NewIRedisRPCClient(redisRPCClient.conn)
	return nil
}

//CloseRPCConnection close the grpc.Conn
func (redisRPCClient *RedisRPCClient)CloseRPCConnection(){
	redisRPCClient.conn.Close()
}