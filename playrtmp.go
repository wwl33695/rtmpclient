package main

import "net"
import "rtmp"

func MessageProc(client *rtmp.RTMPClient) {
    for {
        err := client.RecvMessage()
        if err != nil {
            println("HandleMessage error", err.Error())
            return
        }
    }
}

func main() {
/*
	ns, err := net.LookupHost("live.hkstv.hk.lxdns.com")
	if err != nil {
		println("Err:", err.Error())
		return
	}

	for _, n := range ns {
		println("--", n) 
	}
*/

//    conn, err := net.Dial("tcp", ":1935")
    conn, err := net.Dial("tcp", "116.242.0.29:1935")
    if err != nil {
        println("连接服务端失败:", err.Error())
        return
    }
	defer func() {
		conn.Close()
	}()

    client := rtmp.GetRTMPClient(conn)

    err = client.Handshake()
    if err != nil {
        println("Handshake error", err.Error())
        return
    }

//    err = client.Connect("rtmp://127.0.0.1/live/test")
    err = client.Connect("rtmp://live.hkstv.hk.lxdns.com/live")
    if err != nil {
        println("Connect error", err.Error())
        return
    }

    MessageProc(client)

	println("playrtmp exit")
}
