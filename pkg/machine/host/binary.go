package host

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/go-cmd/cmd"
	log "github.com/golang/glog"
)

type urlBinary struct {
	Os      string
	Kind    string
	Kubectl string
}

func NewBinary() *Binary {
	return &Binary{
		localKindBinaryPath:    filepath.Join(os.TempDir(), "kind"),
		localKubectlBinaryPath: filepath.Join(os.TempDir(), "kubectl"),
		urlBinary:              urlBinaryRes,
	}
}

type Binary struct {
	localKindBinaryPath    string
	localKubectlBinaryPath string
	urlBinary              urlBinary
}

func (b *Binary) ensureBinaries() error {
	if !fileExists(b.localKindBinaryPath) {
		// download kind
		if err := b.downloadPlainBinary(b.urlBinary.Kind, b.localKindBinaryPath); err != nil {
			return err
		}
	}
	if !fileExists(b.localKubectlBinaryPath) {
		// download kubectl
		if err := b.downloadPlainBinary(b.urlBinary.Kubectl, b.localKubectlBinaryPath); err != nil {
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

func (b *Binary) downloadPlainBinary(sourceURL, localpath string) error {
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
	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	log.Infof("%s: %d bytes downloaded\n", localpath, n)
	return nil
}

func (b *Binary) Kind(args []string) error {
	if err := b.ensureBinaries(); err != nil {
		return err
	}

	cmdAndArgs := []string{b.localKindBinaryPath}
	cmdAndArgs = append(cmdAndArgs, args)
	return b.runCmd(cmdAndArgs, true)
}

func (b *Binary) runCmd(cmdAndArgs []string, blocking bool) error {
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}
	runcmd := cmd.NewCmdOptions(cmdOptions, cmdAndArgs...)
	runStatus := runcmd.Start()
	if blocking {
		<-runStatus
	}
	stdoutCh := runcmd.Stdout
	stderrCh := runcmd.Stderr

	ch := mergeCh(stdoutCh, stderrCh)
	for {
		msg, more := <-ch
		if !more {
			break
		}
		log.Infof("bianry: %s\n", msg)
	}
	return nil
}

func mergeCh(ch1, ch2 chan string) chan string {
	out := make(chan string)
	dequeChFunc := func(outChan, inputChan chan string, done chan<- struct{}) {
		for {
			deque, more := <-inputChan
			if more {
				outChan <- deque
			} else {
				done <- struct{}{}
				return
			}
		}
	}

	go func() {
		done := make(chan struct{}, 2)

		go dequeChFunc(out, ch1, done)
		go dequeChFunc(out, ch2, done)

		<-done
		<-done
		close(out)
	}()
	return out
}
