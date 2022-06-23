package vagrantclient

import (
	"fmt"
	"io/ioutil"

	govagrant "github.com/bmatcuk/go-vagrant"
	fssh "github.com/footprintai/multikf/pkg/ssh"
	"golang.org/x/crypto/ssh"
	"sigs.k8s.io/kind/pkg/log"
)

func NewVagrantCli(machineName string, vagrantMachineDir string, logger log.Logger, verbose bool) (*VagrantCli, error) {
	cli, err := govagrant.NewVagrantClient(vagrantMachineDir)
	if err != nil {
		return nil, err
	}
	return &VagrantCli{
		logger:  logger,
		name:    machineName,
		client:  cli,
		Verbose: verbose,
	}, nil
}

type VagrantCli struct {
	name    string
	logger  log.Logger
	client  *govagrant.VagrantClient
	Verbose bool
}

func (v *VagrantCli) Up() error {
	v.logger.V(0).Infof("vagrantmachine(%s): start machines...\n", v.name)
	cmd := v.client.Up()
	cmd.MachineName = v.name
	cmd.Verbose = v.Verbose
	if err := cmd.Run(); err != nil {
		return err
	}
	v.logger.V(0).Infof("vagrantmahcine(%s) is ready\n", v.name)
	return nil
}

type vagrantStatus string

func (vs vagrantStatus) String() string {
	return string(vs)
}

const (
	vagrantStatusInvalid    vagrantStatus = "invalid_status"
	vagrantStatusNotCreated vagrantStatus = "not_created"
	vagrantStatusUp         vagrantStatus = "up"
	vagrantStatusRunning    vagrantStatus = "running"
)

func (v *VagrantCli) Status() string {
	cmd := v.client.Status()
	cmd.MachineName = v.name
	cmd.Verbose = v.Verbose
	if err := cmd.Run(); err != nil {
		return vagrantStatusInvalid.String()
	}
	return cmd.StatusResponse.Status[v.name]
}

func (v *VagrantCli) TryUp() error {
	status := v.Status()
	v.logger.V(0).Infof("vagrantmachine(%s): status:%s\n", v.name, status)

	if status == vagrantStatusNotCreated.String() || status == vagrantStatusInvalid.String() {
		return v.Up()
	}
	if status == vagrantStatusUp.String() || status == vagrantStatusRunning.String() {
		v.logger.V(0).Infof("vagrantmachine(%s) already up and running\n", v.name)
		return nil
	}
	v.logger.V(0).Infof("vagrantmachine(%s) clean up vagrant previous state\n", v.name)
	v.Destroy()
	return v.TryUp()
}

func (v *VagrantCli) Destroy() error {
	cmd := v.client.Destroy()
	cmd.Verbose = v.Verbose
	cmd.MachineName = v.name
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (v *VagrantCli) SSHConfig() (SSHConfigFile, error) {
	cmd := v.client.SSHConfig()
	cmd.Verbose = v.Verbose
	cmd.MachineName = v.name
	if err := cmd.Run(); err != nil {
		return SSHConfigFile{}, err
	}
	return SSHConfigFile{cmd.SSHConfigResponse.Configs[v.name]}, nil
}

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

func (v *VagrantCli) Scp(fromRemotePath string, toHostPath string) error {
	sshconfg, err := v.SSHConfig()
	if err != nil {
		return err
	}
	clientconfig, err := sshconfg.SSHClientConfig()
	if err != nil {
		return err
	}

	conn, err := fssh.NewSSHConn(sshconfg.Addr(), clientconfig)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Scp(fromRemotePath, toHostPath)
}

func (v *VagrantCli) SshExec(command string) (string, error) {
	sshconfg, err := v.SSHConfig()
	if err != nil {
		return "", err
	}
	clientconfig, err := sshconfg.SSHClientConfig()
	if err != nil {
		return "", err
	}

	conn, err := fssh.NewSSHConn(sshconfg.Addr(), clientconfig)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return conn.Exec(command)
}
