package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *server {
	return &server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *server) handleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client to order book feed:", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second *2)
	}
}

func (s *server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client: ", ws.RemoteAddr())
	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error: ", err)
			continue
		}
		msg := buf[:n]
		// fmt.Println(string(msg))
		// ws.Write([]byte("Thank you for the message!!!"))
		s.broadcast(msg)
	}
}

func (s *server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}
		} (ws)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/orderbookfeed", websocket.Handler((server.handleWSOrderbook)))
	http.ListenAndServe(":3000", nil)

}