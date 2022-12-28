package Utils

import (
	"Funnull/Moduls"
	"fmt"
	"golang.org/x/crypto/ssh"
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

func RemoteConnect(nodes Moduls.Nodes, command string) string {
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            nodes.NodeIP.SSHUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(nodes.NodeIP.SSHPass)}
	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", nodes.NodeIP.SSHRemoteHost, nodes.NodeIP.SSHPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return "error info:" + err.Error()
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		return "error info:" + err.Error()
	}
	defer session.Close()
	executeres := CommandExecute(session, command)
	return executeres
}
