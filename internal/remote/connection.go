/*
 * @Author: duanzt
 * @Date: 2023-07-14 10:27:51
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-17 13:12:18
 * @FilePath: connection.go
 * @Description: 远程ssh连接
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package remote

import (
	"context"
	"embed"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/duanztop/gossh/internal"
	"github.com/duanztop/gossh/internal/tools"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (

	// defaultUsername 默认用户名
	defaultUsername = "root"

	// defaultPrivateKey 默认私钥存放路径
	defaultPrivateKey = ".ssh/id_rsa"
)

var (
	//go:embed .ssh
	multi embed.FS
)

type connection struct {
	client *ssh.Client
	addr   string // 地址信息
}

// Close 关闭连接
// @author duanzt
// @date 2023-07-14 09:48:57
// @return error 关闭异常时返回
func (c *connection) Close() error {
	return c.client.Close()
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
	sess, err := c.generateSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()

	if err := fn(sess); err != nil {
		return "", err
	}

	if err := sess.Wait(); err != nil {
		return "", err
	}
	return sess.Output(), err
}

// ExecShell 执行shell
//
//	@author duanzt
//	@date 2023-07-14 09:54:35
//	@param ctx context.Context 上下文context
//	@param shell string shell脚本
//	@return string 执行shell输出结果
//	@return error ssh异常时返回
func (c *connection) ExecShell(ctx context.Context, shell string) (string, error) {
	return c.Exec(ctx, func(i internal.ISession) error {
		return i.Exec(shell)
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
	sftpClient, err := sftp.NewClient(c.client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	err = sftpClient.MkdirAll(filepath.Dir(dest))
	if err != nil {
		return err
	}
	fd, err := sftpClient.Create(dest)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = io.Copy(fd, src)
	// buf := make([]byte, 1024)
	// for {
	// 	n, _ := src.Read(buf)
	// 	if n == 0 {
	// 		break
	// 	}
	// 	fd.Write(buf[:n])
	// }
	return err
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
func (c *connection) CopyFileITRMon(src io.Reader, dest string, mode string, destSizeChan chan int64) (err error) {
	sftpClient, err := sftp.NewClient(c.client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	if err := sftpClient.MkdirAll(filepath.Dir(dest)); err != nil {
		return err
	}
	fd, err := sftpClient.Create(dest)
	if err != nil {
		return err
	}
	defer fd.Close()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		_, err = io.Copy(fd, src)
		waitGroup.Done()
	}()

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	go func() {
		defer func() {
			_ = recover()
		}()
		for range ticker.C {
			stat, _ := fd.Stat()
			destSizeChan <- stat.Size()
		}
	}()
	waitGroup.Wait()
	stat, _ := fd.Stat()
	destSizeChan <- stat.Size()
	close(destSizeChan)
	return err
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
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	return c.CopyFileITR(file, dest, mode)
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
func (c *connection) CopyFileLTRMon(src string, dest string, mode string, destSizeChan chan int64) (err error) {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	return c.CopyFileITRMon(file, dest, mode, destSizeChan)
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
	sftpClient, err := sftp.NewClient(c.client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	file, err := sftpClient.Open(src)
	if err != nil {
		return err
	}
	destFile, err := tools.FileTools.CreateFile(dest)
	if err != nil {
		return err
	}

	_, err = io.Copy(destFile, file)
	return err
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
func (c *connection) CopyFileRTLMon(src string, dest string, mode string, destSizeChan chan int64) (err error) {
	sftpClient, err := sftp.NewClient(c.client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()
	file, err := sftpClient.Open(src)
	if err != nil {
		return err
	}

	destFile, err := tools.FileTools.CreateFile(dest)
	if err != nil {
		return err
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		_, err = io.Copy(destFile, file)
		waitGroup.Done()
	}()

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	go func() {
		defer func() {
			_ = recover()
		}()
		for range ticker.C {
			stat, _ := destFile.Stat()
			destSizeChan <- stat.Size()
		}
	}()
	waitGroup.Wait()
	stat, _ := destFile.Stat()
	destSizeChan <- stat.Size()
	close(destSizeChan)
	return err

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

// generateSession 生成session对象
//
//	@author duanzt
//	@date 2023-07-14 11:23:02
//	@receiver c *connection
//	@return *session
//	@return error
func (c *connection) generateSession() (*session, error) {
	sshSess, err := c.client.NewSession()
	if err != nil {
		c.client.Close()
		return nil, err
	}
	return &session{sshSess: sshSess}, nil
}

// NewConnection1 新建连接（通过username+password方式）
//
//	@author duanzt
//	@date 2023-07-14 05:15:59
//	@param username string 用户名
//	@param password string 密码
//	@param addr string ssh连接地址
//	@return internal.IConnection ssh连接对象
//	@return error 连接异常时返回
func NewConnection1(username, password, addr string) (internal.IConnection, error) {
	auth := []ssh.AuthMethod{}
	keyboardInteractiveChallenge := func(
		username,
		instruction string,
		questions []string,
		echos []bool,
	) (answers []string, err error) {
		if len(questions) == 0 {
			return []string{}, nil
		}
		return []string{password}, nil
	}
	auth = append(auth, ssh.Password(password))
	auth = append(auth, ssh.KeyboardInteractive(keyboardInteractiveChallenge))
	return newConnectionBasic(auth, username, addr)
}

// NewConnection2 新建连接（通过username+私钥方式）
//
//	@author duanzt
//	@date 2023-07-14 05:16:52
//	@param username string 用户名
//	@param privateKey string 私钥文件地址
//	@param addr string ssh连接地址
//	@return internal.IConnection ssh连接
//	@return error 连接异常时返回
func NewConnection2(username, privateKey, addr string) (internal.IConnection, error) {
	auth := []ssh.AuthMethod{}
	pkData, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}
	pk, err := ssh.ParsePrivateKey(pkData)
	if err != nil {
		return nil, err
	}
	auth = append(auth, ssh.PublicKeys(pk))
	return newConnectionBasic(auth, username, addr)
}

// NewConnectionDefault 使用默认方式新建连接（默认用户名：root，默认使用私钥连接，私钥地址为当前目录下的.ssh/id_rsa文件）
//
//	@author duanzt
//	@date 2023-07-14 06:16:26
//	@param addr string ssh连接地址
//	@return internal.IConnection ssh连接
//	@return error 连接异常时返回
func NewConnectionDefault(addr string) (internal.IConnection, error) {
	auth := []ssh.AuthMethod{}
	f, _ := multi.Open(defaultPrivateKey)
	defer f.Close()
	data, err := tools.FileTools.ReadFile(f)
	if err != nil {
		return nil, err
	}
	pk, err := ssh.ParsePrivateKey(data)
	if err != nil {
		return nil, err
	}
	auth = append(auth, ssh.PublicKeys(pk))
	return newConnectionBasic(auth, defaultUsername, addr)
}

// newConnectionBasic 新建连接（默认方法，auth需要前置组装）
//
//	@author duanzt
//	@date 2023-07-14 05:13:11
//	@param auth []ssh.AuthMethod auth方法
//	@param username string 用户名
//	@param addr string ssh连接地址
//	@return internal.IConnection ssh连接
//	@return error 连接异常时返回
func newConnectionBasic(auth []ssh.AuthMethod, username, addr string) (internal.IConnection, error) {
	config := ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}

	clientConfig := &ssh.ClientConfig{
		User:    username,
		Auth:    auth,
		Timeout: time.Duration(1) * time.Minute,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err != nil {
		return nil, err
	}
	return &connection{client: client, addr: addr}, nil
}
