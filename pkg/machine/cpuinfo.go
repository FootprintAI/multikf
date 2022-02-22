package machine

import (
	"regexp"
	"strconv"
	"strings"
)

func NewCpuInfoParserHelper(str string, err error) (*CpuInfo, error) {
	if err != nil {
		return nil, err
	}
	return NewCpuInfoParser(str)
}

var cpuinfoRegExp = regexp.MustCompile("([^:]*?)\\s*:\\s*(.*)$")

func NewCpuInfoParser(str string) (*CpuInfo, error) {
	lines := strings.Split(str, "\n")

	var cpuinfo = &CpuInfo{}
	var processor = newProcessorInfo()

	for i, line := range lines {
		var key string
		var value string

		if len(line) == 0 && i != len(lines)-1 {
			// end of processor
			cpuinfo.Processors = append(cpuinfo.Processors, processor)
			processor = newProcessorInfo()
			continue
		} else if i == len(lines)-1 {
			continue
		}

		submatches := cpuinfoRegExp.FindStringSubmatch(line)
		key = submatches[1]
		value = submatches[2]

		switch key {
		case "processor":
			processor.Id, _ = strconv.ParseInt(value, 10, 64)
		case "physical id":
			processor.PhysicalId, _ = strconv.ParseInt(value, 10, 64)
		case "core id":
			processor.CoreId, _ = strconv.ParseInt(value, 10, 64)
		}
	}
	if processor.Id != -1 {
		// valid but not append yet
		cpuinfo.Processors = append(cpuinfo.Processors, processor)
	}
	return cpuinfo, nil
}

func newProcessorInfo() *ProcessorInfo {
	return &ProcessorInfo{
		Id:         -1,
		CoreId:     -1,
		PhysicalId: -1,
	}
}

type CpuInfo struct {
	Processors []*ProcessorInfo
}

func (c *CpuInfo) NumCPUs() int {
	return len(c.Processors)
}

type ProcessorInfo struct {
	Id         int64
	PhysicalId int64
	CoreId     int64
}
