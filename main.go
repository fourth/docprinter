package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"flag"

	"github.com/codegangsta/negroni"
	"github.com/goincremental/negroni-oauth2"
	"github.com/goincremental/negroni-sessions"
	"github.com/gorilla/mux"
	"github.com/shurcooL/go/github_flavored_markdown"
)

type config struct {
	clientID     string
	secret       string
	redirectHost string
}

func (c *config) Validate() error {
	c.clientID = strings.TrimSpace(c.clientID)
	c.secret = strings.TrimSpace(c.secret)
	c.redirectHost = strings.TrimSpace(c.redirectHost)
	if c.clientID == "" || c.secret == "" || c.redirectHost == "" {
		return errors.New("All flags must be provided.")
	}

	return nil
}

func main() {
	cfg := &config{}
	flag.StringVar(&cfg.clientID, "clientID", "", "Github client ID")
	flag.StringVar(&cfg.secret, "secret", "", "Github client secret")
	flag.StringVar(&cfg.redirectHost, "host", "localhost", "the host where this is running")

	flag.Parse()

	err := cfg.Validate()

	if err != nil {
		log.Fatal(err)
	}

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
		q := req.URL.Query()
		ref, ok := q["ref"]
		var renderUrl string
		if ok {
			renderUrl = fmt.Sprintf("http://%s/render/%s/path/%s?token=%s&ref=%s", host, repo, path, token, ref[0])
		} else {
			renderUrl = fmt.Sprintf("http://%s/render/%s/path/%s?token=%s", host, repo, path, token)
		}

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

	redirectURL := fmt.Sprintf("http://%s/oauth2callback", cfg.redirectHost)

	n := negroni.New()
	n.Use(sessions.Sessions("my_session", sessions.NewCookieStore([]byte("secret123"))))
	github := oauth2.Github(&oauth2.Options{
		ClientID:     cfg.clientID,
		ClientSecret: cfg.secret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"repo"},
	})
	n.Use(github)

	router := mux.NewRouter()

	//routes added to mux do not require authentication
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		token := oauth2.GetToken(req)
		if token == nil || token.IsExpired() {
			fmt.Fprintf(w, "not logged in, or the access token is expired")
			return
		}
		fmt.Fprintf(w, "logged in")
		return
	})

	assetSrv := negroni.New(&negroni.Static{
		Dir:    http.Dir("assets"),
		Prefix: "/assets",
	})

	assetSrv.UseHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not Found")
	}))

	router.Handle("/assets/{assetpath}", assetSrv)

	router.HandleFunc("/render/{repo:.*}/path/{path:.*}", func(w http.ResponseWriter, req *http.Request) {
		log.Println("rendering the pdf")
		vars := mux.Vars(req)
		repo := vars["repo"]
		path := vars["path"]
		query := req.URL.Query()
		token := query.Get("token")
		ref, refpresent := query["ref"]
		urlStr := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo, path)
		if refpresent {
			urlStr += fmt.Sprintf("?ref=%s", ref[0])
		}
		log.Println(urlStr)
		githubreq, _ := http.NewRequest("GET", urlStr, nil)
		githubreq.Header.Add("Authorization", "Bearer "+token)
		githubreq.Header.Add("Accept", "application/vnd.github.v3.raw")

		res, err := client.Do(githubreq)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		defer res.Body.Close()

		var md bytes.Buffer

		log.Println("outputting")
		io.Copy(&md, res.Body)

		io.WriteString(w, `<html><head><meta charset="utf-8">
<link href='http://fonts.googleapis.com/css?family=Merriweather' rel='stylesheet' type='text/css'>
<link href="/assets/print.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)

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
