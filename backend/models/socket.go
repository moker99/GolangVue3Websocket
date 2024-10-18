package models

import (
	"context"
	"encoding/json"
	"fmt"
	"ginchat/utils"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"gopkg.in/fatih/set.v0"
)

// WebSocket 連接的節點
type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 用戶節點映射表
var clientMap map[int64]*Node = make(map[int64]*Node)

// 讀寫鎖
var rwLocker sync.RWMutex

// 全局 UDP 發送通道
var udpsendChan chan []byte = make(chan []byte, 1024)

// 載入聊天視窗時升級 WebSocket 連接
func UpgradeWebSocket(writer http.ResponseWriter, request *http.Request) (*Node, error) {
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	isvalidate := true // 假設 token 驗證通過
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalidate
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		fmt.Println("WebSocket 升級失敗: ", err)
		return nil, err
	}

	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	fmt.Printf("userId:%v初始化node節點成功", userId)

	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	go sendProc(node)
	go recvProc(node)

	return node, nil
}

// 初始化節點，但不升級 WebSocket 連接
func InitNode(userId int64) *Node {
	node := &Node{
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	return node
}

// 發送消息處理
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("8.步驟完成，websocket發送訊息 : ", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 儲存消息到 Redis
func storeMessageInRedis(data []byte) {
	ctx := context.Background()
	msg := Message{}
	json.Unmarshal(data, &msg)

	var key string
	if msg.FromId > msg.TargetId {
		key = fmt.Sprintf("msg_%d_%d_%d", msg.TargetId, msg.FromId, msg.Type)
	} else {
		key = fmt.Sprintf("msg_%d_%d_%d", msg.FromId, msg.TargetId, msg.Type)
	}

	// 存儲訊息到 Redis
	utils.Red.ZAdd(ctx, key, redis.Z{
		Score:  float64(time.Now().Unix()), // 使用時間戳作為 Score 來排序
		Member: data,
	})
	fmt.Println("2.訊息已寫入 Redis 緩存")
}

// 接收消息處理
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("1.recvProc接收訊息 : ", string(data))

		// 儲存消息到 Redis
		storeMessageInRedis(data)

		// 廣播消息
		broadMsg(data)
	}
}

// 廣播消息處理
func broadMsg(data []byte) {
	fmt.Println("3.將接收訊息放進UDP的管道 : ", string(data))
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("監聽UDP地址協程啟動 !!!")
}

// 發送UDP數據
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case data := <-udpsendChan:
			fmt.Println("4.UDP的管道收到訊息並發送至UDP地址 : ", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 接收UDP數據
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("5.監聽UDP地址接收訊息並分發 : ", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

// 根據訊息類型進行調度
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("6.解析訊息種類並調用指定的發送方法 : ", string(data))
	switch msg.Type {
	case 1: // 私信
		sendMsg(msg.TargetId, data)
	}
}

// 發送私信
func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	fmt.Println("7.將訊息給websocket發送私信 : ", string(msg))
	if ok {
		node.DataQueue <- msg
	}
}
