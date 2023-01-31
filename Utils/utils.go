package Utils

import (
	"Funnull/Moduls"
	"fmt"
	websocket2 "github.com/gorilla/websocket"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

func CommandExecute(session *ssh.Session, command string) string {
	//执行远程命令
	combo, err := session.CombinedOutput(command)
	if err != nil {
		return "error info:" + err.Error()
	}
	return string(combo)
}

func RemoteConnect(nodes Moduls.Nodes, command string, ch chan string, wg *sync.WaitGroup, ws_msg_type int, ws_conn *websocket2.Conn) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 15, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            nodes.NodeIP.SSHUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	if nodes.NodeIP.SSHPass != "" {
		config.Auth = []ssh.AuthMethod{ssh.Password(nodes.NodeIP.SSHPass)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc("~/.ssh/id_rsa")}
	}
	//config.Auth = []ssh.AuthMethod{ssh.Password(nodes.NodeIP.SSHPass)}
	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", nodes.NodeIP.SSHRemoteHost, nodes.NodeIP.SSHPort)
	sshClient, sc_err := ssh.Dial("tcp", addr, config)
	if sc_err != nil {
		er := ws_conn.WriteMessage(ws_msg_type, []byte("[Error info]:"+"节点："+nodes.NodeIP.SSHRemoteHost+"\n错误内容："+sc_err.Error()+"\n"))
		if er != nil {
			log.Fatalf("error info:" + er.Error())
		}
		ch <- "error info:" + "节点：" + nodes.NodeIP.SSHRemoteHost + sc_err.Error()
	}
	if sshClient != nil {
		//创建ssh-session
		session, err := sshClient.NewSession()
		if err != nil {
			er := ws_conn.WriteMessage(ws_msg_type, []byte("[Error info]:"+"节点："+nodes.NodeIP.SSHRemoteHost+"\n错误内容："+err.Error()+"\n"))
			if er != nil {
				log.Fatalf("error info:" + er.Error())
			}
			ch <- "error info:" + "节点：" + nodes.NodeIP.SSHRemoteHost + "\n错误内容：" + err.Error()
		}
		if session != nil {
			executeres := CommandExecute(session, command)
			er := ws_conn.WriteMessage(ws_msg_type, []byte(executeres+"\n"))
			if er != nil {
				log.Fatalf("error info:" + er.Error())
			}
			ch <- "[DEBUG INFO]:" + "节点：" + nodes.NodeIP.SSHRemoteHost + "\n" + executeres
		}
	}
	wg.Done()
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
