/*
 * @Author: duanzt
 * @Date: 2023-07-14 16:38:35
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-14 16:54:52
 * @FilePath: iptools.go
 * @Description: ip处理的工具
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package tools

import "net"

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
