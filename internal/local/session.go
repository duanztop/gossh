/*
 * @Author: duanzt
 * @Date: 2023-07-14 10:27:28
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-14 18:07:33
 * @FilePath: session.go
 * @Description: 本地session
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package local

import (
	"bufio"
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"sync"
)

// session 本地session
type session struct {
	sess   *exec.Cmd
	ctx    context.Context
	output *bytes.Buffer
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
	var stdout bytes.Buffer
	sess := exec.CommandContext(s.ctx, "sh", "-c", shell)
	if runtime.GOOS == "windows" {
		sess = exec.CommandContext(s.ctx, "cmd", "/c", shell)
	}
	if logFunc == nil {
		sess.Stdout = &stdout
	}
	sess.Stderr = &stdout
	s.output = &stdout
	s.sess = sess
	var waitGroup sync.WaitGroup
	if logFunc != nil {
		waitGroup.Add(1)
		stdout, err := sess.StdoutPipe()
		if err != nil {
			return err
		}
		go func() {
			logFunc(bufio.NewScanner(stdout))
			waitGroup.Done()
		}()
	}
	err := sess.Start()
	waitGroup.Wait()
	return err
}

// Wait 等待执行
//
//	@author duanzt
//	@date 2023-07-14 10:12:51
//	@return error 异常时返回
func (s *session) Wait() error {
	return s.sess.Wait()
}

// Close 关闭ssh连接
//
//	@author duanzt
//	@date 2023-07-14 10:13:28
//	@return error
func (s *session) Close() error {
	return nil
}

// Output Exec执行完后调用，获取执行shell输出结果
//
//	@author duanzt
//	@date 2023-07-14 10:13:37
//	@return string 执行shell输出结果
func (s *session) Output() string {
	return s.output.String()
}
