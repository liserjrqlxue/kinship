package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/liserjrqlxue/goUtil/simpleUtil"
)

var (
	port = flag.String(
		"port",
		":9091",
		"web server listen port",
	)
)

var StaticDir = make(map[string]string)

func main() {
	flag.Parse()
	StaticDir["/static"] = "static"
	StaticDir["/public"] = "public"

	// 设置访问的路由
	//http.HandleFunc("/Web_url_name", func_name)
	http.HandleFunc("/kinship", kinship)
	// 设置主页与ServeFile
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// static file server
		for prefix, staticDir := range StaticDir {
			if strings.HasPrefix(r.URL.Path, prefix) {
				file := staticDir + r.URL.Path[len(prefix):]
				fmt.Println(file)
				http.ServeFile(w, r, file)
				return
			}
		}
		homepage(w, r)
	})

	//设置监听的端口
	fmt.Printf("start http://localhost%v\n", *port)
	simpleUtil.CheckErr(http.ListenAndServe(*port, nil))
}
