/*
 * @Author: duanzt
 * @Date: 2023-07-14 16:38:35
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-17 18:48:08
 * @FilePath: iptools.go
 * @Description: ip处理的工具
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package tools

import (
	"net"
	"regexp"
)

const (

	// ipv4Regex ipv4地址校验正则
	ipv4Regex = `^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?).){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

	// ipv6Regex ipv6地址校验正则
	ipv6Regex = `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:([0-9a-fA-F]{1,4}|:)|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]).){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]).){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`

	localhost   = "localhost"
	localhostIp = "127.0.0.1"
)

type iptools struct{}

var (
	IpTools = iptools{}
)

// GetFlagUpV4Interface 获取本机启用的ipv4网卡信息
//
//	@author duanzt
//	@date 2023-07-14 04:48:36
//	@receiver iptools
//	@return []net.Interface 机启用的ipv4网卡信息
//	@return error 异常时返回
func (iptools) GetFlagUpV4Interface() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	result := make([]net.Interface, 0)
	for i := range ifaces {
		if ifaces[i].Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if ifaces[i].Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		result = append(result, ifaces[i])
	}
	return result, nil
}

// V4Address 获取网卡ipv4地址
//
//	@author duanzt
//	@date 2023-07-14 04:52:15
//	@receiver iptools
//	@param addr net.Addr 地址信息
//	@return net.IP ipv4
func (iptools) V4Address(addr net.Addr) net.IP {
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
	return ip.To4()
}

// IsV4 校验ip是否是ipv4
//
//	@author duanzt
//	@date 2023-07-17 06:37:12
//	@receiver iptools
//	@param ip string ip字符串
//	@return bool 是ipv4则返回true
func (iptools) IsV4(ip string) bool {
	match, err := regexp.MatchString(ipv4Regex, ip)
	return err == nil && match
}

// IsV6 校验ip是否是ipv6
//
//	@author duanzt
//	@date 2023-07-17 06:38:28
//	@receiver iptools
//	@param ip string ip字符串
//	@return bool 是ipv6则返回true
func (iptools) IsV6(ip string) bool {
	match, err := regexp.MatchString(ipv6Regex, ip)
	return err == nil && match
}

// CheckIpIsLocal 判断ip是否为本机
//
//	@author duanzt
//	@date 2023-07-17 06:42:42
//	@receiver iptools
//	@param ip string ip字符串
//	@return bool ip为本机则返回true
func (iptools iptools) CheckIpIsLocal(ip string) bool {
	if ip == localhost || ip == localhostIp {
		return true
	}
	interfaces, err := iptools.GetFlagUpV4Interface()
	if err != nil {
		return false
	}
	for i := range interfaces {
		if netAddr, err := interfaces[i].Addrs(); err != nil {
			for j := range netAddr {
				if netIp := iptools.V4Address(netAddr[j]); netIp != nil && netIp.String() == ip {
					return true
				}
			}
		}
	}
	return false
}
