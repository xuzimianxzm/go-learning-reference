package main

import (
	"discovery/string-service/config"
	"flag"
	"fmt"
	"github.com/longjoy/micro-go-book/common/discover"
	uuid "github.com/satori/go.uuid"
	"os"
	"os/signal"
	"syscall"
)

var (
	servicePort = flag.Int("service.port", 10085, "service port")
	serviceHost = flag.String("service.host", "127.0.0.1", "service host")
	consulPort  = flag.Int("consul.port", 8500, "consul port")
	consulHost  = flag.String("consul.host", "127.0.0.1", "consul host")
	serviceName = flag.String("service.name", "string", "service name")
	errChan     = make(chan error)
)

func main() {
	flag.Parse()
	go onlyStartASingleHttpService(servicePort)
	//go  startHttpServerAndRegisterServiceToConsul()

	go listenerExitCommand()

	exit(nil, "")
}

func startHttpServerAndRegisterServiceToConsul() {
	var discoveryClient discover.DiscoveryClient
	discoveryClient, err := discover.NewKitDiscoverClient(*consulHost, *consulPort)

	if err != nil {
		config.Logger.Println("Get Consul Client failed")
		os.Exit(-1)

	}

	instanceId := *serviceName + "-" + uuid.NewV4().String()

	//http server
	go func() {
		//启动前执行注册
		if !discoveryClient.Register(*serviceName, instanceId, "/health", *serviceHost, *servicePort, nil, config.Logger) {
			config.Logger.Printf("string-service for service %s failed.", serviceName)
			// 注册失败，服务启动失败
			os.Exit(-1)
		}
	}()

	exit(discoveryClient, instanceId)
}

func exit(discoveryClient discover.DiscoveryClient, instanceId string) {
	error := <-errChan
	//服务退出取消注册
	if discoveryClient != nil {
		discoveryClient.DeRegister(instanceId, config.Logger)
		config.Logger.Println(error)
	}
	os.Exit(-1)
}

func listenerExitCommand() {
	func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
}
