package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

/*
编写一个 HTTP 服务器

1. 接收客户端 request，并将 request 中带的 header 写入 response header
2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
4. 当访问 localhost/healthz 时，应返回 200

*/
func main() {
	log.Println("server starting...")

	// 1-3
	http.HandleFunc("/test", handleTest)

	//4. 当访问 localhost/healthz 时，应返回 200
	http.HandleFunc("/healthz", healthz)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalln("ListenAndServe: " + err.Error())
		return
	}
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	//1. 接收客户端 request，并将 request 中带的 header 写入 response header
	//r.Header.Write(w)
	// for...range 写法
	for k, v := range r.Header {
		// v 为数组
		for _, vv := range v {
			log.Println("add header", k, vv)
			w.Header().Add(k, vv)
		}
	}

	//2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	version, exists := os.LookupEnv("VERSION")
	if exists {
		log.Println("found env VERSION", version)
		w.Header().Add("VERSION", version)
		w.Header().Write(w)
	}

	//3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	ip, err := getIP(r)
	if err != nil {
		log.Println("parse IP error,", err.Error())
	}
	log.Println("client IP:", ip, ", HTTP status: ", w.Header().Get("status"))
}

// 根据头信息返回IP
func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}

func healthz(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "200")
}
