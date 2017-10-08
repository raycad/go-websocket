package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/mem"
)

// MaxGoroutines -- define the max goroutines number
const MaxGoroutines int64 = 1e4

// RecvMsgTimer -- timer interval for counting received messages (in second)
const RecvMsgTimer int = 10

// GoRoutinesPool -- goroutines pool
var GoRoutinesPool = make(chan struct{}, MaxGoroutines)

// Connected clients
var clients = make(map[*websocket.Conn]bool)

// Broadcast channel
var broadcast = make(chan Message)

// TestBroadcastMsg -- for testing broadcasting messages to clients
const TestBroadcastMsg bool = false

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	// EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var mu sync.Mutex

// Message -- define our message object
type Message struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Result   string `json:"result"`
}

// Stats -- Server statistics information
type Stats struct {
	recvMsgPrint int // Received messages printed in an interval time
	recvMsgCount int // Received messages in an interval time
	elapsedTime  int64
	startTime    time.Time
}

var stats Stats

func main() {
	// Create a simple file server
	// fs := http.FileServer(http.Dir("../public"))
	// http.Handle("/", fs)

	http.HandleFunc("/stats", handleWSStats)

	// Configure websocket route
	http.HandleFunc("/aep", handleWSConnections)

	if TestBroadcastMsg == true {
		// Start listening for incoming messages
		go handleWSMessages()
	}

	// Execute server statistics
	go handleServerStats()

	// Start the server on localhost port 2706 and log any errors
	log.Println("http server started on :2706")
	err := http.ListenAndServe(":2706", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Get server statistics information
func handleWSStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	v, _ := mem.VirtualMemory()

	calculateRecvMsg()

	stats := fmt.Sprintf("Websocket clients: %d<br>Goroutines number: %d<br>"+
		"Received messages in %d seconds: %d<br>"+
		"Total VM: %v<br>Free VM:%v<br>UsedPercent VM: %v<br>",
		len(clients), runtime.NumGoroutine(), stats.elapsedTime,
		stats.recvMsgPrint, v.Total, v.Free, v.UsedPercent)

	io.WriteString(w, stats)
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
			return
		}

		mu.Lock()
		stats.recvMsgCount++
		mu.Unlock()

		fmt.Printf("Received message: %s in GoroutineId %d, Goroutine number %d\n",
			msg, getGID(), runtime.NumGoroutine())

		// Would block if GoRoutinesPool channel is already filled
		GoRoutinesPool <- struct{}{}

		go parseMessage(msg)

		time.Sleep(10 * time.Millisecond)
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

func parseMessage(msg Message) {
	// Update message
	msg.Result = strconv.Itoa(heavyComputation())

	if TestBroadcastMsg == true {
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}

	// Update pool
	<-GoRoutinesPool
}

func calculateRecvMsg() {
	mu.Lock()
	stats.recvMsgPrint = stats.recvMsgCount
	stats.recvMsgCount = 0
	mu.Unlock()
	stats.elapsedTime = int64(time.Since(stats.startTime) / time.Second)
	// Restart the startTime
	stats.startTime = time.Now()
}

func handleServerStats() {
	for {
		time.Sleep(time.Duration(RecvMsgTimer) * time.Second)
		calculateRecvMsg()
		fmt.Printf("handleServerStats in GoroutineId %d. Received %d messages"+
			" in %d seconds\n", getGID(), stats.recvMsgPrint, stats.elapsedTime)
	}
}
