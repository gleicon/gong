// Copyright 2014 gong authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	html "html/template"
	text "text/template"

	"github.com/apexskier/httpauth"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var (
	VERSION = "tip"
	APPNAME = "gong"

	// Templates
	HTML *html.Template
	TEXT *text.Template

	backend     httpauth.GobFileAuthBackend
	aaa         httpauth.Authorizer
	roles       map[string]httpauth.Role
	backendfile = "auth.gob"
)

func createDefaultUsers() {
}

func main() {
	configFile := flag.String("c", "gong.conf", "")
	flag.Usage = func() {
		fmt.Println("Usage: gong [-c gong.conf] [-l logfile]")
		os.Exit(1)
	}
	flag.Parse()

	var err error
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Parse templates.
	HTML = html.Must(html.ParseGlob(config.TemplatesDir + "/*.html"))
	TEXT = text.Must(text.ParseGlob(config.TemplatesDir + "/*.txt"))

	// Set up databases.
	db, err := sql.Open("mysql", config.DB.MySQL)
	if err != nil {
		log.Fatal(err)
	}

	// Set GOMAXPROCS and show server info.
	var cpuinfo string
	if n := runtime.NumCPU(); n > 1 {
		runtime.GOMAXPROCS(n)
		cpuinfo = fmt.Sprintf("%d CPUs", n)
	} else {
		cpuinfo = "1 CPU"
	}
	log.Printf("%s %s (%s)", APPNAME, VERSION, cpuinfo)
	// Initialize auth
	os.Create(backendfile)
	defer os.Remove(backendfile)

	// create the backend
	backend, err = httpauth.NewGobFileAuthBackend(backendfile)
	if err != nil {
		log.Println(err)
	}

	// create some default roles
	roles = make(map[string]httpauth.Role)
	roles["user"] = 30
	roles["admin"] = 80
	aaa, err = httpauth.NewAuthorizer(backend, []byte("cookie-encryption-key"), "user", roles)

	// create a default user
	hash, err := bcrypt.GenerateFromPassword([]byte("adminadmin"), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	defaultUser := httpauth.UserData{Username: "admin", Email: "admin@localhost", Hash: hash, Role: "admin"}
	err = backend.SaveUser(defaultUser)
	if err != nil {
		log.Println(err)
	}

	// Start HTTP server.
	s := new(httpServer)
	s.init(config, db)
	go s.ListenAndServe()
	go s.ListenAndServeTLS()

	// Sleep forever.
	select {}
}
