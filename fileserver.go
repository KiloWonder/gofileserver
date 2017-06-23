package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
)

const usage string = "fileserver usage:\neasy_server [-port port] [-root rootdirectory]\nIf has no argument, will use ./ as root folder as default, and 31100 as default port.\nUse \"gofileserver -help\" to show this message."

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
	http.Handle("/", svrHandler)
	println("#Go: Running http server...\n#Go: Root folder is: "+*rootDir, "\n#Go: Port is :", *port)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), svrHandler)
	if err != nil {
		log.Fatal("\n#Go: LstenAndServer: ", err, "\n#Go: port is :", string(*port))
	}
}