package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

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

var (
	seededGUIDGen     *rand.Rand
	seededGUIDGenOnce sync.Once
	seededGUIDLock    sync.Mutex
	logOneError       sync.Once
)

func genSeededGUID() uint64 {
	// Golang does not seed the rng for us. Make sure it happens.
	seededGUIDGenOnce.Do(func() {
		seededGUIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	})

	// The golang random generators are *not* intrinsically thread-safe.
	seededGUIDLock.Lock()
	defer seededGUIDLock.Unlock()
	return uint64(seededGUIDGen.Int63())
}

func genSeededGUID2() (uint64, uint64) {
	// Golang does not seed the rng for us. Make sure it happens.
	seededGUIDGenOnce.Do(func() {
		seededGUIDGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	})

	seededGUIDLock.Lock()
	defer seededGUIDLock.Unlock()
	return uint64(seededGUIDGen.Int63()), uint64(seededGUIDGen.Int63())
}
