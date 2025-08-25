package main

import (
	"fmt"

	"github.com/nhh-7/go-chat/internal/config"
	httpsserver "github.com/nhh-7/go-chat/internal/https_server"
	"github.com/nhh-7/go-chat/utils/zlog"
)

func main() {
	conf := config.GetConfig()
	host := conf.MainConfig.Host
	port := conf.MainConfig.Port

	if err := httpsserver.GE.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		zlog.Fatal("server run err: ")
		return
	}
}
