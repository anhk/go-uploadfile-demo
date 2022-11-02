package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

// 处理 /upload  逻辑
func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) // 获取请求的方法
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			panic(err)
		}
		fmt.Printf("#Form +%v\n", r.Form)
		fmt.Printf("#PostForm +%v\n", r.PostForm)
		fmt.Printf("#MultipartForm +%v\n", r.MultipartForm)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println("+++", err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
		if err != nil {
			fmt.Println("---", err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	//fmt.Println(32 << 20)
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":8081", nil)
}
