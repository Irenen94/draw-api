package utils

import (
	"aiot-service-for-mfp/pkg/logger"
	"bytes"
	_net "net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/docker"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func GetOutboundIP() string {
	conn, err := _net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.Logger.Error(err.Error())
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*_net.UDPAddr)
	return localAddr.IP.String()
}

func GetLocalIP() (ip string, err error) {
	addrs, err := _net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*_net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() {
			continue
		}
		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		return ipAddr.IP.String(), nil
	}
	return
}

/*
func getNetInfo() {
	info, _ := net.IOCounters(true)
	for index, v := range info {
		fmt.Printf("%v:%v send:%v recv:%v\n", index, v, v.BytesSent, v.BytesRecv)
	}
}
*/

// ok , cnt
func GetCpus() (bool, int) {
	t := true
	cnt, err := cpu.Counts(t)
	if err == nil {
		return true, cnt
	} else {
		logger.Logger.Error(err.Error())
		return false, 0
	}
}

// ok , pct
func GetCpuUsage() (bool, float64) {
	percpu := false
	pct, err := cpu.Percent(5000*time.Millisecond, percpu)
	if err == nil {
		return true, pct[0]
	} else {
		logger.Logger.Error(err.Error())
		return false, 0.0
	}
}

// ok , total G , pct
func GetMemInfo() (bool, float64, float64) {
	v, err := mem.VirtualMemory()
	if err == nil {
		return true, float64(v.Total) / 1024.0 / 1024.0 / 1024.0, v.UsedPercent
	} else {
		logger.Logger.Error(err.Error())
		return false, 0.0, 0.0
	}
}

/*
type InfoStat struct {
	Hostname             string `json:"hostname"`
	Uptime               uint64 `json:"uptime"`
	BootTime             uint64 `json:"bootTime"`
	Procs                uint64 `json:"procs"`           // number of processes
	OS                   string `json:"os"`              // ex: freebsd, linux
	Platform             string `json:"platform"`        // ex: ubuntu, linuxmint
	PlatformFamily       string `json:"platformFamily"`  // ex: debian, rhel
	PlatformVersion      string `json:"platformVersion"` // version of the complete OS
	KernelVersion        string `json:"kernelVersion"`   // version of the OS kernel (if available)
	KernelArch           string `json:"kernelArch"`      // native cpu architecture queried at runtime, as returned by `uname -m` or empty string in case of error
	VirtualizationSystem string `json:"virtualizationSystem"`
	VirtualizationRole   string `json:"virtualizationRole"` // guest or host
	HostID               string `json:"hostid"`             // ex: uuid
}
*/
// ok , hostname
func GetHostName() (bool, string) {
	hInfo, err := host.Info()
	if err == nil {
		return true, hInfo.Hostname
	} else {
		logger.Logger.Error(err.Error())
		return false, ""
	}
}

// ok , docker info
/*
[{"containerID":"93b621d95cc7c887056add519398da56a17e337eb2dfdd7d4542fcfb53d3053e","name":"kafka_backend","image":"cv_dev_basic:6","status_code":"Up 6 months","running":true}
 {"containerID":"8de351cfdc0b6a5d79ac06641275336b0f9e4dcd9dd0938ff6ba9d14cf0bdca7","name":"metabase","image":"metabase","status_code":"Up 6 months","running":true}]
type CgroupDockerStat struct {
	ContainerID string `json:"containerID"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Status      string `json:"status_code"`
	Running     bool   `json:"running"`
}
*/
func GetDockerTasks() (bool, []docker.CgroupDockerStat) {
	stat, err := docker.GetDockerStat()
	if err == nil {
		/*
			task_list := []map[string]string{}
			for _, item := range stat {
				task_list = append(task_list, item["name"])
			}
			return true, task_list
		*/
		return true, stat
	} else {
		logger.Logger.Error(err.Error())
		return false, nil
	}
}

func Shellout(command string) (error, string, string) {
	ShellToUse := "bash"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

// ok , driver , gpu_name , count , mem_total , mem_used , ut_gpu , ut_mem
func GetGpuInfo() (bool, string, string, int, float64, float64, float64, float64) {
	err, out, errout := Shellout("nvidia-smi --query-gpu=driver_version,gpu_name,memory.total,memory.used,utilization.memory,utilization.gpu,count --format=csv,noheader,nounits")
	if err != nil {
		logger.Logger.Error(err.Error())
		logger.Logger.Error(errout)
		return false, "", "", 0, 0.0, 0.0, 0.0, 0.0
	}
	// parse
	fields := strings.Split(out, ",")
	if len(fields) != 7 {
		return false, "", "", 0, 0.0, 0.0, 0.0, 0.0
	}
	count, err1 := strconv.Atoi(strings.Trim(strings.Trim(fields[6], " "), "\n"))
	if err1 != nil {
		logger.Logger.Error(err1.Error())
	}
	mem_total, err2 := strconv.ParseFloat(strings.Trim(fields[2], " "), 64)
	if err2 != nil {
		logger.Logger.Error(err2.Error())
	}
	mem_used, err3 := strconv.ParseFloat(strings.Trim(fields[3], " "), 64)
	if err3 != nil {
		logger.Logger.Error(err3.Error())
	}
	ut_gpu, err4 := strconv.ParseFloat(strings.Trim(fields[5], " "), 64)
	if err4 != nil {
		logger.Logger.Error(err4.Error())
	}
	ut_mem, err5 := strconv.ParseFloat(strings.Trim(fields[4], " "), 64)
	if err5 != nil {
		logger.Logger.Error(err5.Error())
	}
	return true, fields[0], fields[1], count, mem_total, mem_used, ut_gpu, ut_mem
}

/*
func main() {
	fmt.Println(GetLocalIP())
	fmt.Println(GetOutboundIP())
	fmt.Println(GetCpus())
	fmt.Println(GetCpuUsage())
	fmt.Println(GetMemInfo())
	fmt.Println(GetHostName())
	fmt.Println(GetDockerTasks())
	fmt.Println(GetGpuInfo())
}
*/
