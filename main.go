package main

import (
	"os"
	"log"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"os/exec"
	"io"
	"time"
)

func main() {
	confJson, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error:", err)
	}
	type Config struct {
		ListenHost   string
		ListenPort   int
		Tokens       []string
	}
	var config Config

	err = json.Unmarshal(confJson, &config)
	var srv http.Server
	srv.Addr = fmt.Sprintf("%s:%d", config.ListenHost, config.ListenPort)

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			authentication := req.Header.Get("Authentication")
			authStr := strings.SplitN(authentication, " ", 2)
			if strings.ToLower(authStr[0]) != "bearer" {
				rw.WriteHeader(400)
				rw.Write([]byte("Undefined Authentication Type"))
				return
			}
			found := false
			for _, token := range config.Tokens {
				if authStr[1] == token {
					found = true
					break
				}
			}
			if found == false {
				rw.WriteHeader(400)
				rw.Write([]byte("Invalid token"))
				return
			}
			err := req.ParseMultipartForm(2 * 1024)
			if err != nil {
				rw.WriteHeader(400)
				rw.Write([]byte(err.Error()))
				return
			}
			commandStr := req.PostForm.Get("command")
			reader := strings.NewReader("Queued\r\n")
			io.Copy(rw, reader)
			oc := make(chan []byte)
			done := 0
			go func() {
				cmd := exec.Command("bash", "-c", commandStr)
				log.Println("Command:", cmd.String())
				op, err := cmd.CombinedOutput()
				if err != nil {
					log.Println("Error:", err)
				}
				oc <- op
				close(oc)
			}()
			go func() {
				// Keep the connection alive. Is this required?
				for {
					io.Copy(rw, strings.NewReader(""))
					time.Sleep(30 * time.Second)
					if done == 1 {
						break
					}
				}
			}()
			output := <- oc
			rw.Write(output)
			done = 1
		} else if req.Method == "GET" {
			html, err := os.ReadFile("index.html")
			if err != nil {
				log.Fatal("HTML file is missing")
			}
			rw.Write(html)
		}
	})
	fmt.Println("Starting server")
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Server done")
}
