package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	// EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message Define our message object
type Message struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Result   string `json:"result"`
}

func main() {
	// Create a simple file server
	// fs := http.FileServer(http.Dir("../public"))
	// http.Handle("/", fs)

	http.HandleFunc("/", serveHome)

	// Configure websocket route
	http.HandleFunc("/aep", handleWSConnections)

	// Start listening for incoming messages
	go handleWSMessages()

	// Start the server on localhost port 2706 and log any errors
	log.Println("http server started on :2706")
	err := http.ListenAndServe(":2706", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found.", 404)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, "<html><body>Hi! I'm Raycad!</body></html>")
}

// Handle Web Socket Connection
func handleWSConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	fmt.Printf("Client %d\n", len(clients))

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			// Break or return the current Goroutine
			break
		} else {
			fmt.Printf("Received message: %s in GoroutineId %d, Goroutine number %d\n",
				msg, getGID(), runtime.NumGoroutine())
		}

		// Update message
		msg.Result = strconv.Itoa(heavyComputation())

		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleWSMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			} else {
				fmt.Printf("Sent message in GoroutineId %d, Goroutine number %d\n",
					getGID(), runtime.NumGoroutine())
			}
		}
	}
}

// Get current Goroutine Id
func getGID() uint64 {
	startTime := time.Now()
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	elapsedTime := time.Since(startTime)
	fmt.Printf("getGID took %s\n", elapsedTime)

	return n
}

// Simulate long computation
func heavyComputation() int {
	// Set the size larger to make longer computation to test performance
	const min = 100
	size := rand.Intn(min) + min
	var result int

	for i := 0; i < size; i++ {
		for k := 0; k < size; k++ {
			result++
		}
	}

	return result
}
