package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	fmt.Println("命令行参数数量:", len(os.Args))
	for k, v := range os.Args {
		fmt.Printf("args[%v]=[%v]\n", k, v)
	}

	var mode string

	var s_innerListenAddress string
	var s_outerListenAddress string
	var s_signalServerAddress string

	var c_signalServerAddress string
	var c_serverAddress string
	var c_socksServerAddress string

	var key string

	flag.StringVar(&s_innerListenAddress, "i", "0.0.0.0:1080", "内网监听地址,默认0.0.0.0:1080")
	flag.StringVar(&s_outerListenAddress, "o", "0.0.0.0:1086", "外网监听地址,默认0.0.0.0:1086")
	flag.StringVar(&s_signalServerAddress, "ssi", "0.0.0.0:1084", "信令服务器地址,默认0.0.0.0:1084")

	flag.StringVar(&c_signalServerAddress, "csi", "s2.lrsj.fun:1084", "信令服务器地址,默认s2.lrsj.fun:1082")
	flag.StringVar(&c_serverAddress, "s", "s2.lrsj.fun:1086", "服务监听地址,默认s2.lrsj.fun:1086")
	flag.StringVar(&c_socksServerAddress, "socks", "s2.lrsj.fun:10808", "socks5服务地址,默认127.0.0.1:1080")
	flag.StringVar(&mode, "m", "client", "模式,默认server,可选client")
	flag.StringVar(&key, "k", "WyuIRGXac98iFfbbOx30x3RLfiNLm8zN", "加密隧道key(32字节)")

	flag.Parse()

	switch mode {
	case "server":
		{
			signalServer := NewSignalServer(s_signalServerAddress)
			signalServer.Start()
			server := NewServer(s_innerListenAddress, s_outerListenAddress, signalServer, []byte(key))
			server.Start()
		}
	case "client":
		{
			client := NewClient(c_signalServerAddress, c_serverAddress, c_socksServerAddress, []byte(key))
			client.Start()
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	signal := <-signals
	log.Printf("receive signal %s, graceful ending...\n", signal)

}
