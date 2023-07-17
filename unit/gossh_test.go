/*
 * @Author: duanzt
 * @Date: 2023-07-14 17:21:14
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-17 12:27:02
 * @FilePath: gossh_test.go
 * @Description: 单元测试相关代码
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package unit

import (
	"context"
	"testing"

	"github.com/duanztop/gossh"
)

// TestGetRemote1 测试获取远程连接1，并执行shell
//
//	@author duanzt
//	@date 2023-07-14 05:21:52
//	@param t *testing.T
func TestGetRemote1(t *testing.T) {
	con, err := gossh.Remote1("root", "password", "xxx.xxx.xxx.xxx:22")
	if err != nil {
		t.Error(err)
		return
	}
	defer con.Close()
	s, err2 := con.ExecShell(context.Background(), "df -h")
	if err2 != nil {
		t.Error(err2)
		return
	}
	t.Logf(s)
}

// TestGetRemote2 测试获取远程连接2，并执行shell
//
//	@author duanzt
//	@date 2023-07-14 05:53:26
//	@param t *testing.T
func TestGetRemote2(t *testing.T) {
	con, err := gossh.Remote2("root", "/root/.ssh/id_rsa", "xxx.xxx.xxx.xxx:22")
	if err != nil {
		t.Error(err)
		return
	}
	defer con.Close()
	s, err2 := con.ExecShell(context.Background(), "df -h")
	if err2 != nil {
		t.Error(err2)
		return
	}
	t.Logf(s)
}

// TestGetRemoteDefault 测试获取远程连接默认方式，并执行shell
//
//	@author duanzt
//	@date 2023-07-14 06:29:04
//	@param t *testing.T
func TestGetRemoteDefault(t *testing.T) {
	con, err := gossh.RemoteDefault("xxx.xxx.xxx.xxx:22")
	if err != nil {
		t.Error(err)
		return
	}
	defer con.Close()
	s, err2 := con.ExecShell(context.Background(), "df -h")
	if err2 != nil {
		t.Error(err2)
		return
	}
	t.Logf(s)
}

// TestGetLocal 测试获取本地连接，并执行shell
//
//	@author duanzt
//	@date 2023-07-14 05:53:32
//	@param t *testing.T
func TestGetLocal(t *testing.T) {
	l := gossh.Local()
	defer l.Close()
	s, err := l.ExecShell(context.Background(), "df -h")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf(s)
}

func TestLocalFile(t *testing.T) {
	src := "./Mario.jpeg"
	dest := "./Mario_copy.jpeg"
	mode := "0755"
	l := gossh.Local()
	defer l.Close()
	err := l.CopyFileRTL(src, dest, mode)
	if err != nil {
		t.Error("CopyFileRTL failed")
	} else {
		t.Log("CopyFileRTL succeeded")
	}
	// dataChan, err := l.CopyFileRTLMon(src, dest, mode)
}
