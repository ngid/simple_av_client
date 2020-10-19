/**
 * @Author: mjzheng
 * @Description:
 * @File:  main.go
 * @Version: 1.0.0
 * @Date: 2020/6/23 上午11:49
 */

package main

import (
	"flag"
	"fmt"
)

// audience: ./simple_av_client -u 252238532 -r 1000

// anchor:  ./simple_av_client -u 88881811 -r 1000 -a 1

var pRoomId = flag.Int64("r", 0, "room id")
var pUid = flag.Int64("u", 0, "uid")
var pUpload = flag.Int64("a", 0, "anchor")

func main() {
	flag.Parse()
	fmt.Println("other", flag.Args())
	uid := *pUid
	roomId := *pRoomId
	upload := *pUpload
	if uid == 0 || roomId == 0 {
		fmt.Println(uid, roomId)
		return
	}

	fmt.Println("get pararms: ", roomId, uid, upload)

	StartTRPCClient(uid, roomId, upload)
}
