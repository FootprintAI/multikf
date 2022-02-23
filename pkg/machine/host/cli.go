package host

import (
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
	log "github.com/golang/glog"
)

type urlBinary struct {
	Os      string
	Kind    string
	Kubectl string
}

func NewCLI(binpath string, verbose bool) (*CLI, error) {
	if binpath == "" {
		binpath = os.TempDir()
	}
	if err := os.MkdirAll(binpath, os.ModePerm); err != nil {
		return nil, err
	}
	cli := &CLI{
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
	verbose                bool
	localKindBinaryPath    string
	localKubectlBinaryPath string
	urlBinary              urlBinary
}

func (cli *CLI) ensureBinaries() error {
	if !fileExists(cli.localKindBinaryPath) {
		log.Infof("can't found binary from %s, download from intenrnet...\n", cli.localKindBinaryPath)
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
		log.Infof("can't found binary from %s, download from intenrnet...\n", cli.localKubectlBinaryPath)
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
	_, stdout, err := cli.runCmd(cmdAndArgs)
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
	b := make([]byte, 0, 1024*1024*10 /*10M buffer*/)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
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
	_, stdout, err := cli.runCmd(cmdAndArgs)
	if err != nil {
		return err
	}
	return stdout.Stdout()
}

func (cli *CLI) InstallKubeflow(kfmanifestFile string) error {
	cmdAndArgs := []string{
		cli.localKubectlBinaryPath,
		"apply",
		"-f",
		kfmanifestFile,
	}
	for {
		log.Info("this command will keep retry for every 10s until it succeed.\n")
		exitcode, stdout, err := cli.runCmd(cmdAndArgs)
		if err != nil {
			return err
		}
		stdout.Stdout()
		if exitcode != 0 {
			time.Sleep(10)
		} else {
			return nil
		}
	}
}

func (cli *CLI) RemoveCluster(clustername string) error {
	cmdAndArgs := []string{
		cli.localKindBinaryPath,
		"delete",
		"cluster",
		"--name",
		clustername,
	}
	_, stdout, err := cli.runCmd(cmdAndArgs)
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
	_, stdout, err := cli.runCmd(cmdAndArgs)
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
	_, stdout, err := cli.runCmd(cmdAndArgs)
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
	_, stdout, err := cli.runCmd(cmdAndArgs)
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
	ch chan string
}

func newOutputStream(ch chan string) *outputStream {
	return &outputStream{ch: ch}
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
	log.Infof("cli: %s\n", string(all))
	return nil
}

func (cli *CLI) runCmd(cmdAndArgs []string) (int, *outputStream, error) {
	if cli.verbose {
		log.Infof("cmdandargs:%s\n", cmdAndArgs)
	}
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	runcmd := cmd.NewCmdOptions(cmdOptions, cmdAndArgs[0], cmdAndArgs[1:]...)
	runStatus := runcmd.Start()
	// run and output stderr
	for stderrline := range runcmd.Stderr {
		log.Infof("cli: %s\n", stderrline)
	}
	stat := <-runStatus

	return stat.Exit, newOutputStream(runcmd.Stdout), nil
}
