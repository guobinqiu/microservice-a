package registry

import (
	"fmt"
	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
)

type ConsulRegisterCenter struct {
	name      string
	port      int
	apiclient *consulapi.Client
	serviceid string
}

//临时用用
var Port int

func NewConsulRegistry(name string, port int) (*ConsulRegisterCenter, error) {
	cli, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, err
	}
	Port = port
	return &ConsulRegisterCenter{name, port, cli, name + uuid.New().String()}, nil
}

func (c *ConsulRegisterCenter) RegService() error {
	//consul运行在容器内，在容器内访问宿主机ip需要使用host.docker.internal地址（macos）
	check := consulapi.AgentServiceCheck{
		HTTP:     fmt.Sprintf("http://host.docker.internal:%d/ping", c.port),
		Interval: "5s",
	}

	reg := consulapi.AgentServiceRegistration{
		ID:      c.serviceid, //唯一
		Name:    c.name,
		Address: "127.0.0.1", //供客户端调用的微服务地址
		Port:    c.port,
		Tags:    []string{"primary"},
		Check:   &check,
	}

	return c.apiclient.Agent().ServiceRegister(&reg)
}

func (c *ConsulRegisterCenter) Deregister() error {
	return c.apiclient.Agent().ServiceDeregister(c.serviceid)
}
