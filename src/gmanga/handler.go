package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type AppContext struct {
	DataDir   string
	StaticDir string
	Lib       *Library
}

type ContextHandler struct {
	context *AppContext
	handler ContextHandlerFunc
}

type ContextHandlerFunc func(c *AppContext, w http.ResponseWriter, r *http.Request)

func (ch ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch.handler(ch.context, w, r)
}

func NewAppContext() *AppContext {
	c := AppContext{"", "", nil}
	c.Lib = nil
	return &c
}

func DefineRoutes(c *AppContext) *mux.Router {
	r := mux.NewRouter()

	r.Handle("/", ContextHandler{c, HomeHandler})
	r.Handle("/global.js", ContextHandler{c, GlobalJSHandler})
	r.Handle("/api/books", ContextHandler{c, ListBooksHandler}).Methods("OPTIONS", "GET")
	r.Handle("/api/books/{book}/{chapter}", ContextHandler{c, GetChapterHandler}).Methods("OPTIONS", "GET")
	r.Handle("/api/books/{book}", ContextHandler{c, GetBookHandler}).Methods("OPTIONS", "GET")
	r.Handle("/pages/{path:.*}", ContextHandler{c, GetPageHandler}).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(c.StaticDir))))
	return r
}

func HomeHandler(c *AppContext, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hellow world")
	return
}

func GlobalJSHandler(c *AppContext, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "var baseAPIURL=\"http://%s\" ;", *HOST)
	return
}

func ListBooksHandler(c *AppContext, w http.ResponseWriter, r *http.Request) {
	all := c.Lib.GetBooks()
	ulist, err := json.Marshal(all)
	if err != nil {
		log.Println("[ERR] ", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
			fmt.Fprintf(w, "OK")
		} else {
			w.Write(ulist)
		}
	}
}

func GetBookHandler(c *AppContext, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["book"]
	u := c.Lib.GetChapters(uid)

	if u == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
			fmt.Fprintf(w, "OK")
		} else {
			ujson, err := json.Marshal(u)
			if err != nil {
				log.Println("[ERR] marshal:", err)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Write(ujson)
			}
		}
	}
}

func GetChapterHandler(c *AppContext, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookTitle := vars["book"]
	chapterID := vars["chapter"]
	log.Println("Inside chapter handlers:", bookTitle, chapterID)
	u := c.Lib.GetPages(bookTitle, chapterID)

	if u == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
			fmt.Fprintf(w, "OK")
		} else {
			ujson, err := json.Marshal(u)
			if err != nil {
				log.Println("[ERR] marshal:", err)
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Write(ujson)
			}
		}
	}
}

func GetPageHandler(c *AppContext, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pathToImg := vars["path"]
	log.Println("Inside page handlers:", pathToImg)
	ext := strings.ToLower(filepath.Ext(pathToImg))
	fullPath := fmt.Sprintf("%s/%s", c.DataDir, pathToImg)
	p, err := filepath.Abs(fullPath)
	if err != nil {
		log.Println("path err:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if strings.HasPrefix(p, c.DataDir) {
		if ext == ".thumb" {
			baseFilename := strings.TrimSuffix(p, ext)
			log.Println("baseFilename:", baseFilename)
			if _, err := os.Stat(baseFilename + ".png"); !os.IsNotExist(err) {
				// path/to/whatever exists
				log.Println("Serve:", baseFilename+".png")
				http.ServeFile(w, r, baseFilename+".png")
			} else if _, err := os.Stat(baseFilename + ".jpg"); !os.IsNotExist(err) {
				log.Println("Serve:", baseFilename+".jpg")
				http.ServeFile(w, r, baseFilename+".jpg")
			} else if _, err := os.Stat(baseFilename + ".jpeg"); !os.IsNotExist(err) {
				log.Println("Serve:", baseFilename+".jpeg")
				http.ServeFile(w, r, baseFilename+".jpeg")
			} else {
				w.Header().Set("Content-Type", "image/png")
				w.Header().Set("Content-Length", strconv.Itoa(len(EMPTY_IMG_PNG)))
				if _, err := w.Write(EMPTY_IMG_PNG); err != nil {
					log.Println("unable to write empty image.")
				}
			}
		} else {
			if _, err := os.Stat(p); !os.IsNotExist(err) {
				// path/to/whatever exists
				log.Println("Serve:", p)
				http.ServeFile(w, r, p)
			} else {
				w.Header().Set("Content-Type", "image/png")
				w.Header().Set("Content-Length", strconv.Itoa(len(EMPTY_IMG_PNG)))
				if _, err := w.Write(EMPTY_IMG_PNG); err != nil {
					log.Println("unable to write empty image.")
				}
			}
		}
	} else {
		log.Println("path not in:", c.DataDir)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
