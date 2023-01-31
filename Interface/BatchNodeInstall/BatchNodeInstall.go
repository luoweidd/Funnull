package BatchNodeInstall

import (
	"Funnull/Controller"
	"github.com/gin-gonic/gin"
	websocket2 "github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func BatchNodeInstall(c *gin.Context) {
	websocket := websocket2.Upgrader{ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	r := c.Request
	w := c.Writer
	conn, err := websocket.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("error info:" + err.Error())
	}

	defer func() {
		// 发送 websocket 结束包
		conn.WriteMessage(websocket2.CloseMessage,
			websocket2.FormatCloseMessage(websocket2.CloseNormalClosure, ""))
		// 真正关闭 conn
		conn.Close()
	}()
	msg_type, req_data, err := conn.ReadMessage()
	if err != nil {
		log.Fatalf("error info:" + err.Error())
	}
	if string(req_data) != "ping" {
		data := strings.Split(string(req_data), ",")
		if data[3] == "" || data[0] == "" || data[2] == "" || data[4] == "" {
			data_err := conn.WriteMessage(msg_type, []byte("数据内容有缺失，请补齐内容，除密码外均为必填项！\n"))
			if data_err != nil {
				log.Fatalf("error info:" + data_err.Error())
			}
		} else {
			node_name_ip := data[3]
			node_user := data[0]
			node_pass := data[1]
			node_port, _ := strconv.Atoi(data[2])
			Command := data[4]
			Controller.DataProc(node_name_ip, node_user, node_pass, node_port, Command, msg_type, conn)
		}

	} else {
		pone_err := conn.WriteMessage(msg_type, []byte("pong"))
		if pone_err != nil {
			log.Fatalf("error info:" + pone_err.Error())
		}
	}
}
