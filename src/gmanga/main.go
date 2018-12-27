package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
)

var (
	HOST     = flag.String("l", "127.0.0.1:8080", "Listening interface, default: 127.0.0.1:8080")
	DATA     = flag.String("d", "-", "Data dir: directory contains manga directory")
	RESOURCE = flag.String("s", "static", "Resource dir: javascript & html resource, default: static dir in current dir")
)

func validateConfig() {
	fi, err := os.Stat(*DATA)
	if os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Fatal(err)
		os.Exit(1)
	}
	if !fi.IsDir() {
		log.Fatal(*DATA, "is not a directory")
		os.Exit(1)
	}
	fi, err = os.Stat(*RESOURCE)
	if os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Fatal(err)
		os.Exit(1)
	}
	if !fi.IsDir() {
		log.Fatal(*RESOURCE, "is not a directory")
		os.Exit(1)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()

	if *DATA == "-" {
		flag.PrintDefaults()
		return
	}
	validateConfig()
	log.Println("Listen from:", *HOST)
	log.Println("Loading files from:", *DATA)
	log.Println("Loading resources from:", *RESOURCE)
	globalContext := NewAppContext()
	globalContext.DataDir = *DATA
	globalContext.StaticDir = *RESOURCE
	globalContext.Lib = NewLibrary(*DATA)
	scanResult := globalContext.Lib.LoadAll()
	log.Println("Scan result:", scanResult)
	//dd := globalContext.Lib.Dump()
	//log.Println("Dump result:\n", dd)
	log.Println("----")
	/// watcher
	done := make(chan bool)
	defer close(done)

	/*
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		go WatchDir(watcher, done, *DATA, globalContext.Items)
		err = watcher.Add(*DATA)
		if err != nil {
			log.Fatal(err)
		}
		//// end watcher
	*/
	//log.Printf("%+v\n", globalContext)

	// start server
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	serverCfg := NewServerCfg(*HOST, globalContext)

	go func() {
		serverCfg.Serve()
	}()
	<-stop
	log.Println("graceful shutting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), 15)
	defer cancel()

	if err := serverCfg.Server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("final down")
}
