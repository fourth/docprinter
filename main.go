package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-oauth2"
	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
	"github.com/shurcooL/go/github_flavored_markdown"
)

func main() {
	secureMux := mux.NewRouter()
	client := &http.Client{}

	secureMux.HandleFunc("/print/{repo:.*}/path/{path:.*}", func(w http.ResponseWriter, req *http.Request) {
		log.Println("Making the pdf")
		token := oauth2.GetToken(req).Access()
		log.Println(req.URL)
		vars := mux.Vars(req)
		repo := vars["repo"]
		path := vars["path"]
		host := req.Host
		renderUrl := fmt.Sprintf("http://%s/render/%s/path/%s?token=%s", host, repo, path, token)
		cmd := exec.Command("phantomjs", "./renderpdf.js", renderUrl)
		log.Println(renderUrl)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("starting the command")
		if err = cmd.Start(); err != nil {
			log.Fatal(err)
		}
		log.Println("outputting")
		io.Copy(w, stdout)

		log.Println("waiting")
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
		log.Println("done")
	})

	secure := negroni.New()
	secure.Use(oauth2.LoginRequired())
	secure.UseHandler(secureMux)

	n := negroni.New()
	n.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("secret123"))))
	github := oauth2.Github(&oauth2.Options{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://docprint.sjjdev.com/oauth2callback",
		Scopes:       []string{"repo"},
	})
	n.Use(github)

	router := mux.NewRouter()

	//routes added to mux do not require authentication
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		token := oauth2.GetToken(req)
		s := sessions.GetSession(req)
		fmt.Fprintln(w, "token", s.Get("oauth2_token"))
		if token == nil || token.IsExpired() {
			fmt.Fprintf(w, "not logged in, or the access token is expired")
			return
		}
		fmt.Fprintf(w, "logged in")
		return
	})

	router.HandleFunc("/render/{repo:.*}/path/{path:.*}", func(w http.ResponseWriter, req *http.Request) {
		log.Println("rendering the pdf")
		vars := mux.Vars(req)
		repo := vars["repo"]
		path := vars["path"]
		query := req.URL.Query()
		token := query.Get("token")
		urlStr := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo, path)
		log.Println(urlStr)
		githubreq, _ := http.NewRequest("GET", urlStr, nil)
		githubreq.Header.Add("Authorization", "Bearer "+token)
		githubreq.Header.Add("Accept", "application/vnd.github.v3.raw")
		res, err := client.Do(githubreq)
		if err != nil {
			fmt.Fprintln(w, err)
		}
		defer res.Body.Close()

		var md bytes.Buffer

		log.Println("outputting")
		io.Copy(&md, res.Body)

		io.WriteString(w, `<html><head><meta charset="utf-8"><link href="https://github.com/assets/github.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)

		w.Write(github_flavored_markdown.Markdown(md.Bytes()))

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		io.WriteString(w, `</article></body></html>`)

		return
	})

	router.Handle("/{wildcard:.*}", secure)

	n.UseHandler(router)

	n.Run(":3000")
}
