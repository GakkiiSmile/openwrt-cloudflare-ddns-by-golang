package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"main/model"
	"net"
	"os/exec"
	"strings"
)

// Shellout 执行shell命令
func Shellout(command string, arg ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(command, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// Ipv6NeighborInfo 网上邻居ipv6地址结构体
type Ipv6NeighborInfo struct {
	Addr   string
	Mac    string
	Status string
}

// Ipv6Neighbor里status的枚举类型
const (
	STALE     = "STALE"
	REACHABLE = "REACHABLE"
	DELAY     = "DELAY"
	FAILED    = "FAILED"
)

// GetIpv6NeighborList Openwrt 获取ipv6的网络邻居地址信息 一般邻居都是分配ipv6 不考虑ipv4情况
// 一般解析格式为 2409:8a55:xxxx:xxxx::1 dev br-lan lladdr 34:08:bc:xx:xx:xx STALE
func GetIpv6NeighborList() ([]Ipv6NeighborInfo, error) {
	outString, _, err := Shellout("ip", "-6", "neighbor", "show")

	if err != nil {
		return nil, errors.New("当前命令可能不支持当前系统或者系统未安装ip软件")
	}

	lines := strings.Split(outString, "\n")

	ipv6Infos := make([]Ipv6NeighborInfo, 0, len(lines))

	for _, line := range lines {
		if strings.Trim(line, " ") != "" {
			ipInfo := strings.Split(line, " ")

			addr := ipInfo[0]
			mac := ""
			status := ipInfo[len(ipInfo)-1]

			// 如果状态是错误，则丢弃
			if status != FAILED {
				mac = ipInfo[len(ipInfo)-2]
			}

			ipv6Infos = append(ipv6Infos, Ipv6NeighborInfo{
				Addr:   addr,
				Mac:    mac,
				Status: status,
			})
		}
	}

	return ipv6Infos, nil
}

// 简易判断ip是否为一个私有地址
func isLocal(ip net.IP) bool {

	if ip4 := ip.To4(); ip4 != nil {

		return (ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1]&0xf0 == 16) ||
			(ip4[0] == 192 && ip4[1] == 168))
	}
	return len(ip) == 16 && ip[0]&0xfe == 0xfc
}

// IPs 存储ip字段信息
type IPs struct {
	IPv4 []string

	IPv6 []string
}

// GetCurrentIPs 获取本机网卡IP
func GetCurrentIPs() (ips IPs, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIPNet bool
	)

	ips = IPs{}

	// 获取所有网卡
	addrs, err = net.InterfaceAddrs()

	if err != nil {
		return ips, err
	}

	// 取网卡IP
	for _, addr = range addrs {
		if ipNet, isIPNet = addr.(*net.IPNet); isIPNet && ipNet.IP.IsGlobalUnicast() && !isLocal(ipNet.IP) {
			ipString := ipNet.IP.String()

			if len(ipNet.IP) == 4 {
				ips.IPv4 = append(ips.IPv4, ipString)
			} else if len(ipNet.IP) == 16 {
				ips.IPv6 = append(ips.IPv6, ipString)
			}
		}
	}

	if len(ips.IPv4) == 0 && len(ips.IPv6) == 0 {
		err = errors.New("没有获取到一个公网ip")
	}

	return ips, err
}

// LoadConfig 从本地加载config
func LoadConfig(filePath string) (config *model.Config, err error) {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		errString := fmt.Sprintf("文件打开失败，原因：%v", err)

		return config, errors.New(errString)
	}

	config = &model.Config{}

	err = json.Unmarshal(file, config)

	if err != nil {
		errString := fmt.Sprintf("JSON解析失败，原因：%v", err)

		return config, errors.New(errString)
	}

	return config, nil
}
