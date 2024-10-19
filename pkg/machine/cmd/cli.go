package cmd

import (
	"io"
	"net/http"
	"os"

	"github.com/footprintai/multikf/pkg/k8s"
)

type BinaryResource struct {
	Os      string
	Kind    string
	Kubectl func(k8s.KindK8sVersion) string
}

func DownloadPlainBinary(sourceURL, localpath string) error {
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
