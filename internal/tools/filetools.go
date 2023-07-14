/*
 * @Author: duanzt
 * @Date: 2023-07-14 18:21:26
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-14 18:23:00
 * @FilePath: filetools.go
 * @Description: 文件处理工具
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package tools

import (
	"io"
	"io/fs"
)

type filetools struct{}

var (
	FileTools = filetools{}
)

// ReadFile 读取文件
//
//	@author duanzt
//	@date 2023-07-14 06:23:21
//	@receiver filetools
//	@param f fs.File 文件
//	@return []byte 文件内容
//	@return error 读取时发生异常
func (filetools) ReadFile(f fs.File) ([]byte, error) {
	var size int
	if info, err := f.Stat(); err == nil {
		size64 := info.Size()
		if int64(int(size64)) == size64 {
			size = int(size64)
		}
	}
	size++

	if size < 512 {
		size = 512
	}

	data := make([]byte, 0, size)
	for {
		if len(data) >= cap(data) {
			d := append(data[:cap(data)], 0)
			data = d[:len(data)]
		}
		n, err := f.Read(data[len(data):cap(data)])
		data = data[:len(data)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return data, err
		}
	}
}
