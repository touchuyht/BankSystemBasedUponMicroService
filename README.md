## 系统设计

### 系统概述

​	本系统主要是针对银行的基本业务进行一个实现，包括存款、取款、转账以及贷款。传统的银行业务都是以单体核心业务的方式进行的，银行的基本业务以单体应用的方式运行在高性能的小型机上，能够保证系统的强一致性，本系统主要是对银行的核心业务采用分布式系统进行的一次试验。

###  系统架构与交互流程

#### 系统架构

​	系统基于微服务理念和业务的边界划分除了多个子系统，划分出的子系统包括统一认证中心、用户信息中心以及交易中心。

​	相关子系统的介绍：

​	1）统一认证中心主要完成用户登录时密码校对、生成token和分发token。

​	2）用户信息中心提供用户基本信息的查询和修改。

​	3）交易中心提供存款、取款、转账和贷款的业务，其中贷款会根据用户的信息决定是否给用户贷款以及贷款的金额。

​	4）数据推送与持久化层主要负责向Service推送数据、分发token、想HBase更新修改的信息等。

#### 系统间的交互流程

​	用户登录，token需要在所有微服务进行同步，容器启动之后先向nacos注册自己，并向消息队列发送一条消息，数据推送与持久化层在监听到新的容器启动之后会去拉去它的信息，并向其redis数据库推送相应的数据，并在推送完成之后向消息队列插入消息。在统一认正中心的数据发生变化之后，其需要消息队列发送一条消息，数据推送与持久化层在接收到这条消息以后需要向其本地redis中的数据进行修改并写入hbase，并通过grpc的方式同步的修改其他redis数据库的数据。

​	当用户访问用户信息中心时，用户信息中心根据用户请求头中携带的token进行校验，在校验通过之后提供对应的服务，并在校验失败是返回出错提示。用户信息中心主要就是用户个人信息的查看和修改，并在用户信息发生修改时向消息队列发送消息，数据推送与持久化层在接收到消息之后向将修改写入自己的数据库并同步到hbase，并通过grpc向指定的数据库进行数据的同步。

​	当用户访问交易中心时，交易中心根据用户请求头中携带的token进行校验，在校验通过之后提供对应的服务，比如余额查询、行内转账，在进行行内转账时，采用的是异步的方式，先本地扣款，然后插入一条转账消息到消息队列，消费者被接受到之后进行转账，如果成功插入一条消息，我在接收到这条成功转账的消息之后完成交易，否则回滚。
