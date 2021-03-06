/**
 * @Author: mjzheng
 * @Description:
 * @File:  simple_client
 * @Version: 1.0.0
 * @Date: 2020/10/19 上午10:25
 */

package main

import (
	"bytes"
	"context"
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

const (
	CLIENT_STATUS_JOIN      = 1
	CLIENT_STATUS_UPLOAD    = 2
	CLIENT_STATUS_SEND_DATA = 3
)

func ComposeMsg(msg proto.Message) (data []byte) {
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

	//fmt.Println(len(pData))

	data = buf.Bytes()
	return
}

func JoinRoom(conn net.Conn, uid int64, roomId int64) {
	headReq := &simple_msg.HeadReq{
		Cmd:    int32(simple_av.BIG_CMD_SIMPLE_AV),
		Subcmd: int32(simple_av.SUB_CMD_JoinRoom),
		Seq:    1,
	}
	bodyReq := &simple_av.JoinRoomReq{
		RoomId: roomId,
		Uid:    uid,
	}
	var err error
	headReq.Ex, err = proto.Marshal(bodyReq)
	if err != nil {
		return
	}
	fmt.Println(headReq, bodyReq)
	pData := ComposeMsg(headReq)
	conn.Write(pData)
}

func OnReceive(conn net.Conn) {
	ctx := context.Background()
	buf := make([]byte, 1024)
	fmt.Println("len", len(buf))
	from := 0
	for {
		total, err := conn.Read(buf[from:])
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return //终止程序
		}

		buf, from, _ = ParseMsg(ctx, buf, total+from)
	}
}

func Upload(conn net.Conn, uid int64, roomId int64) {
	msg := &simple_msg.HeadReq{
		Cmd:    int32(simple_av.BIG_CMD_SIMPLE_AV),
		Subcmd: int32(simple_av.SUB_CMD_Upload),
		Seq:    int32(1),
	}
	req := &simple_av.UploadReq{
		Uid:    uid,
		RoomId: roomId,
	}
	var err error
	msg.Ex, err = proto.Marshal(req)
	if err != nil {
		return
	}

	pData := ComposeMsg(msg)
	conn.Write(pData)
}

func SendData(conn net.Conn, uid int64, roomId int64, seq int32, payload []byte) {
	msg := &simple_msg.HeadReq{
		Cmd:    int32(simple_av.BIG_CMD_SIMPLE_AV),
		Subcmd: int32(simple_av.SUB_CMD_SendData),
		Seq:    seq,
	}
	req := &simple_av.SendDataReq{
		Uid:    uid,
		RoomId: roomId,
		Data:   []byte(payload),
	}
	var err error
	msg.Ex, err = proto.Marshal(req)
	if err != nil {
		return
	}

	pData := ComposeMsg(msg)
	conn.Write(pData)
}

func IntToBytes(n int32) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func StartSimpleClient(uid int64, roomId int64, upload int64) {
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

	go OnReceive(conn)

	JoinRoom(conn, uid, roomId)

	i := int32(2)
	for {
		time.Sleep(time.Millisecond * 1000)
		if upload > 0 {
			SendData(conn, uid, roomId, i, IntToBytes(i))
		}
		i++
	}
}