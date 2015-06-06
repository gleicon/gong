// Copyright 2014 gong authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (s *httpServer) route() {
	// Assets
	http.Handle("/js/", http.FileServer(http.Dir(s.config.DocumentRoot)))
	http.Handle("/css/", http.FileServer(http.Dir(s.config.DocumentRoot)))
	http.Handle("/img/", http.FileServer(http.Dir(s.config.DocumentRoot)))
	http.Handle("/fonts/", http.FileServer(http.Dir(s.config.DocumentRoot)))

	// Static routes
	http.HandleFunc("/login", s.loginHandler)
	http.HandleFunc("/logout", s.logoutHandler)
	http.HandleFunc("/signup", s.serveFileByURI)
	http.HandleFunc("/dashboard", s.serveFileByURI)
	http.HandleFunc("/index.html", s.indexHandler)
	http.HandleFunc("/", s.indexHandler)

	// API handlers
	http.HandleFunc("/api/v1/login", s.loginHandler)
	http.HandleFunc("/api/v1/logout", s.logoutHandler)
	http.HandleFunc("/api/v1/signup", s.signUpHandler)
}

func (s *httpServer) serveFileByURI(w http.ResponseWriter, r *http.Request) {
	file := fmt.Sprintf("%s/%s.html", s.config.DocumentRoot, r.URL.Path[1:])
	http.ServeFile(w, r, file)
}

func (s *httpServer) indexHandler(w http.ResponseWriter, r *http.Request) {

	err := aaa.Authorize(w, r, true)
	if err != nil && err.Error() != "already authenticated" {
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user, err := aaa.CurrentUser(w, r); err == nil {
		log.Println(user.Username)
		file := s.config.DocumentRoot + "index.html"
		http.ServeFile(w, r, file)
	} else {
		log.Println(err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func (s *httpServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		err := aaa.Login(w, r, username, password, "/")
		if err != nil && strings.Contains(err.Error(), "already authenticated") {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
		break
	case "GET":
		file := fmt.Sprintf("%s/login.html", s.config.DocumentRoot)
		http.ServeFile(w, r, file)
		break
	default:
		http.Error(w, "Method not allowed", 405)
		break
	}
}

func (s *httpServer) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := aaa.Logout(w, r); err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (s *httpServer) signUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Fprintf(w, "signup\r\n")
	} else {
		http.Error(w, "Method not allowed", 405)
	}
}
