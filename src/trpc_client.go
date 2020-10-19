/**
 * @Author: mjzheng
 * @Description:
 * @File:  trpc_client.go
 * @Version: 1.0.0
 * @Date: 2020/10/19 上午10:31
 */

package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mjproto/simple_av"
	"github.com/mjproto/simple_msg"
	"google.golang.org/grpc"
	"log"
	"time"
)

func JoinRoom2(stream simple_msg.SimpleMsg_HeadClient, uid int64, roomId int64) {
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

	stream.SendMsg(headReq)

	reply, err := stream.Recv()
	if err != nil {
		log.Printf("failed to recv: %v", err)
	}
	fmt.Println("join room", headReq, bodyReq, reply)

	HandleTrpcMsg(reply)
}

func SendData2(stream simple_msg.SimpleMsg_HeadClient, uid int64, roomId int64, seq int32, payload []byte) {
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

	stream.SendMsg(msg)

	reply, err := stream.Recv()
	if err != nil {
		log.Printf("failed to recv: %v", err)
	}

	fmt.Println("send data", msg, req, reply)
}

func StartTRPCClient(uid int64, roomId int64, upload int64) {

	address := "localhost:50000"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("faild to connect: %v", err)
	}
	defer conn.Close()

	c := simple_msg.NewSimpleMsgClient(conn)

	stream, err := c.Head(context.Background())
	if err != nil {
		log.Printf("failed to call: %v", err)
		return
	}

	JoinRoom2(stream, uid, roomId)

	i := int32(2)
	for {
		time.Sleep(time.Millisecond * 1000)
		if upload > 0 {
			SendData2(stream, uid, roomId, i, IntToBytes(i))
		} else {
			req := &simple_msg.HeadReq{}
			err := stream.RecvMsg(req)
			if err != nil {
				log.Println("failed to recv request", err)
			}
			HandleTrpcRequest(req)
			//reply, err := stream.Recv()
			//if err != nil {
			//	log.Printf("failed to recv: %v", err)
			//}
			//HandleTrpcMsg(reply)
		}
		i++
	}
}

func HandleTrpcMsg(rsp *simple_msg.HeadRsp) {
	switch rsp.Subcmd {
	case int32(simple_av.SUB_CMD_JoinRoom):
		bodyRsp := &simple_av.JoinRoomRsp{}
		err := proto.Unmarshal(rsp.GetEx(), bodyRsp)
		if err != nil {
		}
		fmt.Println("recv join rsp", rsp, bodyRsp)
	case int32(simple_av.SUB_CMD_SendData):
		bodyRsp := &simple_av.SendDataReq{}
		err := proto.Unmarshal(rsp.GetEx(), bodyRsp)
		if err != nil {

		}
		fmt.Println("recv data ", rsp, bodyRsp)

	case int32(simple_av.SUB_CMD_Upload):
	case int32(simple_av.SUB_CMD_ExitRoom):
	default:
	}
}

func HandleTrpcRequest(req *simple_msg.HeadReq) {
	switch req.Subcmd {
	case int32(simple_av.SUB_CMD_SendData):
		bodyReq := &simple_av.SendDataReq{}
		err := proto.Unmarshal(req.GetEx(), bodyReq)
		if err != nil {

		}
		fmt.Println("recv data ", req, bodyReq)
	default:
	}
}
