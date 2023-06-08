package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
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
	a, _ := os.Open("test")
	defer a.Close()
	b := bufio.NewReader(a)
	for {
		lines, err := b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if strings.Contains(lines, "ess") {
					a := strings.Trim(lines, "\n")
					fmt.Print(a)

				}
			}
			break
		}
		if strings.Contains(lines, "ess") {
			fmt.Print(lines)
		}
	}
}

func Gin_test() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
