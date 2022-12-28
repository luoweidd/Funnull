package Moduls

type SSHInfo struct {
	SSHRemoteHost string
	SSHUser       string
	SSHPass       string
	SSHPort       int
}

type Nodes struct {
	NodeName string
	NodeIP   SSHInfo
}
