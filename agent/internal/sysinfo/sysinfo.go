package sysinfo

import (
	"fmt"
	"net"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

// GetOSInfo 获取操作系统的简要信息
func GetOSInfo() (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	// 返回类似 "Ubuntu 22.04 (linux)" 的格式
	return fmt.Sprintf("%s %s (%s)", info.Platform, info.PlatformVersion, info.PlatformFamily), nil
}

// GetPrimaryIP 获取本机的主要对外 IP 地址
// 这段代码的逻辑是尝试连接到一个公网地址，然后查看本地使用了哪个 IP
func GetPrimaryIP() (string, error) {
	// 方法一：遍历所有网络接口
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				// 找到了第一个非环回的 IPv4 地址
				return ipnet.IP.String(), nil
			}
		}
	}

	// 方法二：如果上面方法找不到，尝试拨号一个公网地址（不会真正发送数据）
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		// 如果拨号失败，说明可能没有外网连接，返回错误
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip := localAddr.IP.String()

	// 有时会返回 IPv6 格式，我们需要清理一下
	if strings.Contains(ip, ":") {
		parts := strings.Split(ip, ":")
		return parts[len(parts)-1], nil
	}

	return ip, nil
}
