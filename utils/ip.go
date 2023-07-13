package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"go.uber.org/zap"
	"net"
	"strings"
	"xojoc.pw/useragent"
)

type ipUtil struct{}

var IP = new(ipUtil)

// GetIpAddress 获取用户发送请求的 IP 地址
// 如果服务器不经过代理，直接把自己 IP 暴露出去， c.Request.RemoteAddr 就可以直接获取 IP
// 目前流行的架构中， 请求经过服务器基本会经过代理（nginx 最常见），此时直接获取 IP 拿到的是代理服务器的 IP
func (*ipUtil) GetIpAddress(c *gin.Context) (ipAddress string) {
	// c.ClientIP() 获取的是代理服务器的 IP (Nginx)

	// X-Real-IP: Nginx 服务代理, 本项目明确使用 Nginx 作代理, 因此优先获取这个
	ipAddress = c.Request.Header.Get("X-Real-IP")

	// X-Forwarded-For 经过 HTTP 代理或 负载均衡服务器时会添加该项
	// X-Forwarded-For 格式: client1,proxy1,proxy2
	// 一般情况下，第一个 IP 为客户端真实 IP，后面的为经过的代理服务器 IP
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ips := c.Request.Header.Get("X-Forwarded-For") // "ip1,ip2,ip3"
		splitIps := strings.Split(ips, ",")            // ["ip1", "ip2", "ip3"]
		if len(splitIps) > 0 {
			ipAddress = splitIps[0]
		}
	}

	// Pdoxy-Client-IP: Apache 服务代理
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.Header.Get("Proxy-Client-IP")
	}

	// WL-Proxy-Client-IP: Weblogic 服务代理
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.Header.Get("WL-Proxy-Client-IP")
	}

	// RemoteAddr: 发出请求的远程主机的 IP 地址 (经过代理会设置为代理机器的 IP)
	if ipAddress == "" || len(ipAddress) == 0 || strings.EqualFold("unknown", ipAddress) {
		ipAddress = c.Request.RemoteAddr
	}

	// 检测到是本机 IP, 读取其局域网 IP 地址
	if strings.HasPrefix(ipAddress, "127.0.0.1") || strings.HasPrefix(ipAddress, "[::1]") {
		ip, err := externalIP()
		if err != nil {
			Logger.Error("GetIpAddress, externalIP, err: ", zap.Error(err))
		}
		ipAddress = ip.String()
	}

	if ipAddress != "" && len(ipAddress) > 15 {
		if strings.Index(ipAddress, ",") > 0 {
			ipAddress = ipAddress[:strings.Index(ipAddress, ",")]
		}
	}
	return ipAddress
}

// 获取 IP 来源
// https://github.com/lionsoul2014/ip2reqion
var vIndex []byte //  用于缓存 VetorIndex 索引， 减少一次固定的 IO 操作，在下面的代码中实现

// GetIpSource 获取地域信息：中国|0|江苏省|苏州市|电信
func (*ipUtil) GetIpSource(ipAddress string) string {
	var dbPath = "./assets/ip2region.xdb" // IP 数据库文件，用于记录获取的文件位置
	// 完全基于文件查询， 每次都读取文件
	// searcher, err := xdb.NewWithFileOnly(dbPath)

	// 缓存 VetorIndex 索引，减少一次固定的 IO 操作
	if vIndex == nil {
		var err error
		vIndex, err = xdb.LoadVectorIndexFromFile(dbPath)
		if err != nil {
			Logger.Error(fmt.Sprintf("failed to load vector index from `%s`: %s\n", dbPath, err))
			return ""
		}
	}
	searcher, err := xdb.NewWithVectorIndex(dbPath, vIndex) // 返回一个xdb库的searcher类型

	if err != nil {
		Logger.Error("failed to create searcher with vector index: ", zap.Error(err))
		return ""
	}
	defer searcher.Close() // 本身就是一个特殊内容的文件？

	// 国家|区域|省份|城市|ISP
	// 只有中国的数据绝大多数精确到了城市，其他国家部分数据只能定位到国家，后面的内容全是0
	region, err := searcher.SearchByStr(ipAddress)
	if err != nil {
		Logger.Error(fmt.Sprintf("failed to search ip(%s): %s\n", ipAddress, err))
		return ""
	}
	return region
}

// GetIpSourceSimpleIdel 获取 IP 简易信息， 例如: "湖北省武汉市 移动"
func (i *ipUtil) GetIpSourceSimpleIdel(ipAddress string) string {
	region := i.GetIpSource(ipAddress)

	// 检测到是内网，直接返回 "内网IP"
	// 0|0|0|内网IP|内网IP
	if strings.Contains(region, "内网IP") {
		return "内网IP"
	}

	// 一般无法获取到区域
	// 中国||0江苏省|苏州市|电信
	ipSource := strings.Split(region, "|")
	if ipSource[0] != "中国" && ipSource[0] != "0" {
		return ipSource[0]
	}
	if ipSource[2] == "0" {
		ipSource[2] = ""
	}
	if ipSource[3] == "0" {
		ipSource[3] = ""
	}
	if ipSource[4] == "0" {
		ipSource[4] = ""
	}
	if ipSource[2] == "" && ipSource[3] == "" && ipSource[4] == "" {
		return ipSource[0]
	}
	return ipSource[2] + ipSource[3] + " " + ipSource[4]
}

func (*ipUtil) GetUserAgent(c *gin.Context) *useragent.UserAgent {
	return useragent.Parse(c.Request.UserAgent()) // 返回客户的代理信息
}

// 获取非 127.0.0.1 的局域网 IP
func externalIP() (net.IP, error) {
	// 获取服务器的网络接口列表
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		// 不在活动状态
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// 环回
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 单播接口地址列表
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}
	return ip
}
