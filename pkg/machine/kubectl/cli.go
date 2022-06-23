package kubectl

import (
	"errors"
	"fmt"
	"io"
	pkgioutil "io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	gocmd "github.com/go-cmd/cmd"
	"sigs.k8s.io/kind/pkg/log"

	"github.com/footprintai/multikf/pkg/machine/cmd"
	"github.com/footprintai/multikf/pkg/machine/ioutil"
)

type binaryResource struct {
	Os      string
	Kind    string
	Kubectl string
}

func NewCLI(logger log.Logger, binpath string, verbose bool) (*CLI, error) {

	if binpath == "" {
		binpath = os.TempDir()
	}
	if err := os.MkdirAll(binpath, os.ModePerm); err != nil {
		return nil, err
	}
	cli := &CLI{
		logger:                 logger,
		verbose:                verbose,
		localKindBinaryPath:    filepath.Join(binpath, localBinaryRes.Kind),
		localKubectlBinaryPath: filepath.Join(binpath, localBinaryRes.Kubectl),
		urlBinary:              urlBinaryRes,
	}
	cli.logger.V(1).Infof("running binary with OS:%s...\n", localBinaryRes.Os)
	if err := cli.ensureBinaries(); err != nil {
		return nil, err
	}
	return cli, nil
}

type CLI struct {
	logger                 log.Logger
	verbose                bool
	localKindBinaryPath    string
	localKubectlBinaryPath string
	urlBinary              binaryResource
}

func (cli *CLI) ensureBinaries() error {
	if !fileExists(cli.localKindBinaryPath) {
		cli.logger.V(0).Infof("can't found binary from %s, download from intenrnet...\n", cli.localKindBinaryPath)
		// download kind
		if err := downloadPlainBinary(cli.urlBinary.Kind, cli.localKindBinaryPath); err != nil {
			return err
		}
		err := os.Chmod(cli.localKindBinaryPath, 0755 /*rwx-rx-rx*/)
		if err != nil {
			return err
		}
	}
	if !fileExists(cli.localKubectlBinaryPath) {
		cli.logger.V(0).Infof("can't found binary from %s, download from intenrnet...\n", cli.localKubectlBinaryPath)
		// download kubectl
		if err := downloadPlainBinary(cli.urlBinary.Kubectl, cli.localKubectlBinaryPath); err != nil {
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

func (cli *CLI) ListClusters() ([]string, error) {
	cmdAndArgs := []string{
		cli.localKindBinaryPath,
		"get",
		"clusters",
	}
	sr, _, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return nil, err
	}
	stdoutblob, err := ioutil.ReadAll(sr)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, token := range strings.Split(string(stdoutblob), "\n") {
		if token != "" {
			out = append(out, token)
		}
	}
	return out, nil
}

func (cli *CLI) ProvisonCluster(kindConfigfile string) error {
	cmdAndArgs := []string{
		cli.localKindBinaryPath,
		"create",
		"cluster",
		"--config",
		kindConfigfile,
	}
	sr, _, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	return sr.Stdout()
}

//func (cli *CLI) InstallRequiredPkgs(containername ContainerName) error {
//	// TODO: check whether we have to install gpu related pkg
//	//_, err := cli.RemoteExec(containername, "apt-get update && apt-get install -y pciutils")
//	//if err != nil {
//	//	log.Error("cli: failed on install required pkg, err:%v\n", err)
//	//}
//	//cli.RemoteExec(containername, "lspci  | grep -i nvidia")
//	return nil
//}

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
		sr.Stdout()

		ps := <-status
		cli.logger.V(1).Infof("kf installation, ps code:%+v\n", ps.Exit)
		if ps.Exit == 0 {
			return nil
		}
		time.Sleep(20 * time.Second)
	}
	return nil
}

func (cli *CLI) PatchKubeflow(kubeConfigFile string) error {
	multiCmdAndArgs := [][]string{
		[]string{
			cli.localKubectlBinaryPath,
			"patch",
			"configmap",
			"workflow-controller-configmap",
			"--patch",
			"{\"data\":{\"containerRuntimeExecutor\":\"emissary\"}}",
			"-n",
			"kubeflow",
			"--kubeconfig",
			kubeConfigFile,
		},
		[]string{
			cli.localKubectlBinaryPath,
			"rollout",
			"restart",
			"deployment/workflow-controller",
			"-n",
			"kubeflow",
			"--kubeconfig",
			kubeConfigFile,
		},
	}
	for _, cmdAndArgs := range multiCmdAndArgs {
		sr, _, err := cli.runCmd(cmdAndArgs)
		if err != nil {
			return err
		}
		sr.Stdout()
		time.Sleep(3 * time.Second)

	}
	return nil
}

func (cli *CLI) Portforward(kubeConfigFile, svc, namespace string, fromPort, toPort int) error {
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
	sr, _, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	return sr.Stdout()
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
	return sr.Stdout()
}

func (cli *CLI) RemoveCluster(clustername string) error {
	cmdAndArgs := []string{
		cli.localKindBinaryPath,
		"delete",
		"cluster",
		"--name",
		clustername,
	}
	sr, _, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	return sr.Stdout()

}

func (cli *CLI) GetKubeConfig(clustername string, exportLocalFilePath string) error {
	cmdAndArgs := []string{
		cli.localKindBinaryPath,
		"get",
		"kubeconfig",
		"--name",
		clustername,
	}
	sr, _, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	contentBlob, err := ioutil.ReadAll(sr)
	if err != nil {
		return err
	}
	return pkgioutil.WriteFile(exportLocalFilePath, contentBlob, 0600)
}

func (cli *CLI) runCmd(cmdAndArgs []string) (ioutil.StreamReader, <-chan gocmd.Status, error) {
	return cmd.NewCmd(cli.logger, cli.verbose).Run(cmdAndArgs...)
}
