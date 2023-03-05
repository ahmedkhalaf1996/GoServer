package notification

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// Room represents a single chat room
type Room struct {
	forward chan []byte

	join chan *Chatter

	leave chan *Chatter

	chatters map[*Chatter]bool

	topic string
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("serving http failed ", err)
		return
	}

	chatter := &Chatter{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- chatter
	defer func() {
		r.leave <- chatter
	}()
	go chatter.write()
	chatter.read()
}

// NewRoom creates a new chat room
func NewRoom(topic string) *Room {
	return &Room{
		forward:  make(chan []byte),
		join:     make(chan *Chatter),
		leave:    make(chan *Chatter),
		chatters: make(map[*Chatter]bool),
		topic:    topic,
	}
}

// run group initilize a group chat room
func (r *Room) RunGroup() {
	log.Printf("running chat room %v", r.topic)
	for {
		select {
		case chatter := <-r.join:
			log.Printf("new chatter in room %v", r.topic)
			// calling data

			r.chatters[chatter] = true
			// GetMessFromDB(r.topic, r)

		case chatter := <-r.leave:
			log.Printf("chatter leaving room %v", r.topic)
			delete(r.chatters, chatter)
			close(chatter.send)
		case msg := <-r.forward:
			// data := FromJSON(msg)
			// fmt.Println("EEEERRR", data)
			// if data != nil {
			// 	if data.Sender != "" && data.MessageType != nil && data.Message != "" {
			// 		// SendMessToGroupDB(data, r)
			// 	}
			// 	log.Printf("chatter '%v' writing message to room %v, message: %v", data.Sender, r.topic, data.Message)
			// }

			for chatter := range r.chatters {
				select {
				case chatter.send <- msg:
				default:
					delete(r.chatters, chatter)
					close(chatter.send)
				}
			}
		}
	}
}

// Run initializes a chat room
func (r *Room) Run() {
	log.Printf("running chat room %v", r.topic)
	for {
		select {
		case chatter := <-r.join:
			log.Printf("new chatter in room %v", r.topic)
			// calling data

			r.chatters[chatter] = true
			// GetMessFromDB(r.topic, r)

		case chatter := <-r.leave:
			log.Printf("chatter leaving room %v", r.topic)
			delete(r.chatters, chatter)
			close(chatter.send)
		case msg := <-r.forward:
			// data := FromJSON(msg)
			// // fmt.Println("EEEERRR", data)
			// if data != nil {
			// 	if data.Sender != "" && data.Received != "" && data.MessageType != nil && data.Message != "" {
			// 		// SendMessToDB(data, r)
			// 	}
			// 	log.Printf("chatter '%v' writing message to room %v, message: %v", data.Sender, r.topic, data.Message)
			// }

			for chatter := range r.chatters {
				select {
				case chatter.send <- msg:
				default:
					delete(r.chatters, chatter)
					close(chatter.send)
				}
			}
		}
	}
}
