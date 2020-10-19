/**
 * @Author: mjzheng
 * @Description:
 * @File:  parse.go
 * @Version: 1.0.0
 * @Date: 2020/6/29 下午7:02
 */

package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mjproto/simple_msg"
)

const (
	STATUS_START_EX = 1
	STATUS_LENGTH   = 2
	STATUS_BODY     = 3
	STATUS_END_EX   = 4
	STATUS_COMPLETE = 5
)

func ParseMsg(ctx context.Context, buf []byte, total int) (remain []byte, remainLen int, msg []byte) {
	useLen := 0
	from := 0
	status := STATUS_START_EX
	needLen := 1
	msg = nil
	for from+needLen <= total {
		switch status {
		case STATUS_START_EX:
			if buf[from] != 0x2 {
				fmt.Println("unexcept start error")
				break
			}
			from += needLen
			needLen = 2
			status = STATUS_LENGTH
		case STATUS_LENGTH:
			msgLen := int(binary.BigEndian.Uint16(buf[from : from+needLen]))
			from += needLen
			needLen = msgLen
			status = STATUS_BODY
		case STATUS_BODY:
			msg = buf[from : from+needLen]
			from += needLen
			needLen = 1
			status = STATUS_END_EX
		case STATUS_END_EX:
			if buf[from] != 0x3 {
				fmt.Println("unexcept end error")
				break
			}
			from += needLen
			HandleMsg(ctx, msg)
			useLen = from

			needLen = 1
			status = STATUS_START_EX
		}
	}
	if useLen < total {
		// move
		remainLen = total - useLen
		for i := 0; i < remainLen; i++ {
			buf[i] = buf[useLen+i]
		}
		//fmt.Println("reamin len", total, useLen, remainLen)
		return buf, remainLen, msg
	} else {
		return buf, 0, msg
	}
}

func HandleMsg(ctx context.Context, pData []byte) {
	headRsp := &simple_msg.HeadRsp{}
	err := proto.Unmarshal(pData, headRsp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(headRsp)
}
