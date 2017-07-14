package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const usage string = `fileserver usage:\neasy_server [-port port] [-root rootdirectory]\n
If has no argument, will use ./ as root folder as default, and 31100 as default port.\n
Use \"gofileserver -help\" to show this message.`

var (
	svrHandler http.Handler
)

func main() {
	port := flag.Int("port", 31100, "[-port xxxx], e.g. ... -port 31100 ...")
	rootDir := flag.String("root", "./", "[-root xxxx], e.g. ... -root C:/ ...")
	help := flag.Bool("help", false, usage)
	flag.Parse()
	if *help || flag.NArg() > 4 {
		println(usage)
		os.Exit(0)
	}

	svrHandler = http.FileServer(http.Dir(*rootDir))
	http.HandleFunc("/", RecordServer)
	println("#Server: Running http file server...\n#Server: Root folder is: " + *rootDir, "\n#Server: Port is :", *port)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatal("\n#Server: LstenAndServer: ", err, "\n#Server: port is :", string(*port))
	}
}

func RecordServer(w http.ResponseWriter, req *http.Request) {
	println("# INFO: ", time.Now().Format("2006-01-02 15:04:05 -0700 MST"))
	println("        Remote address is : ", req.RemoteAddr)
	println("        Request URL is :    ", req.URL.Path)
	println()
	svrHandler.ServeHTTP(w, req)
}
