package notification

import (
	"os"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gorilla/mux"

	"github.com/rs/cors"
	
	"golang.org/x/exp/slices"
	
	"net/http"
)

var Rooms []string

// !!!!!!!!!!!!!! NOTE !!!!!!!!!!!!!!!!!!!
// !!!!!!!!!!!! we need to get existing room id's from db !!!!!!!!!!!!!1
func GetPrivateRoomeData(router *mux.Router) {
	resp, err := http.Get("http://" +string( os.Getenv("MAINHOST")) + "/Private/GetAllUsers")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var d struct {
		Data []string `json:"data" bson:"data"`
	}

	json.Unmarshal(body, &d)
	for _, element := range d.Data {
		fmt.Printf("element: %v\n", element)
		AddNewRoom(element, router)
	}
	// }

}

// !!!!!!!!!!!!!! NOTE !!!!!!!!!!!!!!!!!!!

// Run starts a new notification server with 4 notification rooms, listening on port 8080
func AddNewRoom(id string, router *mux.Router) {
	// fmt.Println("is Exist Or not ", slices.Contains(Rooms, id))
	if slices.Contains(Rooms, id) {
		fmt.Println("Room Already Exist!", id)
	} else {
		fmt.Println("new Room", id)
		Rooms = append(Rooms, id)
		r := NewRoom(id)
		router.Handle("/notification/"+id, r)
		go r.Run()
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

	fs := http.FileServer(http.Dir("socketDocs"))
	router.Handle("/", fs)

	router.HandleFunc("/listenToNotification", func(w http.ResponseWriter, r *http.Request) {

		UserId := r.URL.Query().Get("UserId")
		AddNewRoom(UserId, router)

	})

	// ---
	GetPrivateRoomeData(router)

	// ---
	fmt.Println("Notification Server is ready and is listening at port :8090 . . .")
	http.ListenAndServe(":8090", c.Handler(router))
	// http.ListenAndServeTLS(":8080", "certificate.crt", "private.key", c.Handler(router))
}
