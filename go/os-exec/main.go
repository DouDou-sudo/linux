package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path"
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
	t := time.NewTicker(time.Second)
	// for i := 0; i < 3; i++ {
	// 	<-t.C
	// 	fmt.Println(i)
	// }
	<-t.C
	fmt.Println("123")
	t.Stop()
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
	// }()
	// err = c.Start()
	// wg.Wait()
	// fmt.Println("123")
	// c.Wait()
	return err
}

//一次性获取所有输出(stderr and stdout)再打印，windows下会乱码
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
	if err != nil {
		return fmt.Errorf("读取文件%s失败,%v", s, err)
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			str, err := reader.ReadString('\n') // 每次读取一行
			if err == nil {
				fmt.Print(str)
			} else {
				if err == io.EOF { // 读到文件末尾
					fmt.Print(str)
				} else {
					return fmt.Errorf("读取文件%s错误,%v", s, err)
				}
				break
			}
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

func main1() string {
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
