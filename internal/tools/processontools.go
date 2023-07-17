/*
 * @Author: duanzt
 * @Date: 2023-07-17 09:28:40
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-17 14:18:09
 * @FilePath: processontools.go
 * @Description: 一个简易的进度条工具
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package tools

import (
	"fmt"
	"strings"
	"sync/atomic"
)

const (
	graph           = "#"
	graphBackground = " "
	lable           = "|/-\\"

	// printGraphNumber 打印进度条长度
	printGraphNumber = 100
)

// processontools 进度条工具类
type processontools struct{}

var (
	PorcessOnTools        = processontools{}
	lableUnit      uint64 = 0
)

// Print 一个简易的打印进度条方法
//
//	@author duanzt
//	@date 2023-07-17 12:28:45
//	@receiver processontools
//	@param currValue int64 当前大小
//	@param totalVale int64 期望大小
func (processontools) Print(currValue, totalVale int64) {

	// 比例
	rate := float32((currValue + 1) * 100 / totalVale)
	upperLimit := int(printGraphNumber * rate / 100)

	// [####################################################################################################][100%][|]
	// 拼接进度#号部分
	back := ""
	for i := 0; i < 100; i++ {
		back += graphBackground
	}
	fmt.Printf("\r[%s][%.2f%%][%s]", strings.Replace(back, graphBackground, graph, upperLimit), rate, string(lable[atomic.AddUint64(&lableUnit, 1)%4]))
}

// PrintEnd 进度终止时打印
//
//	@author duanzt
//	@date 2023-07-17 12:29:11
//	@receiver processontools
func (processontools) PrintEnd() {
	fmt.Print("\n")
}
