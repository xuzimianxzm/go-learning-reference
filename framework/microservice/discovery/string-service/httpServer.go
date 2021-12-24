package main

import (
	"context"
	"discovery/common/config"
	"discovery/string-service/endpoint"
	"discovery/string-service/plugins"
	"discovery/string-service/service"
	"discovery/string-service/transport"
	"net/http"
	"strconv"
)

func onlyStartASingleHttpService(servicePort *int) {
	var svc service.Service
	svc = service.StringService{}
	// add logging middleware
	svc = plugins.LoggingMiddleware(config.KitLogger)(svc)

	stringEndpoint := endpoint.MakeStringEndpoint(svc)

	//创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	//把算术运算Endpoint和健康检查Endpoint封装至StringEndpoints
	endpts := endpoint.StringEndpoints{
		StringEndpoint:      stringEndpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	//创建http.Handler
	r := transport.MakeHttpHandler(context.Background(), endpts, config.KitLogger)

	//http server
	config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
	handler := r
	errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
}
