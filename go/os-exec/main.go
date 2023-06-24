package main

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

func CheckError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
func main() {
	// temp("D:/iso")
	// LinuxCommand("ping 127.0.0.1")
	// RetiCommand("ping 127.0.0.1")
	// ctx, cancel := context.WithCancel(context.Background())
	// go func(cancelFunc context.CancelFunc) {
	// 	time.Sleep(3 * time.Second)
	// 	cancelFunc()
	// }(cancel)
	// RetiContCommand(ctx, "ls")
	// SeleCache()
	// a := time.Now().Format("2006-01-02 15:04:05")
	// // fmt.Println(<-time.After(time.Second))
	// <-time.After(time.Second * 3)
	// // time.Now().Format()
	// fmt.Println(a)
	// t := time.NewTicker(time.Second)
	// // for i := 0; i < 3; i++ {
	// // 	<-t.C
	// // 	fmt.Println(i)
	// // }
	// <-t.C
	// fmt.Println("123")
	// t.Stop()
	// time.strftime
	// fmt.Print(exec.LookPath("ls"))
	// IoReadFile("test")
	// BufReadFile("test")
	// fmt.Print(BufReadFile("test"))
	// if err := Command("ping 127.0.0.1"); err != nil {
	// 	fmt.Printf("执行异常，退出代码 1\n%v", err)
	// } else {
	// 	fmt.Println("执行完毕，退出代码 0")
	// }
	fmt.Println(2, 1)
	// 打印log并os.Exit(1)
	log.Fatalln("123")
	// 打印log继续向下运行
	log.Print("123")
	fmt.Print("abc")

}

func upload(w http.ResponseWriter, r *http.Request) {
	h:=md5.New()
	h.Sum(nil)
	r.FormFile()
}

// 可关闭+实时输出
func RetiContCommand(ctx context.Context, cmd string) error {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			select {
			case <-ctx.Done():
				if ctx.Err() != nil {
					fmt.Printf("程序出现错误: %q", ctx.Err())
				} else {
					fmt.Println("程序被终止")
				}
				return
			default:
				readString, err := reader.ReadString('\n')
				if err == nil {
					fmt.Print(readString)
				} else {
					if err == io.EOF {
						fmt.Print(readString)
						// return
						break
					} else {
						// return
						break
					}
				}
			}
		}
	}(&wg)
	err = c.Start()
	wg.Wait()
	return err
}

// 实时显示，top，ping...，windows下会乱码
func RetiCommand(cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	err = c.Start()
	// c.Run()
	fmt.Println("123")
	reader := bufio.NewReader(stdout)
	for {
		readString, err := reader.ReadString('\n')
		fmt.Print(readString)
		if err == io.EOF {
			fmt.Print(readString)
			// return
			break
		} else {
			// return
			break
		}
	}
	// }()
	// err = c.Start()
	// wg.Wait()
	// fmt.Println("123")
	// c.Wait()
	return err
}

// 一次性获取所有输出(stderr and stdout)再打印，windows下会乱码
func Command(cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	output, err := c.CombinedOutput() //linux
	fmt.Print(string(output))
	return err
}

// 将标准输入标准错误重定向到当前终端的标准输出和标准错误，可以实时输出，ping，top等命令
// windows下运行也不会乱码
func LinuxCommand(cmd string) {
	c := exec.Command("sh", "-c", cmd)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Run()
	// os.Args()
}

// 遍历目录
func temp(dirname string) {
	FileInfo, _ := ioutil.ReadDir(dirname)
	for _, v := range FileInfo {
		fmt.Printf("%s\n", v.Name())
		if v.IsDir() {
			temp(path.Join(dirname, v.Name()))
		}
	}
}

// ioutil读取小文件并打印
func IoReadFile(s string) error {
	a, err := ioutil.ReadFile(s)
	if err != nil {
		return fmt.Errorf("读取文件%s失败,%v", s, err)
	}
	fmt.Print(string(a))
	return nil
}

// bufio读取大文件并打印
func BufReadFile(s string) error {
	file, err := os.Open(s)
	file.Sync()
	if err != nil {
		return fmt.Errorf("读取文件%s失败,%v", s, err)
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			str, err := reader.ReadString('\n') // 每次读取一行
			fmt.Print(str)
			if err == io.EOF { // 读到文件末尾
				// fmt.Print(str)
				break
			}
			if err != nil {
				return fmt.Errorf("读取文件%s错误,%v", s, err)
			}
			// if err == nil {
			// 	fmt.Print(str)
			// } else {
			// 	if err == io.EOF { // 读到文件末尾
			// 		fmt.Print(str)
			// 	} else {
			// 		return fmt.Errorf("读取文件%s错误,%v", s, err)
			// 	}

			// }
			// break
		}
	}
	// fmt.Print("文件读取结束!")
	return nil
}

// func BufReadFile(s string) (a string, e error) {
// 	file, err := os.Open(s)
// 	// var a string
// 	if err != nil {
// 		// return a, errors.New(fmt.Sprintf("打开文件%s失败,%v", s, err))
// 		e = fmt.Errorf("读取文件%s失败,%v", s, err)
// 		return
// 	} else {
// 		defer file.Close()
// 		reader := bufio.NewReader(file)
// 		for {
// 			str, err := reader.ReadString('\n')
// 			if err == nil {
// 				a += str
// 			} else if err == io.EOF {
// 				a += str
// 				break
// 			} else {
// 				e = fmt.Errorf("读取文件%s错误,%v", s, err)
// 				break
// 			}
// 		}
// 	}
// 	// fmt.Println("文件读取完毕")
// 	return
// }

