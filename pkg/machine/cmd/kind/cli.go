package kind

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	gocmd "github.com/go-cmd/cmd"
	"sigs.k8s.io/kind/pkg/log"

	"github.com/footprintai/multikf/pkg/machine/cmd"
	"github.com/footprintai/multikf/pkg/machine/ioutil"
)

func NewCLI(logger log.Logger, binpath string, verbose bool) (*CLI, error) {

	if binpath == "" {
		binpath = os.TempDir()
	}
	if err := os.MkdirAll(binpath, os.ModePerm); err != nil {
		return nil, err
	}
	cli := &CLI{
		logger:              logger,
		verbose:             verbose,
		localKindBinaryPath: filepath.Join(binpath, cmd.OSLocalBinaryRes.Kind),
		urlBinary:           cmd.OSUrlBinaryRes,
	}
	cli.logger.V(1).Infof("running binary with OS:%s...\n", cmd.OSLocalBinaryRes.Os)
	if err := cli.ensureBinaries(); err != nil {
		return nil, err
	}
	return cli, nil
}

type CLI struct {
	logger              log.Logger
	verbose             bool
	localKindBinaryPath string
	urlBinary           cmd.BinaryResource
}

func (cli *CLI) ensureBinaries() error {
	if !fileExists(cli.localKindBinaryPath) {
		cli.logger.V(0).Infof("can't found binary from %s, download from intenrnet...\n", cli.localKindBinaryPath)
		// download kind
		if err := cmd.DownloadPlainBinary(cli.urlBinary.Kind, cli.localKindBinaryPath); err != nil {
			return err
		}
		err := os.Chmod(cli.localKindBinaryPath, 0755 /*rwx-rx-rx*/)
		if err != nil {
			return err
		}
	}
	return nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
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
	return ioutil.StderrOnError(sr)
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
	return ioutil.StderrOnError(sr)

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
	return os.WriteFile(exportLocalFilePath, contentBlob, 0600)
}

func (cli *CLI) runCmd(cmdAndArgs []string) (*ioutil.CmdOutputStream, <-chan gocmd.Status, error) {
	return cmd.NewCmd(cli.logger).Run(cmdAndArgs...)
}
