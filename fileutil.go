// Copyright 2016 zxfonline@sina.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package fileutil

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	//DefaultFileMode 默认的文件权限 0644
	DefaultFileMode os.FileMode = 0644

	//DefaultFileFlag linux下需加上O_WRONLY或是O_RDWR
	DefaultFileFlag int = os.O_APPEND | os.O_CREATE | os.O_WRONLY
)
var (
	defaultDirs []string
)

func init() {
	wd, _ := os.Getwd()
	arg0 := path.Clean(os.Args[0])
	var exeFile string
	if strings.HasPrefix(arg0, "/") {
		exeFile = arg0
	} else {
		exeFile = path.Join(wd, arg0)
	}
	parent, _ := path.Split(exeFile)
	//1：命令执行所在目录
	wdDir := TransPath(wd)
	defaultDirs = append(defaultDirs, wdDir)
	exeDir := TransPath(parent)
	if wdDir != exeDir { //2：可执行文件所在目录
		defaultDirs = append(defaultDirs, exeDir)
	}
	exeLastDir := TransPath(path.Join(parent, ".."))
	if wdDir != exeLastDir { //3：可执行文件上一级目录
		defaultDirs = append(defaultDirs, exeLastDir)
	}
}

//InitPathDirs 初始化工程文件根目录 retset：是否重置默认的目录列表 。 rootPaths：新增的目录列表，如果没有设置指定根目录并未重置默认的目录的话，文件搜索规则为: 1：命令执行所在目录、2：可执行文件所在目录、3：可执行文件上一级目录...新增的路径列表，否则直接按照给的文件名称查找文件。
func InitPathDirs(retset bool, rootPaths ...string) {
	if retset {
		defaultDirs = defaultDirs[:0]
	}
	for _, rootPath := range rootPaths {
		if len(rootPath) > 0 {
			rootPath = TransPath(rootPath)
			rootPath, _ = filepath.Abs(rootPath)
			rootPath = TransPath(rootPath)
		}
		if rootPath != "" {
			found := false
			for _, defaultDir := range defaultDirs {
				if defaultDir == rootPath {
					found = true
					break
				}
			}
			if !found {
				defaultDirs = append(defaultDirs, rootPath)
			}
		}
	}
}

//FindFile 查找文件，根据初始化的文件目录顺序查找文件
func FindFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	for _, dir := range defaultDirs {
		fpath := PathJoin(dir, name)
		if FileExists(fpath) {
			f, err := os.OpenFile(fpath, flag, perm)
			if err != nil {
				return nil, fmt.Errorf("open file err,path:%s,err:%v.", fpath, err)
			}
			return f, nil
		}
	}
	fpath := TransPath(name)
	if FileExists(fpath) {
		f, err := os.OpenFile(fpath, flag, perm)
		if err != nil {
			return nil, fmt.Errorf("open file err,path:%s,err:%v.", fpath, err)
		}
		return f, nil
	}
	return nil, fmt.Errorf("file no found,file:%s,dirs:%v.", name, defaultDirs)
}

//OpenFile 打开文件，如果目录文件不存在则创建一个文件
func OpenFile(pathfile string, fileflag int, filemode os.FileMode) (wc *os.File, err error) {
	pathfile = strings.Replace(filepath.Clean(pathfile), "\\", "/", -1)
	dir := path.Dir(pathfile)
	if _, err = os.Stat(dir); err != nil && !os.IsExist(err) {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if err = os.MkdirAll(dir, filemode); err != nil {
			return nil, err
		}
		if _, err = os.Stat(dir); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(pathfile, fileflag, filemode)
}

//DirExists 指定目录是否存在
func DirExists(dir string) bool {
	d, e := os.Stat(dir)
	switch {
	case e != nil:
		return false
	case !d.IsDir():
		return false
	}
	return true
}

//FileExists 指定文件是否存在
func FileExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

//ChangeFileExt 更改文件后缀名 eg:filename=`../a/test/aa.txt` newExt=`.csv` -->return=`../a/test/aa.csv`
func ChangeFileExt(filename, newExt string) string {
	filename = strings.Replace(filepath.Clean(filename), "\\", "/", -1)
	file := path.Base(filename)
	file = strings.TrimSuffix(file, path.Ext(file)) + newExt
	dir := path.Dir(filename)
	return path.Join(dir, file)
}

//PathJoin 路径合并 并将 “\\” 转换成 “/”
func PathJoin(dir, filename string) string {
	return strings.Replace(path.Join(filepath.Clean(dir), filename), "\\", "/", -1)
}

//TransPath 路径连接符转换 将路径 “\\” 转换成 “/”
func TransPath(path string) string {
	return strings.Replace(filepath.Clean(path), "\\", "/", -1)
}
