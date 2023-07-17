/*
 * @Author: duanzt
 * @Date: 2023-07-17 16:04:32
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-17 16:17:36
 * @FilePath: main.go
 * @Description: 远程ssh连接操作文件的使用示例
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package main

import (
	"fmt"
	"os"

	"github.com/duanztop/gossh"
	"github.com/duanztop/gossh/internal/tools"
)

func main() {
	// 本段代码步骤解读：
	// 1.将当前目录下Mario.jpeg文件拷贝到远程服务器/tmp/Mario.jpeg文件，并监控拷贝进度
	// 2.将远程服务器/tmp/Mario.jpeg文件拷贝到本地target/Mario.jpeg中，并监控拷贝进度
	// 大文件可以使用testFile进行测试（通过dd if=/dev/zero of=testFile bs=1G count=1创建）
	conn, err := gossh.RemoteDefault("172.31.54.6:22")
	if err != nil {
		panic(err)
	}

	mode := "0755"

	// =======================第一步=======================
	remoteSrc := "./Mario.jpeg"
	remoteDest := "/tmp/Mario.jpeg"
	destSizeChan := make(chan int64)
	go func() {
		if err := conn.CopyFileLTRMon(remoteSrc, remoteDest, mode, destSizeChan); err != nil {
			panic(err)
		}
	}()

	file, err := os.Open(remoteSrc)
	if err != nil {
		panic(fmt.Sprintf("打开文件异常,%s", err.Error()))
	}
	stat, err := file.Stat()
	if err != nil {
		panic(fmt.Sprintf("获取文件大小异常,%s", err.Error()))
	}
	size := stat.Size()

	for {
		v, ok := <-destSizeChan // 获取数据
		if ok {
			tools.PorcessOnTools.Print(v, size)
		} else {
			tools.PorcessOnTools.PrintEnd()
			break
		}
	}

	fmt.Println("第一步执行成功")

	// =======================第二步=======================
	localDest := "./target/Mario.jpeg"
	destSizeChan = make(chan int64)
	go func() {
		if err := conn.CopyFileRTLMon(remoteDest, localDest, mode, destSizeChan); err != nil {
			panic(err)
		}
	}()
	for {
		v, ok := <-destSizeChan // 获取数据
		if ok {
			tools.PorcessOnTools.Print(v, size)
		} else {
			tools.PorcessOnTools.PrintEnd()
			break
		}
	}
	fmt.Println("第二步执行成功")
}
