package BatchNodeInstall

import (
	"Funnull/Controller"
	"fmt"
	"github.com/gin-gonic/gin"
	websocket2 "github.com/gorilla/websocket"
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
		fmt.Printf("error info:" + err.Error())
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
		fmt.Printf("error info:" + err.Error())
	}
	data := strings.Split(string(req_data), ",")
	node_name_ip := data[3]
	node_user := data[0]
	node_pass := data[1]
	node_port, _ := strconv.Atoi(data[2])
	Controller.DataProc(node_name_ip, node_user, node_pass, node_port, msg_type, conn)
}
