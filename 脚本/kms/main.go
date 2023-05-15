package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Exist(FileName string) bool {
	_, err := os.Stat(FileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	fmt.Println("文件错误")
	return false
}

func CopyFile(dstName, srcName string) error {
	src, err := ioutil.ReadFile(srcName)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dstName, src, 0666)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	dst := "D:/abc.txt"
	src := "D:/111.txt"
	// src1,_:=os.Open(src)
	// dst1,_:=os.Create(dst)
	// io.Copy(src1,dst1)

	bo := Exist(dst)
	if bo {
		fmt.Println("文件已存在")
	} else {
		if b := Exist(src); b {
			if err := CopyFile(dst, src); err != nil {
				fmt.Println("读写文件时发生错误")
			}
			fmt.Println("copy成功")
		} else {
			fmt.Println("源文件不存在，请检查路径")
		}

	}
}
