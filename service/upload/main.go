package main

import (
	cfg "github.com/lilhammer111/hammer-cloud/config"
	"github.com/lilhammer111/hammer-cloud/route"
	"log"
)

func main() {

	router := route.Router()
	log.Fatal(router.Run(cfg.UploadServiceHost))

	//fs := http.FileServer(http.Dir("static/"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))

}
