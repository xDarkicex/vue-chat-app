package main

import (
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
)

type Collections struct {
	Users    []User
	Messages []Message
}

type User struct {
	Nickname string
}

type Message struct {
	Text string
}

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		user := &User{}

		log.Println("Someone Connected")
		so.Join("chat_room")
		so.Emit("users")
		so.On("message", func(text string) {
			data := struct {
				Message string
				User    string
			}{
				Message: text,
				User:    user.Nickname,
			}
			log.Println("emit: ", so.Emit("message", data))
			so.BroadcastTo("chat_room", "message", data)
		})
		so.On("nickname", func(nickname string) {
			log.Println("A user changed their nickname to ", nickname)
			user.Nickname = nickname
		})
	})
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./public")))
	log.Println("serving on :3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil))
}
