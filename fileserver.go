package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const usage string = `fileserver usage:\neasy_server [-port port] [-root rootdirectory]\n
If has no argument, will use ./ as root folder as default, and 31100 as default port.\n
Use \"gofileserver -help\" to show this message.`

var (
	svrHandler http.Handler
	rootDir    *string
)

func main() {
	port := flag.Int("port", 31100, "[-port xxxx], e.g. ... -port 31100 ...")
	rootDir = flag.String("root", "./", "[-root xxxx], e.g. ... -root C:/ ...")
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
		log.Fatal("\n#Server:\nLstenAndServer: ", err, "\n#Server: port is : ", string(*port))
	}
}

func RecordServer(w http.ResponseWriter, req *http.Request) {
	println("# INFO: ", time.Now().Format("2006-01-02 15:04:05 -0700 MST"))
	println("        Remote address is : ", req.RemoteAddr)
	println("        Remote request is : ", req.Method)
	println("        Request URL is :    ", req.URL.Path)
	println()
	if req.Method == "GET" {
		if strings.HasSuffix(req.RequestURI, "/") {
			header := "<head>\n<title>File Server</title>\n\n</head>\n"
			w.Write([]byte(header))

			body1 := "<body>\n"
			w.Write([]byte(body1))

			link := `<git>Link on github: </git><a href="https://github.com/ToolsPlease/gofileserver">github/</a>`
			w.Write([]byte(link))

			w.Write([]byte("<form method=\"POST\" " + " enctype=\"multipart/form-data\">" + "Choose a file to upload: <input name=\"ufile\" type=\"file\" />" + "<input type=\"submit\" value=\"Upload\" />" + "</form>"))

			w.Write([]byte(`<local>`))
			svrHandler.ServeHTTP(w, req)
			w.Write([]byte(`</local>`))

			body2 := "</body>"
			w.Write([]byte(body2))
		} else {
			svrHandler.ServeHTTP(w, req)
		}
	} else if req.Method == "POST" {
		f, h, err := req.FormFile("ufile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sep := strings.LastIndexAny(h.Filename, `\/`)
		filename := string(h.Filename[sep+1:])
		println(filename)
		defer f.Close()
		t, err := os.Create(*rootDir + filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, req, req.RequestURI, http.StatusFound)
	} else {
		svrHandler.ServeHTTP(w, req)
	}
}
