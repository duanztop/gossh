/*
 * @Author: duanzt
 * @Date: 2023-07-14 09:41:38
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-14 10:26:24
 * @FilePath: iconnection.go
 * @Description: 定义connection interface
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package internal

import (
	"context"
	"io"
)

// IConnection connection interface
type IConnection interface {

	// Close 关闭连接
	//  @author duanzt
	//  @date 2023-07-14 09:48:57
	//  @return error 关闭异常时返回
	Close() error

	// Exec 执行(自定义session动作)
	//  @author duanzt
	//  @date 2023-07-14 09:53:16
	//  @param context.Context 上下文context
	//  @param func(isession) error 从该function中获取session进行处理
	//  @return string 执行输出
	//  @return error ssh异常时返回
	Exec(context.Context, func(ISession) error) (string, error)

	// ExecShell 执行shell
	//  @author duanzt
	//  @date 2023-07-14 09:54:35
	//  @param context.Context 上下文context
	//  @param string shell脚本
	//  @return string 执行shell输出结果
	//  @return error ssh异常时返回
	ExecShell(context.Context, string) (string, error)

	// CopyFileLTR 拷贝文件流到远端
	//  @author duanzt
	//  @date 2023-07-14 09:56:42
	//  @param src io.Reader 流
	//  @param dest string 远端目标文件地址
	//  @param mode string 文件权限
	//  @return error ssh异常时返回
	CopyFileITR(src io.Reader, dest, mode string) error

	// CopyFileITRMon 拷贝文件流到远端（监控远端目标文件大小）
	//  @author duanzt
	//  @date 2023-07-14 10:02:16
	//  @param src io.Reader 流
	//  @param dest string 远端目标文件地址
	//  @param mode string 文件权限
	//  @param destSizeChan chan int64 返回远端目标文件大小，单位：byte
	//  @return error ssh异常时返回
	CopyFileITRMon(src io.Reader, dest, mode string, destSizeChan chan int64) error

	// CopyFileLTR 拷贝本地文件到远端
	//  @author duanzt
	//  @date 2023-07-14 10:00:05
	//  @param  src dest 本地文件地址
	//  @param dest string 远端目标文件地址
	//  @param mode string 文件权限
	//  @return error ssh异常时返回
	CopyFileLTR(src, dest, mode string) error

	// CopyFileLTRMon 拷贝本地文件到远端（监控远端目标文件大小）
	//  @author duanzt
	//  @date 2023-07-14 10:00:05
	//  @param src string 本地文件地址
	//  @param dest string 远端目标文件地址
	//  @param mode string 文件权限
	//  @param destSizeChan chan int64 返回远端目标文件大小，单位：byte
	//  @return error ssh异常时返回
	CopyFileLTRMon(src, dest, mode string, destSizeChan chan int64) error

	// CopyFileRTL 拷贝远端文件到本地
	//  @author duanzt
	//  @date 2023-07-14 09:59:07
	//  @param src string 远端文件地址
	//  @param dest string 本地目标文件地址
	//  @param mode string 文件权限
	//  @return error ssh异常时返回
	CopyFileRTL(src string, dest, mode string) error

	// CopyFileRTLMon 拷贝远端文件到本地（监控本地目标文件大小）
	//  @author duanzt
	//  @date 2023-07-14 09:59:07
	//  @param src string 远端文件地址
	//  @param dest string 本地目标文件地址
	//  @param mode string 文件权限
	//  @param destSizeChan chan int64 返回本地目标文件大小，单位：byte
	//  @return error ssh异常时返回
	CopyFileRTLMon(src string, dest, mode string, destSizeChan chan int64) error

	// GetAddr 获取ssh连接地址（例127.0.0.1:22）
	//  @author duanzt
	//  @date 2023-07-14 10:06:15
	//  @return string ssh连接地址
	GetAddr() string

	// GetIp 获取ssh ip（例127.0.0.1）
	//  @author duanzt
	//  @date 2023-07-14 10:06:36
	//  @return string ip地址
	GetIp() string
}
