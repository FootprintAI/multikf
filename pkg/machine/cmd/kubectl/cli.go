package kubectl

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	gocmd "github.com/go-cmd/cmd"
	"sigs.k8s.io/kind/pkg/log"

	"github.com/footprintai/multikf/pkg/k8s"
	"github.com/footprintai/multikf/pkg/machine/cmd"
	"github.com/footprintai/multikf/pkg/machine/ioutil"
)

func NewCLI(logger log.Logger, binpath string, verbose bool, version k8s.KindK8sVersion) (*CLI, error) {

	if binpath == "" {
		binpath = os.TempDir()
	}
	if err := os.MkdirAll(binpath, os.ModePerm); err != nil {
		return nil, err
	}
	cli := &CLI{
		logger:                 logger,
		verbose:                verbose,
		localKubectlBinaryPath: filepath.Join(binpath, cmd.OSLocalBinaryRes.Kubectl(k8s.KindK8sVersion{})),
		urlBinary:              cmd.OSUrlBinaryRes,
		k8sVersion:             version,
	}
	cli.logger.V(1).Infof("running binary with OS:%s...\n", cmd.OSLocalBinaryRes.Os)
	if err := cli.ensureBinaries(); err != nil {
		return nil, err
	}
	return cli, nil
}

type CLI struct {
	logger                 log.Logger
	verbose                bool
	localKubectlBinaryPath string
	urlBinary              cmd.BinaryResource
	k8sVersion             k8s.KindK8sVersion
}

func (cli *CLI) ensureBinaries() error {

	if !fileExists(cli.localKubectlBinaryPath) {
		cli.logger.V(0).Infof("can't found binary from %s, download from intenrnet...\n", cli.localKubectlBinaryPath)
		// download kubectl
		if err := downloadPlainBinary(cli.urlBinary.Kubectl(cli.k8sVersion), cli.localKubectlBinaryPath); err != nil {
			return err
		}
		err := os.Chmod(cli.localKubectlBinaryPath, 0755 /*rwx-rx-rx*/)
		if err != nil {
			return err
		}
	}
	// TODO: ensure version
	return nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func downloadPlainBinary(sourceURL, localpath string) error {
	resp, err := http.Get(sourceURL)
	if err != nil {
		return err
	}
	out, err := os.Create(localpath)
	if err != nil {
		return err
	}
	defer out.Close()
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (cli *CLI) InstallKubeflow(kubeConfigFile string, kfmanifestFile string) error {
	cmdAndArgs := []string{
		cli.localKubectlBinaryPath,
		"apply",
		"-f",
		kfmanifestFile,
		"--kubeconfig",
		kubeConfigFile,
	}
	cli.logger.V(0).Info("this command will keep retry every 20s.\n")
	for i := 0; i < 20; i++ {
		sr, status, err := cli.runCmd(cmdAndArgs)
		if err != nil {
			return err
		}
		ioutil.StderrOnError(sr)

		ps := <-status
		cli.logger.V(1).Infof("kf installation, ps code:%+v\n", ps.Exit)
		if ps.Exit == 0 {
			return nil
		}
		time.Sleep(20 * time.Second)
	}
	return nil
}

func (cli *CLI) RemoveKubeflow(kubeConfigFile string, kfmanifestFile string) error {
	cmdAndArgs := []string{
		cli.localKubectlBinaryPath,
		"delete",
		"-f",
		kfmanifestFile,
		"--kubeconfig",
		kubeConfigFile,
	}
	sr, _, err := cli.runCmd(cmdAndArgs)
	ioutil.StderrOnError(sr)
	return err
}

func (cli *CLI) Portforward(kubeConfigFile, svc, namespace string, address string, fromPort, toPort int) error {
	// TODO: auto reconnect
	cmdAndArgs := []string{
		cli.localKubectlBinaryPath,
		"port-forward",
		svc,
		"-n",
		namespace,
		fmt.Sprintf("%d:%d", toPort, fromPort),
		"--kubeconfig",
		kubeConfigFile,
	}
	if len(address) != 0 {
		cmdAndArgs = append(cmdAndArgs, []string{
			"--address",
			address,
		}...)
	}
	sr, _, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	return ioutil.StderrOnError(sr)
}

func (cli *CLI) GetPods(kindConfigfile string, namespace string) error {
	cmdAndArgs := []string{
		cli.localKubectlBinaryPath,
		"get",
		"pods",
		"--kubeconfig",
		kindConfigfile,
	}
	if namespace == "" {
		cmdAndArgs = append(cmdAndArgs, "--all-namespaces")
	} else {
		cmdAndArgs = append(cmdAndArgs, "--namespace", namespace)
	}
	sr, _, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	return ioutil.StderrOnError(sr)
}

func (cli *CLI) runCmd(cmdAndArgs []string) (*ioutil.CmdOutputStream, <-chan gocmd.Status, error) {
	return cmd.NewCmd(cli.logger).Run(cmdAndArgs...)
}
