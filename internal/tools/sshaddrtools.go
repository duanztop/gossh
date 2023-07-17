/*
 * @Author: duanzt
 * @Date: 2023-07-17 18:27:15
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-17 18:52:24
 * @FilePath: sshaddrtools.go
 * @Description: ssh addr工具
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package tools

import (
	"errors"
	"strings"
)

const (

	// addrSplit ssh地址分隔符
	addrSplit = ":"

	// sshDefaultPort ssh默认连接端口
	sshDefaultPort = "22"
)

type sshAddrtools struct{}

var (
	SshAddrTools = sshAddrtools{}
)

// SetRightAddr 设置正确的ssh连接地址
//
//	@author duanzt
//	@date 2023-07-17 06:30:11
//	@receiver s sshAddrtools
//	@param addr string 传入的ssh连接地址
//	@return string 正确的ssh连接地址
//	@return error ssh连接地址异常时返回
func (s sshAddrtools) SetRightAddr(addr string) (string, error) {
	if addr == "" {
		return "", errors.New("连接地址不可为空")
	}
	ss := strings.Split(addr, addrSplit)
	// 判断ip是否正确
	if ss[0] != localhost {
		isv4 := IpTools.IsV4(ss[0])
		if !isv4 {
			return "", errors.New("ip不合法,请使用ipv4地址")
		}
	}
	if len(ss) == 1 {
		return addr + addrSplit + sshDefaultPort, nil
	}
	return addr, nil
}

// GetAddrSplit 获取ssh地址连接分隔符
//
//	@author duanzt
//	@date 2023-07-17 06:51:28
//	@receiver sshAddrtools
//	@return string  ssh地址连接分隔符
func (sshAddrtools) GetAddrSplit() string {
	return addrSplit
}
