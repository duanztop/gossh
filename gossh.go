/*
 * @Author: duanzt
 * @Date: 2023-07-14 10:26:52
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-14 18:25:48
 * @FilePath: gossh.go
 * @Description: 暴露文件，提供使用的方法
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package gossh

import (
	"github.com/duanztop/gossh/internal"
	"github.com/duanztop/gossh/internal/local"
	"github.com/duanztop/gossh/internal/remote"
)

// Remote1 获取远程ssh连接（使用username+password验证方式）
//
//	@author duanzt
//	@date 2023-07-14 04:59:22
//	@param username string 用户名
//	@param password string 密码
//	@param addr string ssh连接地址，例如：192.168.10.100:22（只传入ip的情况会默认使用22端口）
//	@return internal.IConnection ssh连接
//	@return error 连接异常时返回
func Remote1(username, password, addr string) (internal.IConnection, error) {
	// TODO:判断addr，如果是ip增加默认后缀:22
	// TODO:判断ip，如果是本机ip（或127.0.0.1，或localhost），则直接使用本地ssh连接，降低远程ssh连接损耗

	return remote.NewConnection1(username, password, addr)
}

// Remote2 获取远程ssh连接（使用username+私钥验证方式）
//
//	@author duanzt
//	@date 2023-07-14 05:01:36
//	@param user string 用户名
//	@param privateKey string 私钥地址
//	@param addr string  ssh连接地址，例如：192.168.10.100:22（只传入ip的情况会默认使用22端口）
//	@return internal.IConnection ssh连接
//	@return error 连接异常时返回
func Remote2(username, privateKey, addr string) (internal.IConnection, error) {
	// TODO:判断addr，如果是ip增加默认后缀:22
	// TODO:判断ip，如果是本机ip（或127.0.0.1，或localhost），则直接使用本地ssh连接，降低远程ssh连接损耗
	return remote.NewConnection2(username, privateKey, addr)
}

// RemoteDefault 获取远程ssh连接（使用username+默认私钥验证方式）
// 注意：请添加以下公钥到目标机器的/root/.ssh/authorized_keys文件中
// ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDP3qG9zzNQkhKuWWJTCHAuY04bQ9h/vhZplVrTSnEoWL1SsbT7v/dCDXuyazNDo9ikd9BS6H/nE5lOKp+Omi1W2uWs/kxrCNhouXRO8kMLViRB3DRP2VYFDo36UafzNdGKkH/vW4Ptilga/ucForW05SindT3KeKf+tB1u3RBlRz6rzpeuqrflVtcaWtQ33exWMO8CxgzCtsDexWWLP+TLdeOaWyfn0hj4tf36+K7oENAzGGhQuEwETiMUkJKfykBThBenWgU9mM1/5VbgvGiW7xIoeyDX8RI6Lz5q8mb3+ajuEqPyX/qwiNasYkQ7bWGaLDAVF3yJ20w7EpP54yi9rEoiBt6GAEo2JX5OuibzwMsz2CCykiB8H4YyiOlBY0q5GwrXC87fslvEt4KcYdk/XrZT+ikrJePgTCbQJhUGf8yYe1aDKoTBf/uIuT/O7aJ49KxnfeQdxel3xIykyROxnNisQ8Iz3vdC/QZlZsQUnJzo0UXtwDpwmwLwjpLFCMM= gossh@github.com
//
//	@author duanzt
//	@date 2023-07-14 06:24:43
//	@param addr string
//	@return internal.IConnection
//	@return error
func RemoteDefault(addr string) (internal.IConnection, error) {
	// TODO:判断addr，如果是ip增加默认后缀:22
	// TODO:判断ip，如果是本机ip（或127.0.0.1，或localhost），则直接使用本地ssh连接，降低远程ssh连接损耗
	return remote.NewConnectionDefault(addr)
}

// Local 获取本地ssh连接（）
//
//	@author duanzt
//	@date 2023-07-14 05:05:44
//	@return internal.IConnection
func Local() internal.IConnection {
	return local.NewConnection()
}
