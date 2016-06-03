// Copyright 2015 someonegg. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package fileutil

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	// 默认的文件权限
	DefaultFileMode os.FileMode = os.ModePerm

	// linux下需加上O_WRONLY或是O_RDWR
	DefaultFileFlag int = os.O_APPEND | os.O_CREATE | os.O_WRONLY
)

//构建一个每日写日志文件的写入器
func OpenFile(pathfile string, filemode os.FileMode, fileflag int) (wc *os.File, err error) {
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

func FileExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

//eg:filename=`../a/test/aa.txt` newExt=`.csv` -->return=`../a/test/aa.csv`
func ChangeFileExt(filename, newExt string) string {
	filename = filepath.ToSlash(filename)
	file := path.Base(filename)
	file = strings.TrimSuffix(file, path.Ext(file)) + newExt
	dir := path.Dir(filename)
	if dir == "." {
		dir = dir + "/" + file
	} else {
		dir = path.Join(dir, file)
	}
	return filepath.ToSlash(dir)
}
func PathJoin(dir, filename string) string {
	if dir == "." {
		return filepath.ToSlash(dir + "/" + filename)
	} else {
		return filepath.ToSlash(path.Join(dir, filename))
	}
}
