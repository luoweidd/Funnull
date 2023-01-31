package Controller

import (
	"Funnull/Moduls"
	"Funnull/Utils"
	websocket2 "github.com/gorilla/websocket"
	"log"
	"net"
	"strings"
	"sync"
)

func BatchNodeCommandExecute(node_name_ip string, node_user string, node_pass string, node_port int, ch chan string, Command string, ws_msg_type int, ws_conn *websocket2.Conn) {
	node_list := strings.Split(node_name_ip, "\n")
	wg := sync.WaitGroup{}
	wg.Add(len(node_list))
	for _, v := range node_list {
		nodeinfo := strings.Split(v, " ")
		ip_addr := net.ParseIP(nodeinfo[1])
		if ip_addr != nil {
			node := Moduls.Nodes{nodeinfo[1], Moduls.SSHInfo{nodeinfo[1], node_user, node_pass, node_port}}
			go Utils.RemoteConnect(node, Command, ch, &wg, ws_msg_type, ws_conn)
		} else {
			err := ws_conn.WriteMessage(ws_msg_type, []byte("Node_IP 格式有误，请修正后再执行。"+"\n"))
			if err != nil {
				log.Fatalf("error info:" + err.Error())
			}
		}
	}
	wg.Wait()
	close(ch)
}

func DataProc(node_name_ip string, node_user string, node_pass string, node_port int, Command string, ws_msg_type int, ws_conn *websocket2.Conn) {
	node_list := strings.Split(node_name_ip, "\n")
	chan_num := len(node_list)
	ch := make(chan string, chan_num)
	BatchNodeCommandExecute(node_name_ip, node_user, node_pass, node_port, ch, Command, ws_msg_type, ws_conn)
	er := ws_conn.WriteMessage(ws_msg_type, []byte("错误列表：\n"))
	if er != nil {
		log.Fatalf("error info:" + er.Error())
	}
	for i := range ch {
		if strings.Contains(i, "error info:") {
			err := ws_conn.WriteMessage(ws_msg_type, []byte(i+"\n"))
			if err != nil {
				log.Fatalf("error info:" + err.Error())
			}
		}
	}
	err := ws_conn.WriteMessage(ws_msg_type, []byte("本次任务执行完成，服务端关闭连接！\n"))
	if err != nil {
		log.Fatalf("error info:" + err.Error())
	}

	ws_conn.WriteMessage(websocket2.CloseMessage,
		websocket2.FormatCloseMessage(websocket2.CloseNormalClosure, "本次任务执行完成，服务端关闭连接！"))
	// 真正关闭 conn
	ws_conn.Close()
}
