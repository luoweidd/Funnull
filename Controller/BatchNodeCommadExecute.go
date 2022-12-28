package Controller

import (
	"Funnull/Moduls"
	"Funnull/Utils"
	"fmt"
	websocket2 "github.com/gorilla/websocket"
	"strings"
)

func BatchNodeCommandExecute(nodes Moduls.Nodes, c chan string) {
	exeres := Utils.RemoteConnect(nodes, "sudo curl -s --basic -u rjzs:lt5uf https://ac.funnullv8.com/agent_install.sh | sudo bash")
	c <- exeres
}

func DataProc(node_name_ip string, node_user string, node_pass string, node_port int, ws_msg_type int, ws_conn *websocket2.Conn) {
	node_list := strings.Split(node_name_ip, "\n")
	for _, v := range node_list {
		nodeinfo := strings.Split(v, " ")
		node := Moduls.Nodes{nodeinfo[0], Moduls.SSHInfo{nodeinfo[1], node_user, node_pass, node_port}}
		ch := make(chan string)
		go BatchNodeCommandExecute(node, ch)
		res_txt := <-ch
		data := []byte(res_txt + "\n")
		err := ws_conn.WriteMessage(ws_msg_type, data)
		if err != nil {
			fmt.Printf("error info:" + err.Error())
		}
	}
}
