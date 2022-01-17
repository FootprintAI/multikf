package runtime

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	govagrant "github.com/bmatcuk/go-vagrant"
	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

type SSHConfigFile struct {
	govagrant.SSHConfig
}

func (s SSHConfigFile) Addr() string {
	return fmt.Sprintf("%s:%d", s.SSHConfig.HostName, s.SSHConfig.Port)
}

func (s SSHConfigFile) PrivateKeySigner() (ssh.Signer, error) {
	key, err := ioutil.ReadFile(s.IdentityFile)
	if err != nil {
		return nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

func (s SSHConfigFile) SSHClientConfig() (*ssh.ClientConfig, error) {
	signer, err := s.PrivateKeySigner()
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config, nil
}

func NewSSHConn(hostAddr string, clientconfig *ssh.ClientConfig) (*SSHConn, error) {
	conn, err := ssh.Dial("tcp", hostAddr, clientconfig)
	if err != nil {
		return nil, err
	}
	return &SSHConn{Client: conn}, nil
}

type SSHConn struct {
	*ssh.Client
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
