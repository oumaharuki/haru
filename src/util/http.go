package util

import (
	"io"
	"net"
	"net/http"
	"time"
)

var httpClient http.Client

func init() {
	httpClient = http.Client{

		Transport: &http.Transport{

			Dial: func(netw, addr string) (net.Conn, error) {

				c, err := net.DialTimeout(netw, addr, time.Second*1) //设置建立连接超时

				if err != nil {

					return nil, err

				}

				c.SetDeadline(time.Now().Add(1 * time.Second)) //设置发送接收数据超时

				return c, nil

			},
		},
	}
}
func HttpGet(url string) (rs string, err error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf := make([]byte, 4096)
	for {
		n, err1 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err1 != nil && err1 != io.EOF {
			err = err1
			return
		}

		rs += string(buf[:n])
	}
	return
}
