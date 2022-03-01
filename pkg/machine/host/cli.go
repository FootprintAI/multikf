package host

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
	"sigs.k8s.io/kind/pkg/log"
)

type urlBinary struct {
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
		localKindBinaryPath:    filepath.Join(binpath, "kind"),
		localKubectlBinaryPath: filepath.Join(binpath, "kubectl"),
		urlBinary:              urlBinaryRes,
	}
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
	urlBinary              urlBinary
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
	stdout, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return nil, err
	}
	stdoutblob, err := readall(stdout)
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

func readall(r io.Reader) ([]byte, error) {
	buf := &bytes.Buffer{}
	b := make([]byte, 1024*1024*10 /*10M buffer*/)
	for {
		n, err := r.Read(b)
		if err == nil {
			buf.Write(b[:n])
		} else {
			if err == io.EOF {
				return buf.Bytes(), nil
			} else {
				return nil, err
			}
		}
	}
}

func (cli *CLI) ProvisonCluster(kindConfigfile string) error {
	cmdAndArgs := []string{
		cli.localKindBinaryPath,
		"create",
		"cluster",
		"--config",
		kindConfigfile,
	}
	stdout, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	return stdout.Stdout()
}

func (cli *CLI) InstallRequiredPkgs(containername ContainerName) error {
	// TODO: check whether we have to install gpu related pkg
	//_, err := cli.RemoteExec(containername, "apt-get update && apt-get install -y pciutils")
	//if err != nil {
	//	log.Error("cli: failed on install required pkg, err:%v\n", err)
	//}
	//cli.RemoteExec(containername, "lspci  | grep -i nvidia")
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
	cli.logger.V(0).Info("this command will keep retry 2 times for every 30s.\n")
	for i := 0; i < 3; i++ {
		stdout, err := cli.runCmd(cmdAndArgs)
		if err != nil {
			return err
		}
		stdout.Stdout()
		time.Sleep(30 * time.Second)
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
		stdout, err := cli.runCmd(cmdAndArgs)
		if err != nil {
			return err
		}
		stdout.Stdout()
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
	stdout, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	return stdout.Stdout()
}

func (cli *CLI) RemoveCluster(clustername string) error {
	cmdAndArgs := []string{
		cli.localKindBinaryPath,
		"delete",
		"cluster",
		"--name",
		clustername,
	}
	stdout, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	return stdout.Stdout()

}

type dockerState struct {
	Status string `json:"status"`
}

func (cli *CLI) GetClusterStatus(containername ContainerName) (string, error) {
	cmdAndArgs := []string{
		"docker",
		"inspect",
		containername.Name(),
		"--format='{{json .State}}'",
	}
	stdout, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return "", err
	}
	d := dockerState{}
	blob, _ := readall(stdout)
	stripped := blob[1 : len(blob)-2] // remove ' xxx '\n
	if err := json.Unmarshal(stripped, &d); err != nil {
		return "", err
	}
	return d.Status, nil
}

func (cli *CLI) RemoteExec(containername ContainerName, cmd string) (resp string, err error) {
	cmdAndArgs := []string{
		"docker",
		"exec",
		containername.Name(),
		"sh",
		"-c",
		cmd,
	}
	stdout, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return "", err
	}
	all, _ := readall(stdout)
	return string(all), nil
}

func (cli *CLI) GetKubeConfig(clustername string, exportLocalFilePath string) error {
	cmdAndArgs := []string{
		cli.localKindBinaryPath,
		"get",
		"kubeconfig",
		"--name",
		clustername,
	}
	stdout, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	contentBlob, err := readall(stdout)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(exportLocalFilePath, contentBlob, 0600)
}

// outputStream implements io.Reader by wrapping the line channel
type outputStream struct {
	logger log.Logger
	ch     chan string
}

func newOutputStream(logger log.Logger, ch chan string) *outputStream {
	return &outputStream{ch: ch, logger: logger}
}

func (o *outputStream) Read(b []byte) (int, error) {
	out, more := <-o.ch
	if !more {
		return 0, io.EOF
	}
	if len(out) > len(b) {
		panic(fmt.Sprintf("insufficient buffer size(buf:%d, data:%d), data could be lost", len(b), len(out)))
	}
	n := copy(b[:len(b)-1], []byte(out))
	b[n] = '\n'
	return n + 1, nil
}

func (o *outputStream) Stdout() error {
	all, err := readall(o)
	if err != nil {
		return err
	}
	o.logger.V(0).Infof("%s\n", string(all))
	return nil
}

func (cli *CLI) runCmd(cmdAndArgs []string) (*outputStream, error) {
	if cli.verbose {
		cli.logger.V(0).Infof("cmd->%s\n", cmdAndArgs)
	}
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	runcmd := cmd.NewCmdOptions(cmdOptions, cmdAndArgs[0], cmdAndArgs[1:]...)
	runcmd.Start()
	// run and output stderr
	for stderrline := range runcmd.Stderr {
		cli.logger.V(1).Infof("%s\n", stderrline)
	}
	//stat := <-runStatus

	return newOutputStream(cli.logger, runcmd.Stdout), nil
}
