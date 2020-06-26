package main

import (
	// "encoding/json"
	"errors"
	"flag"
	// "fmt"
	"html/template"
	"io"
	// "io/ioutil" //dbg
	"log"
	"net/http"
	"os"
	// "path/filepath"
	"strings"
	"devcn.fun/infrastlabs/ingsitemap/asset"
	"devcn.fun/infrastlabs/ingsitemap/ingress"
	"github.com/elazarl/go-bindata-assetfs"
)

var (
	addr  = flag.String("addr", ":9010", "请输入服务端地址")
)

func hasSuffix(url string, prefix []string) bool {
	for _, p := range prefix {
		if strings.HasSuffix(url, p) {
			return true
		}
	}
	return false
}

func handleFuncHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	if hasSuffix(r.URL.Path, []string{".jpg", ".css", ".png", ".png", ".js", ".gif"}) {
		w.Header().Add("Cache-Control", "max-age=604800, must-revalidate")
		w.Header().Add("Pragma", "public")

	} else {
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Add("Pragma", "no-cache")
		w.Header().Add("Expires", "0")
	}
}

func renderhtml(filename string, out io.Writer) error {
	bytes, err := asset.Asset("index-tpl.html")//former org-index //without static/
	// bytes, err := ioutil.ReadFile("./static/index-tpl.html") //dbg
	if err != nil {
		return errors.New("no found home  template  html")
	}

	/* mydata := []data{{
		Url:   "page-1.html",
		Title: "go to page 1",
	}, {
		Url:   "page-2.html",
		Title: "go to page 2",
	}, {
		Url:   "page-3.html",
		Title: "go to page 3",
	}, {
		Url:   "page-5.html",
		Title: "go to page 5",
	}}
	webData := show{mydata} */
	// webData := new(show)
	// webData.Pages= getDatas()
	webData := ingress.Show{ingress.GetIngs()}

	//TODO replace key  //sam: go template
	return template.Must(template.New("markdown").Parse(string(bytes))).Execute(out, webData)
}

func handlePage(w http.ResponseWriter, r *http.Request) {
	handleFuncHttp(w, r)

	var code = 200
	var err error
	defer func() {
		if err != nil {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.WriteHeader(code)
			io.WriteString(w, err.Error())
		}
	}()

	r.ParseForm()
	err = renderhtml("emp", w)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	if os.IsNotExist(err) {
		code = 404
	}
	return
}

func main() {
	flag.Parse()
	ADDR := os.Getenv("ADDR")
	if ADDR != "" {
		addr = &ADDR
	}

	fs := assetfs.AssetFS{
        Asset:     asset.Asset,
        AssetDir:  asset.AssetDir,
        AssetInfo: asset.AssetInfo,
	}
	// http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs")))) //yaml|json
	http.HandleFunc("/page", handlePage)//handel to index-tpl.html
	
	// http.Handle("/", http.FileServer(http.Dir("static")))//for dbg
	http.Handle("/", http.FileServer(&fs))
	
	log.Printf("Listening on %s  ", *addr)
	http.ListenAndServe(*addr, nil)
}
