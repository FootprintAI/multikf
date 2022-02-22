package ssh

import (
	"bytes"
	"context"
	"os"
	"time"

	"github.com/bramvdbogaerde/go-scp"
	gossh "golang.org/x/crypto/ssh"
)

func NewSSHConn(hostAddr string, clientconfig *gossh.ClientConfig) (*SSHConn, error) {
	conn, err := gossh.Dial("tcp", hostAddr, clientconfig)
	if err != nil {
		return nil, err
	}
	return &SSHConn{Client: conn}, nil
}

type SSHConn struct {
	*gossh.Client
}

func (s *SSHConn) Scp(fromRemotePath string, toHostPath string) error {
	client, err := scp.NewClientBySSH(s.Client)
	defer client.Close()

	f, err := os.OpenFile(toHostPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := client.CopyFromRemote(timeoutCtx, f, fromRemotePath); err != nil {
		return err
	}
	return nil
}

func (s *SSHConn) Exec(command string) (string, error) {
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		return "", err
	}
	return b.String(), nil
}
