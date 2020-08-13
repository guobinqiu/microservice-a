package main

import (
	"flag"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"microservice-a/endpoints"
	"microservice-a/middlewares"
	"microservice-a/registry"
	"microservice-a/services"
	"microservice-a/transports"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	name string
	port int
)

func init() {
	flag.StringVar(&name, "name", "serviceA", "服务名")
	flag.IntVar(&port, "port", 5678, "服务端口")
	flag.Parse()
}

func main() {
	svc := services.UserService{}
	getUserNameEndpoint := endpoints.MakeGetUserNameEndpoint(svc)
	delUserEndpoint := endpoints.MakeDelUserEndpoint(svc)

	//加入限流中间件
	limiter := rate.NewLimiter(1, 5)
	rateLimitMiddleware := middlewares.RateLimitMiddleware(limiter)
	getUserNameEndpoint = rateLimitMiddleware(getUserNameEndpoint)
	delUserEndpoint = rateLimitMiddleware(delUserEndpoint)

	//统一错误处理
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(transports.ErrorEncoder),
	}

	//1 解码请求发送至endpoint
	//2 编码响应发送至客户端
	getUserNameHandler := httptransport.NewServer(
		getUserNameEndpoint,
		transports.DecodeUserRequest,
		transports.EncodeGetUserNameResponse,
		options...,
	)
	delUserHandler := httptransport.NewServer(
		delUserEndpoint,
		transports.DecodeUserRequest,
		transports.EncodeDelUserResponse,
		options...,
	)

	//独立路由插件。用gin，iris等框架的话就不需要引入了
	r := mux.NewRouter()
	//r.Handle("/users/{uid}", handler)
	r.Methods("GET").Path("/users/{uid}").Handler(getUserNameHandler)
	r.Methods("DELETE").Path("/users/{uid}").Handler(delUserHandler)

	//健康检查
	r.Methods("GET").Path("/ping").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM) //等待关闭信号量来解除主线程的阻塞状态

	reg, err := registry.NewConsulRegistry(name, port)
	if err != nil {
		panic(err)
	}

	//服务注册
	go func() {
		reg.RegService()
		http.ListenAndServe(":"+strconv.Itoa(port), r);
	}()

	<-s //阻塞主线程

	//反注册，删除服务
	reg.Deregister()
}
