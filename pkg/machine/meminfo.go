package machine

import (
	"bufio"
	"bytes"
	"fmt"
)

func NewMemInfoParserHelper(str string, err error) (*MemInfo, error) {
	if err != nil {
		return nil, err
	}
	return NewMemInfoParser(str)
}

func NewMemInfoParser(str string) (*MemInfo, error) {
	var err error
	s := bufio.NewScanner(bytes.NewBufferString(str))
	m := new(MemInfo)
	fieldOfInterested := 4
	for s.Scan() && fieldOfInterested > 0 {
		switch {
		case bytes.HasPrefix(s.Bytes(), []byte(`MemTotal:`)):
			_, err = fmt.Sscanf(s.Text(), "MemTotal:%d", &m.total)
		case bytes.HasPrefix(s.Bytes(), []byte(`MemFree:`)):
			_, err = fmt.Sscanf(s.Text(), "MemFree:%d", &m.free)
		case bytes.HasPrefix(s.Bytes(), []byte(`Buffers:`)):
			_, err = fmt.Sscanf(s.Text(), "Buffers:%d", &m.buffers)
		case bytes.HasPrefix(s.Bytes(), []byte(`Cached:`)):
			_, err = fmt.Sscanf(s.Text(), "Cached:%d", &m.cached)
		default:
			continue
		}
		if err != nil {
			return nil, err
		}
		fieldOfInterested--
	}
	return m, nil
}

type MemInfo struct {
	free    uint32 // in kB
	cached  uint32
	buffers uint32
	total   uint32
}

const bytesInM = 1024

func (m *MemInfo) Free() string {
	return fmt.Sprintf("%.2f Mib", float64(m.free)/bytesInM)
}

func (m *MemInfo) Total() string {
	return fmt.Sprintf("%.2f Mib", float64(m.total)/bytesInM)
}
