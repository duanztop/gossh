/*
 * @Author: duanzt
 * @Date: 2023-07-17 12:23:38
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-17 12:43:27
 * @FilePath: main.go
 * @Description: 进度条的使用示例
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package main

import (
	"time"

	"github.com/duanztop/gossh/internal/tools"
)

func main() {
	var len int64 = 100
	var i int64 = 0
	for ; i < len; i++ {
		tools.PorcessOnTools.Print(i, len)
		time.Sleep(200 * time.Millisecond)
	}
	tools.PorcessOnTools.PrintEnd()
}
