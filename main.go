/**
 * @Author: mjzheng
 * @Description:
 * @File:  main.go
 * @Version: 1.0.0
 * @Date: 2020/6/23 上午11:49
 */

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mjproto/simple_av"
	"github.com/mjproto/simple_msg"
	"net"
	"time"
)

const (
	SX = 0x2
	EX = 0X3
)

func AddSplit(msg *simple_msg.HeadReq) (data []byte) {
	pData, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	buf.WriteByte(0x2)
	lenBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBuf, uint16(len(pData)))
	buf.Write(lenBuf)
	buf.Write(pData)
	buf.WriteByte(0x3)

	fmt.Println(len(pData))

	data = buf.Bytes()
	return
}

func main() {
	strIP := "localhost:50000"
	var conn net.Conn
	var err error

	//连接服务器
	for conn, err = net.Dial("tcp", strIP); err != nil; conn, err = net.Dial("tcp", strIP) {
		fmt.Println("connect", strIP, "fail")
		time.Sleep(time.Second)
		fmt.Println("reconnect...")
	}
	fmt.Println("connect", strIP, "success")
	defer conn.Close()

	msg := &simple_msg.HeadReq{
		Cmd:    0x100,
		Subcmd: int32(simple_av.SUB_CMD_JoinRoom),
		Seq:    0,
	}
	req := &simple_av.JoinRoomReq{
		RoomId: 1000,
		Uid:    88881811,
	}

	msg.Ex, err = proto.Marshal(req)

	//protobuf编码
	pData := AddSplit(msg)
	conn.Write(pData)

	for i := 0; i < 100; i++ {
		msg = &simple_msg.HeadReq{
			Cmd:    0x100,
			Subcmd: int32(simple_av.SUB_CMD_Upload),
			Seq:    int32(i),
		}
		req := &simple_av.UploadReq{
			RoomId: 1000,
		}
		msg.Ex, err = proto.Marshal(req)

		pData := AddSplit(msg)
		conn.Write(pData[0 : len(pData)/2])
		time.Sleep(time.Millisecond * 10)
		conn.Write(pData[len(pData)/2:])
	}
}
