package machine

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// GPU struct is referred to https://github.com/rai-project/nvidia-smi/blob/864eb441c9ae3171601e2e9df942b9ca386612c0/model.go
type gpu struct {
	ID                                                     string   `xml:"id,attr" json:"id"`
	MemClockClocksGpu                                      string   `xml:"clocks>mem_clock" json:"mem_clock_clocks_gpu"`
	L1Cache                                                string   `xml:"ecc_errors>volatile>single_bit>l1_cache" json:"l_1_cache"`
	ProductName                                            string   `xml:"product_name" json:"product_name"`
	FreeFbMemoryUsageGpu                                   string   `xml:"fb_memory_usage>free" json:"free_fb_memory_usage_gpu"`
	PowerState                                             string   `xml:"power_readings>power_state" json:"power_state"`
	Free                                                   string   `xml:"bar1_memory_usage>free" json:"free"`
	RetiredCountDoubleBitRetirementRetiredPagesGpu         string   `xml:"retired_pages>double_bit_retirement>retired_count" json:"retired_count_double_bit_retirement_retired_pages_gpu"`
	ClocksThrottleReasonUnknown                            string   `xml:"clocks_throttle_reasons>clocks_throttle_reason_unknown" json:"clocks_throttle_reason_unknown"`
	ClocksThrottleReasonApplicationsClocksSetting          string   `xml:"clocks_throttle_reasons>clocks_throttle_reason_applications_clocks_setting" json:"clocks_throttle_reason_applications_clocks_setting"`
	Processes                                              string   `xml:"processes" json:"processes"`
	MemClockApplicationsClocksGpu                          string   `xml:"applications_clocks>mem_clock" json:"mem_clock_applications_clocks_gpu"`
	L2CacheSingleBitAggregateEccErrorsGpu                  string   `xml:"ecc_errors>aggregate>single_bit>l2_cache" json:"l_2_cache_single_bit_aggregate_ecc_errors_gpu"`
	CurrentLinkGen                                         string   `xml:"pci>pci_gpu_link_info>pcie_gen>current_link_gen" json:"current_link_gen"`
	TotalSingleBitVolatileEccErrorsGpu                     string   `xml:"ecc_errors>volatile>single_bit>total" json:"total_single_bit_volatile_ecc_errors_gpu"`
	TextureMemoryDoubleBitVolatileEccErrorsGpu             string   `xml:"ecc_errors>volatile>double_bit>texture_memory" json:"texture_memory_double_bit_volatile_ecc_errors_gpu"`
	L1CacheSingleBitAggregateEccErrorsGpu                  string   `xml:"ecc_errors>aggregate>single_bit>l1_cache" json:"l_1_cache_single_bit_aggregate_ecc_errors_gpu"`
	PendingGom                                             string   `xml:"gpu_operation_mode>pending_gom" json:"pending_gom"`
	AutoBoostDefault                                       string   `xml:"clock_policy>auto_boost_default" json:"auto_boost_default"`
	GraphicsClockApplicationsClocksGpu                     string   `xml:"applications_clocks>graphics_clock" json:"graphics_clock_applications_clocks_gpu"`
	PciBusID                                               string   `xml:"pci>pci_bus_id" json:"pci_bus_id"`
	PowerManagement                                        string   `xml:"power_readings>power_management" json:"power_management"`
	DeviceMemoryDoubleBitAggregateEccErrorsGpu             string   `xml:"ecc_errors>aggregate>double_bit>device_memory" json:"device_memory_double_bit_aggregate_ecc_errors_gpu"`
	BoardID                                                string   `xml:"board_id" json:"board_id"`
	DeviceMemoryDoubleBitVolatileEccErrorsGpu              string   `xml:"ecc_errors>volatile>double_bit>device_memory" json:"device_memory_double_bit_volatile_ecc_errors_gpu"`
	SupportedGraphicsClock                                 []string `xml:"supported_clocks>supported_mem_clock>supported_graphics_clock" json:"supported_graphics_clock"`
	PersistenceMode                                        string   `xml:"persistence_mode" json:"persistence_mode"`
	MemClock                                               string   `xml:"max_clocks>mem_clock" json:"mem_clock"`
	GraphicsClockClocksGpu                                 string   `xml:"clocks>graphics_clock" json:"graphics_clock_clocks_gpu"`
	Used                                                   string   `xml:"bar1_memory_usage>used" json:"used"`
	ImgVersion                                             string   `xml:"inforom_version>img_version" json:"img_version"`
	UsedFbMemoryUsageGpu                                   string   `xml:"fb_memory_usage>used" json:"used_fb_memory_usage_gpu"`
	TotalDoubleBitAggregateEccErrorsGpu                    string   `xml:"ecc_errors>aggregate>double_bit>total" json:"total_double_bit_aggregate_ecc_errors_gpu"`
	MinorNumber                                            string   `xml:"minor_number" json:"minor_number"`
	ProductBrand                                           string   `xml:"product_brand" json:"product_brand"`
	GraphicsClockDefaultApplicationsClocksGpu              string   `xml:"default_applications_clocks>graphics_clock" json:"graphics_clock_default_applications_clocks_gpu"`
	TotalFbMemoryUsageGpu                                  string   `xml:"fb_memory_usage>total" json:"total_fb_memory_usage_gpu"`
	RegisterFileDoubleBitVolatileEccErrorsGpu              string   `xml:"ecc_errors>volatile>double_bit>register_file" json:"register_file_double_bit_volatile_ecc_errors_gpu"`
	MinPowerLimit                                          string   `xml:"power_readings>min_power_limit" json:"min_power_limit"`
	TxUtil                                                 string   `xml:"pci>tx_util" json:"tx_util"`
	TextureMemory                                          string   `xml:"ecc_errors>volatile>single_bit>texture_memory" json:"texture_memory"`
	RegisterFileDoubleBitAggregateEccErrorsGpu             string   `xml:"ecc_errors>aggregate>double_bit>register_file" json:"register_file_double_bit_aggregate_ecc_errors_gpu"`
	PerformanceState                                       string   `xml:"performance_state" json:"performance_state"`
	CurrentDm                                              string   `xml:"driver_model>current_dm" json:"current_dm"`
	PciDeviceID                                            string   `xml:"pci>pci_device_id" json:"pci_device_id"`
	AccountedProcesses                                     string   `xml:"accounted_processes" json:"accounted_processes"`
	PendingRetirement                                      string   `xml:"retired_pages>pending_retirement" json:"pending_retirement"`
	TotalDoubleBitVolatileEccErrorsGpu                     string   `xml:"ecc_errors>volatile>double_bit>total" json:"total_double_bit_volatile_ecc_errors_gpu"`
	UUID                                                   string   `xml:"uuid" json:"uuid"`
	PowerLimit                                             string   `xml:"power_readings>power_limit" json:"power_limit"`
	ClocksThrottleReasonHwSlowdown                         string   `xml:"clocks_throttle_reasons>clocks_throttle_reason_hw_slowdown" json:"clocks_throttle_reason_hw_slowdown"`
	BridgeChipFw                                           string   `xml:"pci>pci_bridge_chip>bridge_chip_fw" json:"bridge_chip_fw"`
	ReplayCounter                                          string   `xml:"pci>replay_counter" json:"replay_counter"`
	L2CacheDoubleBitAggregateEccErrorsGpu                  string   `xml:"ecc_errors>aggregate>double_bit>l2_cache" json:"l_2_cache_double_bit_aggregate_ecc_errors_gpu"`
	ComputeMode                                            string   `xml:"compute_mode" json:"compute_mode"`
	FanSpeed                                               string   `xml:"fan_speed" json:"fan_speed"`
	Total                                                  string   `xml:"bar1_memory_usage>total" json:"total"`
	SmClock                                                string   `xml:"max_clocks>sm_clock" json:"sm_clock"`
	RxUtil                                                 string   `xml:"pci>rx_util" json:"rx_util"`
	GraphicsClock                                          string   `xml:"max_clocks>graphics_clock" json:"graphics_clock"`
	PwrObject                                              string   `xml:"inforom_version>pwr_object" json:"pwr_object"`
	PciBus                                                 string   `xml:"pci>pci_bus" json:"pci_bus"`
	DecoderUtil                                            string   `xml:"utilization>decoder_util" json:"decoder_util"`
	PciSubSystemID                                         string   `xml:"pci>pci_sub_system_id" json:"pci_sub_system_id"`
	MaxLinkGen                                             string   `xml:"pci>pci_gpu_link_info>pcie_gen>max_link_gen" json:"max_link_gen"`
	BridgeChipType                                         string   `xml:"pci>pci_bridge_chip>bridge_chip_type" json:"bridge_chip_type"`
	SmClockClocksGpu                                       string   `xml:"clocks>sm_clock" json:"sm_clock_clocks_gpu"`
	CurrentEcc                                             string   `xml:"ecc_mode>current_ecc" json:"current_ecc"`
	PowerDraw                                              string   `xml:"power_readings>power_draw" json:"power_draw"`
	CurrentLinkWidth                                       string   `xml:"pci>pci_gpu_link_info>link_widths>current_link_width" json:"current_link_width"`
	AutoBoost                                              string   `xml:"clock_policy>auto_boost" json:"auto_boost"`
	GpuUtil                                                string   `xml:"utilization>gpu_util" json:"gpu_util"`
	PciDevice                                              string   `xml:"pci>pci_device" json:"pci_device"`
	RegisterFile                                           string   `xml:"ecc_errors>volatile>single_bit>register_file" json:"register_file"`
	L2Cache                                                string   `xml:"ecc_errors>volatile>single_bit>l2_cache" json:"l_2_cache"`
	L1CacheDoubleBitAggregateEccErrorsGpu                  string   `xml:"ecc_errors>aggregate>double_bit>l1_cache" json:"l_1_cache_double_bit_aggregate_ecc_errors_gpu"`
	RetiredCount                                           string   `xml:"retired_pages>multiple_single_bit_retirement>retired_count" json:"retired_count"`
	PendingDm                                              string   `xml:"driver_model>pending_dm" json:"pending_dm"`
	AccountingModeBufferSize                               string   `xml:"accounting_mode_buffer_size" json:"accounting_mode_buffer_size"`
	GpuTempSlowThreshold                                   string   `xml:"temperature>gpu_temp_slow_threshold" json:"gpu_temp_slow_threshold"`
	OemObject                                              string   `xml:"inforom_version>oem_object" json:"oem_object"`
	TextureMemorySingleBitAggregateEccErrorsGpu            string   `xml:"ecc_errors>aggregate>single_bit>texture_memory" json:"texture_memory_single_bit_aggregate_ecc_errors_gpu"`
	RegisterFileSingleBitAggregateEccErrorsGpu             string   `xml:"ecc_errors>aggregate>single_bit>register_file" json:"register_file_single_bit_aggregate_ecc_errors_gpu"`
	MaxLinkWidth                                           string   `xml:"pci>pci_gpu_link_info>link_widths>max_link_width" json:"max_link_width"`
	TextureMemoryDoubleBitAggregateEccErrorsGpu            string   `xml:"ecc_errors>aggregate>double_bit>texture_memory" json:"texture_memory_double_bit_aggregate_ecc_errors_gpu"`
	ClocksThrottleReasonGpuIdle                            string   `xml:"clocks_throttle_reasons>clocks_throttle_reason_gpu_idle" json:"clocks_throttle_reason_gpu_idle"`
	MultigpuBoard                                          string   `xml:"multigpu_board" json:"multigpu_board"`
	GpuTempMaxThreshold                                    string   `xml:"temperature>gpu_temp_max_threshold" json:"gpu_temp_max_threshold"`
	MaxPowerLimit                                          string   `xml:"power_readings>max_power_limit" json:"max_power_limit"`
	L2CacheDoubleBitVolatileEccErrorsGpu                   string   `xml:"ecc_errors>volatile>double_bit>l2_cache" json:"l_2_cache_double_bit_volatile_ecc_errors_gpu"`
	PciDomain                                              string   `xml:"pci>pci_domain" json:"pci_domain"`
	MemClockDefaultApplicationsClocksGpu                   string   `xml:"default_applications_clocks>mem_clock" json:"mem_clock_default_applications_clocks_gpu"`
	VbiosVersion                                           string   `xml:"vbios_version" json:"vbios_version"`
	RetiredPageAddresses                                   string   `xml:"retired_pages>multiple_single_bit_retirement>retired_page_addresses" json:"retired_page_addresses"`
	GpuTemp                                                string   `xml:"temperature>gpu_temp" json:"gpu_temp"`
	AccountingMode                                         string   `xml:"accounting_mode" json:"accounting_mode"`
	L1CacheDoubleBitVolatileEccErrorsGpu                   string   `xml:"ecc_errors>volatile>double_bit>l1_cache" json:"l_1_cache_double_bit_volatile_ecc_errors_gpu"`
	DeviceMemorySingleBitAggregateEccErrorsGpu             string   `xml:"ecc_errors>aggregate>single_bit>device_memory" json:"device_memory_single_bit_aggregate_ecc_errors_gpu"`
	DisplayActive                                          string   `xml:"display_active" json:"display_active"`
	DefaultPowerLimit                                      string   `xml:"power_readings>default_power_limit" json:"default_power_limit"`
	EncoderUtil                                            string   `xml:"utilization>encoder_util" json:"encoder_util"`
	Serial                                                 string   `xml:"serial" json:"serial"`
	EnforcedPowerLimit                                     string   `xml:"power_readings>enforced_power_limit" json:"enforced_power_limit"`
	RetiredPageAddressesDoubleBitRetirementRetiredPagesGpu string   `xml:"retired_pages>double_bit_retirement>retired_page_addresses" json:"retired_page_addresses_double_bit_retirement_retired_pages_gpu"`
	EccObject                                              string   `xml:"inforom_version>ecc_object" json:"ecc_object"`
	Value                                                  []string `xml:"supported_clocks>supported_mem_clock>value" json:"value"`
	DisplayMode                                            string   `xml:"display_mode" json:"display_mode"`
	DeviceMemory                                           string   `xml:"ecc_errors>volatile>single_bit>device_memory" json:"device_memory"`
	PendingEcc                                             string   `xml:"ecc_mode>pending_ecc" json:"pending_ecc"`
	ClocksThrottleReasonSwPowerCap                         string   `xml:"clocks_throttle_reasons>clocks_throttle_reason_sw_power_cap" json:"clocks_throttle_reason_sw_power_cap"`
	TotalSingleBitAggregateEccErrorsGpu                    string   `xml:"ecc_errors>aggregate>single_bit>total" json:"total_single_bit_aggregate_ecc_errors_gpu"`
	CurrentGom                                             string   `xml:"gpu_operation_mode>current_gom" json:"current_gom"`
	MemoryUtil                                             string   `xml:"utilization>memory_util" json:"memory_util"`
}

