package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// JsonTest("test")
	// fmt.Print(BufReadFile("test"))
	// Gin_test()
	// GinHtml()
	// time.Now().Date()
	// GinHtmls()
	// GinCustom()
	// MultiGin()
	// MultiGinForm()
	// QueryPostFormGin()
	// SecureJson()
	// UploadFileGin()
	// t := time.Now()
	// // fmt.Print(time.Now().Sub(t).Hours())
	// fmt.Println(t)
	// ts := t.Format("2006-01-02 15:04:05")
	// fmt.Printf("%q\n", ts)
	// time, _ := time.ParseInLocation("2006-01-02 15:04:05", ts, time.Local)
	// fmt.Printf("%q", time)
	TimeAf()
}

func TimeAf() {
	start := time.Now()
	c := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		// time.Sleep(3 * time.Second)
		<-c
	}()
	select {
	case c <- 1:
		fmt.Println("channel...")
	case <-time.After(2 * time.Second):
		close(c)
		fmt.Println("timeout...")
	}
	fmt.Println(time.Now().Sub(start).Seconds())
}

func UploadFilesGin() {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", func(ctx *gin.Context) {
		form, _ := ctx.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Panicln(file.Filename)
			dst := "./" + file.Filename
			ctx.SaveUploadedFile(file, dst)
		}
		ctx.String(200, fmt.Sprintf("'%s' uploaded!", len(files)))
	})
	r.Run()
}

func UploadFileGin() {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20
	r.POST("/upload", func(ctx *gin.Context) {
		file, _ := ctx.FormFile("file")
		log.Println(file.Filename)
		dst := "./" + file.Filename
		ctx.SaveUploadedFile(file, dst)
		ctx.String(200, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	r.Run()
}

func SecureJson() {
	r := gin.Default()
	r.GET("/securejson", func(ctx *gin.Context) {
		ctx.SecureJSON(200, []string{"lena", "ays", "fo"})
	})
	r.Run()
}

func QueryPostFormGin() {
	r := gin.Default()
	r.POST("/post", func(ctx *gin.Context) {
		id := ctx.Query("id")
		page := ctx.DefaultQuery("page", "0")
		name := ctx.PostForm("name")
		messgae := ctx.PostForm("message")
		fmt.Printf("id: %s;page: %s;name: %s;messge: %s", id, page, name, messgae)
	})
	r.Run()
}

func MultiGinForm() {
	r := gin.Default()
	r.POST("/form_post", func(ctx *gin.Context) {
		message := ctx.PostForm("message")
		nick := ctx.DefaultPostForm("nick", "anonymous")
		ctx.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	r.Run()
}

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password binding:"required"`
}

func MultiGin() {
	router := gin.Default()
	router.POST("/login", func(c *gin.Context) {
		// 你可以使用显式绑定声明绑定 multipart form：
		// c.ShouldBindWith(&form, binding.Form)
		// 或者简单地使用 ShouldBind 方法自动绑定：
		var form LoginForm
		// 在这种情况下，将自动选择合适的绑定
		if c.ShouldBind(&form) == nil {
			if form.User == "user" && form.Password == "password" {
				c.JSON(200, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(401, gin.H{"status": "unauthorized"})
			}
		}
	})
	router.Run(":8080")
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()

	return fmt.Sprintf("%d/%02d/%02d ", year, month, day)
}

func GinCustom() {
	r := gin.Default()
	r.Delims("{[{", "}]}")
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHTMLFiles("./testdata/template/raw.tmpl")
	r.GET("/raw", func(ctx *gin.Context) {
		ctx.HTML(200, "raw.tmpl", gin.H{
			"now": time.Now(),
		})
	})
	r.Run()
}

func GinHtmls() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	r.GET("/posts/index", func(ctx *gin.Context) {
		ctx.HTML(200, "posts/index.tmpl", gin.H{
			"title": "posts",
		})
	})
	r.GET("/user/index", func(ctx *gin.Context) {
		ctx.HTML(200, "user/index.tmpl", gin.H{
			"title": "user",
		})
	})
	r.Run()
}

func GinHtml() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
	router.Run(":8080")
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
