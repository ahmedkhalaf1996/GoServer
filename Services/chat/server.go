package chat

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/exp/slices"

	// "main/database"

	"net/http"
)

var Rooms []string
var GroupRooms []string

// !!!!!!!!!!!!!! NOTE !!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!! we need to get existing room id's from db !!!!!!!!!!!!!1
func GetPrivateRoomeData(router *mux.Router) {
	// resp, err := http.Get("http://localhost:5000/Private/GetAllPrivateRoomsID")
	// if err != nil {
	// 	log.Fatalln(err)
	// } else {
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}

	// 	var d struct {
	// 		Data []string `json:"data" bson:"data"`
	// 	}

	// 	json.Unmarshal(body, &d)
	// 	for _, element := range d.Data {
	// 		fmt.Printf("element: %v\n", element)
	// 		AddNewRoom(element, router)
	// 	}
	// 	// }
	// }

}

//// group
func GetGroupRoomData(router *mux.Router) {
	// resp, err := http.Get("http://localhost:5000/Gorup/GetAllGroupRoomsID")
	// if err != nil {
	// 	log.Fatalln(err)
	// } else {
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}

	// 	var d struct {
	// 		Data []string `json:"data" bson:"data"`
	// 	}

	// 	json.Unmarshal(body, &d)
	// 	for _, element := range d.Data {
	// 		fmt.Printf("G element: %v\n", element)
	// 		AddNewGroupRoom(element, router)
	// 	}
	// }

}

// !!!!!!!!!!!!!! NOTE !!!!!!!!!!!!!!!!!!!

// Run starts a new chat server with 4 chat rooms, listening on port 8080
func AddNewRoom(name string, router *mux.Router) {
	// fmt.Println("is Exist Or not ", slices.Contains(Rooms, name))
	if slices.Contains(Rooms, name) {
		fmt.Println("Room Already Exist!", name)
	} else {
		fmt.Println("new Room", name)
		Rooms = append(Rooms, name)
		r := NewRoom(name)
		router.Handle("/chat/"+name, r)
		go r.Run()
	}

}

func AddNewGroupRoom(name string, router *mux.Router) {
	// fmt.Println("is Exist Or not ", slices.Contains(Rooms, name))
	if slices.Contains(Rooms, name) {
		fmt.Println("Room Group Already Exist!", name)
	} else {
		fmt.Println("new Group Room", name)
		Rooms = append(Rooms, name)
		r := NewGroupRoom(name)
		router.Handle("/chatGroup/"+name, r)
		go r.RunGroup()
	}

}

func Run() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                   // All origins
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH"}, // Allowing only get, just an example
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"}, // All origins
	})

	router := mux.NewRouter()

	// fs := http.FileServer(http.Dir("socketDocs"))
	// router.Handle("/", fs)

	router.HandleFunc("/.well-known/pki-validation/87EBE7B7073FA2663E5AC42FCF3AFC9D.txt", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("3519D7C18CE919D0E47F16E4F6AEA3216A13BB4BE13DFDBD7602D02F1BB991AC\ncomodoca.com\n90bd8af47dbef34"))
		return
	})

	router.HandleFunc("/AddChat", func(w http.ResponseWriter, r *http.Request) {

		ConvID := r.URL.Query().Get("ConvID")
		AddNewRoom(ConvID, router)

	})

	router.HandleFunc("/AddGroup", func(w http.ResponseWriter, r *http.Request) {

		ConvID := r.URL.Query().Get("ConvID")
		AddNewGroupRoom(ConvID, router)

	})

	// ---
	GetPrivateRoomeData(router)
	GetGroupRoomData(router)
	// ---
	fmt.Println("host in chat ", os.Getenv("MAINHOST"))
	fmt.Println("Server is ready and is listening at port :8080 . . .")
	http.ListenAndServe(":8080", c.Handler(router))
	// http.ListenAndServeTLS(":8080", "certificate.crt", "private.key", c.Handler(router))
}

// !! important TEst !!//

// go https://websocketking.com/

// connection ws://localhost:8080/chat/arduino
// 44.201.87.128
// payload
// {
// 	"message": "hi ok",
// 	"sender": "anumesous",
// 	"received": ""
//   }
// end