type nvidiaSmi struct {
	Timestamp     string `xml:"timestamp" json:"timestamp"`
	DriverVersion string `xml:"driver_version" json:"driver_version"`
	AttachedGpus  string `xml:"attached_gpus" json:"attached_gpus"`
	GPUS          []gpu  `xml:"gpu" json:"gpus"`
}

func NewGpuInfoParserHelper(str string, err error) (*GpuInfo, error) {
	if err != nil {
		return nil, err
	}
	return NewGpuInfoParser(str)
}

func NewGpuInfoParser(str string) (*GpuInfo, error) {
	res := new(nvidiaSmi)
	err := xml.Unmarshal([]byte(str), res)
	if err != nil {
		return nil, err
	}
	gpuInfo := &GpuInfo{
		DriverVersion: res.DriverVersion,
	}
	for _, g := range res.GPUS {
		gpuInfo.Gpus = append(gpuInfo.Gpus, &PerGpuInfo{
			GpuName: g.ProductName,
			GpuMemInfo: GpuMemInfo{
				total: g.TotalFbMemoryUsageGpu,
				free:  g.FreeFbMemoryUsageGpu,
				used:  g.UsedFbMemoryUsageGpu,
			},
		})
	}
	return gpuInfo, nil
}

type GpuInfo struct {
	DriverVersion string
	Gpus          []*PerGpuInfo
}

func (g *GpuInfo) Info() string {
	if g.DriverVersion == "" {
		return ""
	}

	buf := &bytes.Buffer{}
	buf.Write([]byte(fmt.Sprintf("driver:%s ", g.DriverVersion)))
	for idx, gpu := range g.Gpus {
		buf.Write([]byte(fmt.Sprintf("%d:%s/%s ", idx, gpu.GpuMemInfo.Free(), gpu.GpuMemInfo.Total())))
	}
	return buf.String()
}

type PerGpuInfo struct {
	GpuName    string
	GpuMemInfo GpuMemInfo
}

type GpuMemInfo struct {
	total string
	free  string
	used  string
}

func (g *GpuMemInfo) Free() string {
	return g.free
}

func (g *GpuMemInfo) Total() string {
	return g.total
}
