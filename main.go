package main

import (
	"github.com/gorilla/websocket"
	"github.com/ngaut/log"
	"net/http"
	"os"
	"text/template"
	"bufio"
	"io"

)

var upgrader = websocket.Upgrader{}

func home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("index.html").ParseFiles("./html/index.html")
	t.Execute(w, "ws://"+r.Host+"/echo")
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}
	defer c.Close()

	f, _ := os.Open("./xxx.txt")
	defer f.Close()

	bi:=bufio.NewReader(f)
	for {
		l,err:=bi.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		c.WriteMessage(websocket.TextMessage, l)
		log.Info(l)
		//time.Sleep(time.Millisecond*64)
	}

	return
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			break
		}
		f.WriteString(string(message) +"\n")
		log.Infof("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Error("write:", err)
			break
		}
	}
}

func StaticServer(w http.ResponseWriter, req *http.Request) {
	log.Info(111)
	staticHandler := http.FileServer(http.Dir("./public/"))
	staticHandler.ServeHTTP(w, req)
	return
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/echo", echo)
	http.Handle("/img/", http.FileServer(http.Dir("./public/")))
	//http.HandleFunc("/public", StaticServer)
	http.ListenAndServe(":8088", nil)
}
