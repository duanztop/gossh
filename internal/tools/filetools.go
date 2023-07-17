/*
 * @Author: duanzt
 * @Date: 2023-07-14 18:21:26
 * @LastEditors: duanzt
 * @LastEditTime: 2023-07-17 13:16:55
 * @FilePath: filetools.go
 * @Description: 文件处理工具
 *
 * Copyright (c) 2023 by duanzt, All Rights Reserved.
 */
package tools

import (
	"io"
	"io/fs"
	"os"
	"strings"
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

// PathExists 判断一个文件或文件夹是否存在，输入文件路径，根据返回的bool值来判断文件或文件夹是否存在
//
//	@author duanzt
//	@date 2023-07-17 08:55:18
//	@receiver filetools
//	@param path string 文件/文件夹路径
//	@return bool true表示存在
//	@return error 判断异常时返回
func (filetools) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CreateFile 创建文件
//
//	@author duanzt
//	@date 2023-07-17 12:48:55
//	@receiver f filetools
//	@param path string 文件路径
//	@return *os.File 返回创建完的文件对象
//	@return error 创建失败时返回
func (f filetools) CreateFile(path string) (*os.File, error) {
	// 截取最后一个[/]判断目录是否存在，不存在则进行创建
	ss := strings.Split(path, "/")
	dir := ss[len(ss)-2]
	if dirF, err := f.PathExists(dir); err != nil {
		return nil, err
	} else if !dirF {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 判断文件是否存在，不存在则进行创建
	var destFile *os.File
	if destF, err := f.PathExists(path); err != nil {
		return nil, err
	} else if !destF {
		destFile, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	}
	return destFile, nil
}
