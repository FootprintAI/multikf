package machine

import (
	"errors"
	"fmt"
	"net"

	log "github.com/golang/glog"
)

func FindFreeSSHPort() (int, error) {
	start := 2022
	nextIncr := 100
	for {
		if start > 65535 {
			return 0, errors.New("no available port")
		}
		log.Infof("sshport: try %d\n", start)
		if err := isPortAvaialble(start); err == nil {
			log.Infof("sshport: found %d\n", start)
			return start, nil
		}
		start = start + nextIncr
	}
}

func FindFreeKubeApiPort() (int, error) {
	start := 16443
	nextIncr := 1000
	for {
		if start > 65535 {
			return 0, errors.New("no available port")
		}
		log.Infof("kubeapiport: try %d\n", start)
		if err := isPortAvaialble(start); err == nil {
			log.Infof("kubeapiport: found %d\n", start)
			return start, nil
		}
		start = start + nextIncr
	}
}

func isPortAvaialble(port int) error {
	addr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	l.Close()
	return nil
}
