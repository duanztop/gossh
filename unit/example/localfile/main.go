/*
 * @Author: duanzt
 * @Date: 2023-07-17 12:30:24
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-17 14:14:57
 * @FilePath: main.go
 * @Description: 本地连接操作文件的使用示例
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

	// 本段代码步骤解读： 将当前目录下Mario.jpeg文件拷贝到target/Mario.jpeg中，并监控拷贝进度
	// 由于本地ssh文件的操作底层都是调用的CopyFileITR/CopyFileITRMon。所以我们仅测试一种场景：本地拷贝到目标
	// 大文件可以使用testFile进行测试（通过dd if=/dev/zero of=testFile bs=1G count=1创建）
	l := gossh.Local()
	src := "/Users/duanzt/code/golang/gossh/unit/example/localfile/testFile"
	dest := "/Users/duanzt/code/golang/gossh/unit/example/localfile/target/testFile"
	os.Remove(dest)
	mode := "0755"
	file, err := os.Open(src)
	if err != nil {
		panic(fmt.Sprintf("打开文件异常,%s", err.Error()))
	}

	destSizeChan := make(chan int64, 100)

	go func() {
		err = l.CopyFileITRMon(file, dest, mode, destSizeChan)
		if err != nil {
			fmt.Println("操作发生异常：")
			panic(err.Error())
		}
	}()

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
}
