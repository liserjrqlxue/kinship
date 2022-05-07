package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

// os
var (
	ex, _        = os.Executable()
	exPath       = filepath.Dir(ex)
	srcPath      = filepath.Join(exPath, "src")
	templatePath = filepath.Join(exPath, "template")
)

type Infos struct {
	Option  string
	Img     string
	Src     string
	Token   string
	Title   string
	Err     string
	Message string
	Href    string
}

func md5sum(str string) string {
	byteStr := []byte(str)
	sum := md5.Sum(byteStr)
	sumStr := fmt.Sprintf("%x", sum)
	return sumStr
}

func createToken() string {
	// token
	return md5sum(strconv.FormatInt(time.Now().Unix(), 10))
}

func logRequest(r *http.Request) {
	//这些信息是输出到服务器端的打印信息
	log.Println(r.Form)
	log.Println("path", r.URL.Path)
	log.Println("scheme", r.URL.Scheme)
	log.Println(r.Form["url_long"])
	for k, v := range r.Form {
		log.Printf("key:%s\t", k)
		if len(v) < 1024 {
			log.Printf("key:[%s]\tval:[%v]\n", k, v)
		} else {
			log.Printf("key:[%s]\tval: large data!\n", k)
		}
	}
}

func printMsg(w http.ResponseWriter, msg ...any) {
	log.Println(msg)
	var _, err = fmt.Fprintln(w, msg)
	if err != nil {
		log.Println(err)
	}
}

func uploadFile(r *http.Request, outDir, key string) (infos []Infos, e error) {
	var fhs, ok = r.MultipartForm.File[key]
	if !ok {
		e = errors.New(key + "not found!")
		return
	}
	for _, fh := range fhs {
		var f io.ReadCloser
		f, e = fh.Open()
		if e != nil {
			return
		}
		var sPath = outDir + "/" + fh.Filename
		var info = Infos{
			Href:    sPath,
			Message: fh.Filename,
		}
		e = upload(f, sPath)
		if e != nil {
			return
		}
		e = f.Close()
		if e != nil {
			return
		}
		infos = append(infos, info)
	}
	return
}

func upload(file io.Reader, dest string) error {
	var f, err = os.Create(dest)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, file)
	if err == nil {
		err = f.Close()
	} else {
		_ = f.Close()
	}
	return err
}

func homepage(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	//解析url传递的参数，对于POST则解析响应包的主体（request body）
	var err = r.ParseForm()
	if err != nil {
		printMsg(w, err)
		return
	}

	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	logRequest(r)

	//fmt.Fprintf(w, "<script>alert('good')</script>") //这个写入到w的是输出到客户端的
	t, err = template.ParseFiles(templatePath+"/header.html", templatePath+"/footer.html", templatePath+"/index.html")
	if err != nil {
		printMsg(w, err)
		return
	}

	var Info Infos
	Info.Title = "Home Page"
	Info.Token = createToken()
	err = t.ExecuteTemplate(w, "index", Info)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func kinship(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	var err error
	var Info Infos

	Info.Title = "Kinship"
	Info.Token = createToken()
	fmt.Println("method:", r.Method)

	if r.Method == "POST" {
		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			printMsg(w, err)
			return
		}

		logRequest(r)

		var outDir = "public/" + Info.Title + "/" + Info.Token
		var workdir = filepath.Join("public", Info.Title, Info.Token)
		err = os.MkdirAll(workdir, 0755)
		if err != nil {
			printMsg(w, err)
			return
		}

		var probandInfos, fatherInfos, motherInfos []Infos
		probandInfos, err = uploadFile(r, outDir, "proband")
		if err != nil {
			printMsg(w, err)
			return
		}
		fatherInfos, err = uploadFile(r, outDir, "father")
		if err != nil {
			printMsg(w, err)
			return
		}
		motherInfos, err = uploadFile(r, outDir, "mother")
		if err != nil {
			printMsg(w, err)
			return
		}

		var cmd = []string{
			"-m", "ss",
			"--child", probandInfos[0].Href,
			"--dad", fatherInfos[0].Href,
			"--mom", motherInfos[0].Href,
			"--out", outDir + "/" + Info.Title,
		}
		var run = exec.Command("python", cmd...)
		log.Println(run.String())
		log.Printf("PYTHONPATH:%s\n", os.Getenv("PYTHONPATH"))
		err = os.Setenv("PYTHONPATH", srcPath)
		if err != nil {
			printMsg(w, err)
			return
		}
		log.Printf("PYTHONPATH:%s\n", os.Getenv("PYTHONPATH"))
		err = os.Setenv("PATH", srcPath+"/ss/bin:"+os.Getenv("PATH"))
		if err != nil {
			printMsg(w, err)
			return
		}
		log.Printf("PATH:%s\n", os.Getenv("PATH"))
		var output []byte
		output, err = run.CombinedOutput()
		if err != nil {
			printMsg(w, string(output)+"\n", err)
			return
		} else {
			http.Redirect(w, r, filepath.Join(workdir, Info.Title), http.StatusSeeOther)
			return
		}
	} else {
		t, err = template.ParseFiles(templatePath+"/header.html", templatePath+"/footer.html", fmt.Sprintf("%s/%s.html", templatePath, Info.Title))
		if err != nil {
			printMsg(w, err)
		}
		err = t.ExecuteTemplate(w, Info.Title, Info)
		if err != nil {
			printMsg(w, err)
			return
		}
	}
}
