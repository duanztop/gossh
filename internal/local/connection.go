/*
 * @Author: duanzt
 * @Date: 2023-07-14 10:27:45
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-14 17:56:13
 * @FilePath: connection.go
 * @Description: 本地连接（逻辑上，并没有建立任何连接）
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package local

import (
	"context"
	"io"
	"strings"

	"github.com/duanztop/gossh/internal"
)

const (

	// defaultAddress 默认本地ssh连接地址
	defaultAddress = "127.0.0.1:22"
)

// connection 本地连接（逻辑上，并没有建立任何连接）
type connection struct {
	addr string // 地址信息
}

// Close 关闭连接
// @author duanzt
// @date 2023-07-14 09:48:57
// @return error 关闭异常时返回
func (c *connection) Close() error {
	return nil
}

// Exec 执行(自定义session动作)
//
//	@author duanzt
//	@date 2023-07-14 09:53:16
//	@param ctx context.Context 上下文context
//	@param fn func(isession) error 从该function中获取session进行处理
//	@return string 执行输出
//	@return error ssh异常时返回
func (c *connection) Exec(ctx context.Context, fn func(internal.ISession) error) (string, error) {
	sess := &session{ctx: ctx}
	defer sess.Close()

	if err := fn(sess); err != nil {
		return sess.Output(), err
	}

	err := sess.Wait()
	return sess.Output(), err
}

// ExecShell 执行shell
//
//	@author duanzt
//	@date 2023-07-14 09:54:35
//	@param cxt context.Context 上下文context
//	@param shell string shell脚本
//	@return string 执行shell输出结果
//	@return error ssh异常时返回
func (c *connection) ExecShell(cxt context.Context, shell string) (string, error) {
	return c.Exec(cxt, func(s internal.ISession) error {
		return s.Exec(shell)
	})
}

// CopyFileLTR 拷贝文件流到远端
//
//	@author duanzt
//	@date 2023-07-14 09:56:42
//	@param src io.Reader 流
//	@param dest string 远端目标文件地址
//	@param mode string 文件权限
//	@return error ssh异常时返回
func (c *connection) CopyFileITR(src io.Reader, dest string, mode string) error {
	panic("not implemented") // TODO: Implement
}

// CopyFileITRMon 拷贝文件流到远端（监控远端目标文件大小）
//
//	@author duanzt
//	@date 2023-07-14 10:02:16
//	@param src io.Reader 流
//	@param dest string 远端目标文件地址
//	@param mode string 文件权限
//	@param destSizeChan chan int64 返回远端目标文件大小，单位：byte
//	@return error ssh异常时返回
func (c *connection) CopyFileITRMon(src io.Reader, dest string, mode string, destSizeChan chan int64) error {
	panic("not implemented") // TODO: Implement
}

// CopyFileLTR 拷贝本地文件到远端
//
//	@author duanzt
//	@date 2023-07-14 10:00:05
//	@param  src dest 本地文件地址
//	@param dest string 远端目标文件地址
//	@param mode string 文件权限
//	@return error ssh异常时返回
func (c *connection) CopyFileLTR(src string, dest string, mode string) error {
	panic("not implemented") // TODO: Implement
}

// CopyFileLTRMon 拷贝本地文件到远端（监控远端目标文件大小）
//
//	@author duanzt
//	@date 2023-07-14 10:00:05
//	@param src string 本地文件地址
//	@param dest string 远端目标文件地址
//	@param mode string 文件权限
//	@param destSizeChan chan int64 返回远端目标文件大小，单位：byte
//	@return error ssh异常时返回
func (c *connection) CopyFileLTRMon(src string, dest string, mode string, destSizeChan chan int64) error {
	panic("not implemented") // TODO: Implement
}

// CopyFileRTL 拷贝远端文件到本地
//
//	@author duanzt
//	@date 2023-07-14 09:59:07
//	@param src string 远端文件地址
//	@param dest string 本地目标文件地址
//	@param mode string 文件权限
//	@return error ssh异常时返回
func (c *connection) CopyFileRTL(src string, dest string, mode string) error {
	panic("not implemented") // TODO: Implement
}

// CopyFileRTLMon 拷贝远端文件到本地（监控本地目标文件大小）
//
//	@author duanzt
//	@date 2023-07-14 09:59:07
//	@param src string 远端文件地址
//	@param dest string 本地目标文件地址
//	@param mode string 文件权限
//	@param destSizeChan chan int64 返回本地目标文件大小，单位：byte
//	@return error ssh异常时返回
func (c *connection) CopyFileRTLMon(src string, dest string, mode string, destSizeChan chan int64) error {
	panic("not implemented") // TODO: Implement
}

// GetAddr 获取ssh连接地址（例127.0.0.1:22）
//
//	@author duanzt
//	@date 2023-07-14 10:06:15
//	@return string ssh连接地址
func (c *connection) GetAddr() string {
	return c.addr
}

// GetIp 获取ssh ip（例127.0.0.1）
//
//	@author duanzt
//	@date 2023-07-14 10:06:36
//	@return string ip地址
func (c *connection) GetIp() string {
	return strings.Split(c.addr, ":")[0]
}

// NewConnection 新建一个本地ssh连接对象
//
//	@author duanzt
//	@date 2023-07-14 05:06:51
//	@return internal.IConnection
func NewConnection() internal.IConnection {
	return &connection{addr: defaultAddress}
}
