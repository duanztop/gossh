/*
 * @Author: duanzt
 * @Date: 2023-07-14 09:50:38
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-14 11:13:37
 * @FilePath: isession.go
 * @Description: 定义session interface
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package internal

import (
	"bufio"
)

// ISession sesion interface
type ISession interface {

	// Exec 执行shell
	//  @author duanzt
	//  @date 2023-07-14 10:09:53
	//  @param shell string shell命令
	//  @return error 执行异常时返回
	Exec(shell string) error

	// ExecOutput 执行并同步获取输出结果
	//  @author duanzt
	//  @date 2023-07-14 10:12:33
	//  @param shell string shell命令
	//  @param logFunc func(scanner *bufio.Scanner) 获取输出结果function
	//  @return error 执行异常时返回
	ExecOutput(shell string, logFunc func(scanner *bufio.Scanner)) error

	// Wait 等待执行
	//  @author duanzt
	//  @date 2023-07-14 10:12:51
	//  @return error 异常时返回
	Wait() error

	// Close 关闭ssh连接
	//  @author duanzt
	//  @date 2023-07-14 10:13:28
	//  @return error
	Close() error

	// Output Exec执行完后调用，获取执行shell输出结果
	//  @author duanzt
	//  @date 2023-07-14 10:13:37
	//  @return string 执行shell输出结果
	Output() string
}