func main11() string {
	file, err := os.Open("./hello.go")
	if err != nil {
		fmt.Println("文件读取错误", err)
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			str, err := reader.ReadString('\n') // 每次读取一行
			if err == io.EOF {                  // 读到文件末尾
				break
			} else {
				fmt.Print(str)
			}
		}
	}
	return fmt.Sprint("")
}

// 自定义error信息
type dualError struct {
	Num     float64
	problem string
}

func (e dualError) Error() string {
	return fmt.Sprintf("Wrong!!!,because \"%f\" is a negative number", e.Num)
}
func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return -1, dualError{Num: f}
	}
	return math.Sqrt(f), nil
}

// select
func SeleCache() {
	bufChan := make(chan int, 5)

	go func() {
		time.Sleep(time.Second)
		for {
			// <-bufChan
			bufChan <- 1
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		select {
		// case bufChan <- 1:
		case <-bufChan:
			fmt.Println("add success")
			time.Sleep(time.Second)
		default:
			fmt.Println("资源已满，请稍后再试")
			time.Sleep(time.Second)
		}
	}
}

func CtxFunc1(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("called1", time.Now())
			return
		default:
			time.Sleep(3 * time.Second)
			fmt.Println("chaoshi1")
		}
	}
}
func CtxFunc2(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("called2", time.Now())
			return
		case <-time.After(time.Second * 3):
			fmt.Println("chaoshi2")
		}
	}
}

func main1() {
	ioutil.ReadFile("test")
	// ioutil.WriteFile()
	// json.Unmarshal()
	// c := time.Now()
	// fmt.Println(c)
	// a := c.Unix()
	// fmt.Println(a)
	// b := time.Unix(a, 0)
	// fmt.Println(b)
	// s := "   \t\na b\nc\td\t "
	// s1 := "a,b c d"
	// fmt.Printf("%q\n", strings.Split(s1, ","))
	// fmt.Println(strings.Fields(s))
	// fmt.Printf("%q\n", strings.TrimSpace(s))
	// fmt.Printf("%q\n", strings.Trim(s, " "))
	// fmt.Println(strings.TrimSpace(s))
	// fmt.Println(strings.Trim(s, " "))
	// a, _ := ioutil.ReadFile("test")
	Server := "test"
	a, _ := os.Open(Server)
	defer a.Close()
	b := bufio.NewReader(a)
	for {
		lines, err := b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if strings.Contains(lines, "UNIQB_KSVD_PUBADDR") {
					a := strings.Trim(lines, "\n")
					// a := strings.TrimSpace(lines)
					// a := strings.Trim(lines, "\n")
					b := strings.Split(a, "=")[1]
					fmt.Print(a)
					fmt.Print(b)
				}
			}
			break
		}
		if strings.Contains(lines, "UNIQB_KSVD_PUBADDR") {
			a := strings.TrimSpace(lines)
			// a := strings.Trim(lines, "\n")
			b := strings.Split(a, "=")[1]
			// fmt.Printf("%q\n", a)
			fmt.Println(a)
			// fmt.Printf("%q", b)
			fmt.Println(b)
		}
	}
}

type Student struct {
	Name  string
	Age   int
	Score float64
}

func JsonTest(s string) {
	fileName := s
	stu := Student{
		Name:  "HaiCoder",
		Age:   109,
		Score: 100.5,
	}
	fileContent, err := json.Marshal(stu)
	if err = ioutil.WriteFile(fileName, fileContent, 0666); err != nil {
		fmt.Println("Writefile Error =", err)
		return
	}
	//读取文件
	fileContent, err = ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Read file err =", err)
		return
	}
	fmt.Print(string(fileContent))
	// var stuRes Student
	// if err := json.Unmarshal(fileContent, &stuRes); err != nil {
	// 	fmt.Println("Read file error =", err)
	// } else {
	// 	fmt.Println("Read file success =", stuRes)
	// }
}
func BufReadFile1(s string) (err1 error) {
	file, err := os.Open(s)
	if err != nil {
		return fmt.Errorf("读取文件%s失败,%v", s, err)
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			str, err := reader.ReadString('\n') // 每次读取一行
			fmt.Print(str)
			if err == io.EOF { // 读到文件末尾
				err1 = nil
				fmt.Print("\n文件读取结束!")
				// fmt.Print(str)
				break
			}
			if err != nil {
				err1 = fmt.Errorf("读取文件%s错误,%v", s, err)
				break
			}
		}
	}
	// fmt.Print("文件读取结束!")
	return err1
}

func BufWriFile(s, con string) error {
	var f *os.File
	_, err := os.Stat(s)
	if os.IsNotExist(err) {
		f, _ = os.Create(s)
	} else {
		f, _ = os.OpenFile(s, os.O_APPEND, 0666)
	}
	defer f.Close()
	// io.WriteString
	// if _, err := io.WriteString(f, con); err != nil {
	// 	return err
	// }
	// f.Sync()
	// 无buf
	// if _, err := f.WriteString(con); err != nil {
	// 	return err
	// }
	// f.Sync()

	// bufio写入
	// writer := bufio.NewWriter(f)
	// if _, err := writer.WriteString(con); err != nil {
	// 	return err
	// }
	// writer.Flush()	//必须Flush到文件

	// ioutil会清空原文件，如果原文件不存在会窗口
	// if err := ioutil.WriteFile(s, []byte(con), 0666); err != nil {
	// 	return err
	// }
	return nil
}
