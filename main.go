package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/uluyol/vplots/conversions"
)

//go:embed index.tmpl.html
var indexTemplStr string

var indexTempl = template.Must(template.New("index").Parse(indexTemplStr))

type plotViewer struct {
	pdfPaths []string
	svgMu    sync.Mutex // protects svgCache
	svgCache [][]byte
	cPNG     *eagerPNGConverter
	c        *conversions.Converter
}

func (v *plotViewer) rootHandler(w http.ResponseWriter, r *http.Request) {
	err := indexTempl.Execute(w, struct {
		Plots []string
	}{
		v.pdfPaths,
	})
	if err != nil {
		log.Println("GET /: %v", err)
	}
}

func (v *plotViewer) pngHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid id: "+err.Error(), http.StatusBadRequest)
		return
	}
	if id < 0 || id >= len(v.pdfPaths) {
		http.Error(w, "image id out of bounds", http.StatusBadRequest)
		return
	}

	width, err := strconv.Atoi(vars["w"])
	if err != nil {
		http.Error(w, "invalid width: "+err.Error(), http.StatusBadRequest)
		return
	}

	out, err := v.cPNG.Get(id, width)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to convert to png: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "image/png")
	w.Write(out)
}

func (v *plotViewer) svgHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id < 0 || id >= len(v.pdfPaths) {
		http.Error(w, "image id out of bounds", http.StatusBadRequest)
		return
	}

	v.svgMu.Lock()
	defer v.svgMu.Unlock()

	if v.svgCache[id] == nil {
		out, err := v.c.ToSVG(v.pdfPaths[id])
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to convert to svg: %v", err), http.StatusInternalServerError)
			return
		}

		v.svgCache[id] = out
	}

	w.Header().Set("content-type", "image/svg+xml")
	w.Write(v.svgCache[id])
}

func (v *plotViewer) pdfHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if id < 0 || id >= len(v.pdfPaths) {
		http.Error(w, "image id out of bounds", http.StatusBadRequest)
		return
	}

	out, err := os.ReadFile(v.pdfPaths[id])
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to read pdf: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "application/pdf")
	w.Write(out)
}

type quitHandler struct {
	s *http.Server
}

func (h quitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)
	h.s.Shutdown(ctx)
}

var listenAddr = flag.String("addr", ":0", "address to listen on")

func Main() {
	log.SetFlags(0)
	log.SetPrefix("vplots: ")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "usage: vplots plot.pdf...")
		os.Exit(2)
	}

	paths := flag.Args()

	c := conversions.NewConverter()
	v := plotViewer{
		pdfPaths: paths,
		svgCache: make([][]byte, len(paths)),
		cPNG:     newEagerPNGConverter(paths, c),
		c:        c,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", v.rootHandler)
	r.HandleFunc("/pngs/{id:[0-9]+}/{w:[0-9]+}", v.pngHandler)
	r.HandleFunc("/svgs/{id:[0-9]+}", v.svgHandler)
	r.HandleFunc("/pdfs/{id:[0-9]+}", v.pdfHandler)

	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Handler: r,
	}
	r.Handle("/quit", quitHandler{s})
	log.Printf("listening on %s", lis.Addr().String())
	go openURL(lis.Addr().String())
	s.Serve(lis)
}

func main() {
	conversions.AppMain(Main)
}

func openURL(url string) {
	url = "http://" + url

	for tries := 5; tries > 5; tries-- {
		if _, err := http.Get(url); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux", "freebsd", "solaris":
		cmd = exec.Command("xdg-open", url)
	default:
		log.Printf("no known open command, connect to %s", url)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Printf("open command failed, connect to %s: %v", url, err)
	}
}
