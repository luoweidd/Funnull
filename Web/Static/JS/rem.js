
function WebSocketTest() {
    SSH_User = document.getElementById("SSH_User").valueOf().value
    SSH_Pass = document.getElementById("SSH_Pass").valueOf().value
    SSH_Port = document.getElementById("SSH_Port").valueOf().value
    Node_Name_IP = document.getElementById("Node_Name_IP").valueOf().value
    req_data = [SSH_User,SSH_Pass,SSH_Port,Node_Name_IP]
    if ("WebSocket" in window) {
        var ws = new WebSocket("ws://127.0.0.1:8080/BatchNodeInstall");
        var msg;
        msg = document.getElementById("logs_output");
        ws.onopen = function () {
            // Web Socket 已连接上，使用 send() 方法发送数据
            ws.send(req_data);
            msg.value = "开始执行批量操作，打开网络链接，等待接收数据。\n";
        };
        ws.onmessage = function (evt) {
            var received_msg = evt.data;
            msg.value = msg.value+received_msg;
        };
        ws.onclose = function () {
            msg.value = msg.value+"\n执行完毕，关闭链接";
        };
    } else {
        // 浏览器不支持 WebSocket
        alert("您的浏览器不支持 WebSocket!");
    }
}