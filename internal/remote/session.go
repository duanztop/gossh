/*
 * @Author: duanzt
 * @Date: 2023-07-14 10:27:38
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-14 11:28:56
 * @FilePath: session.go
 * @Description: 远程ssh session管理
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package remote

import (
	"bufio"
	"bytes"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

// session 远程ssh session
type session struct {
	sshSess *ssh.Session
	output  *bytes.Buffer
}

// Exec 执行shell
// @author duanzt
// @date 2023-07-14 10:09:53
// @param shell string shell命令
// @return error 执行异常时返回
func (s *session) Exec(shell string) error {
	return s.ExecOutput(shell, nil)
}

// ExecOutput 执行并同步获取输出结果
//
//	@author duanzt
//	@date 2023-07-14 10:12:33
//	@param shell string shell命令
//	@param logFunc func(scanner *bufio.Scanner) 获取输出结果function
//	@return error 执行异常时返回
func (s *session) ExecOutput(shell string, logFunc func(scanner *bufio.Scanner)) error {
	var waitGroup sync.WaitGroup
	if logFunc != nil {
		stdout, err := s.sshSess.StdoutPipe()
		if err != nil {
			return err
		}
		waitGroup.Add(1)
		go func() {
			logFunc(bufio.NewScanner(stdout))
			waitGroup.Done()
		}()
	} else {
		var b bytes.Buffer
		s.sshSess.Stdout = &b
		s.output = &b
	}
	err := s.sshSess.Start(shell)
	waitGroup.Wait()
	return err
}

// Wait 等待执行
//
//	@author duanzt
//	@date 2023-07-14 10:12:51
//	@return error 异常时返回
func (s *session) Wait() error {
	return s.sshSess.Wait()
}

// Close 关闭ssh连接
//
//	@author duanzt
//	@date 2023-07-14 10:13:28
//	@return error
func (s *session) Close() error {
	return s.sshSess.Close()
}

// Output Exec执行完后调用，获取执行shell输出结果
//
//	@author duanzt
//	@date 2023-07-14 10:13:37
//	@return string 执行shell输出结果
func (s *session) Output() string {
	o := s.output.String()
	if strings.HasPrefix(o, "[sudo]") {
		return strings.TrimPrefix(strings.SplitN(o, ":", 2)[1], " ")
	}
	return o
}
