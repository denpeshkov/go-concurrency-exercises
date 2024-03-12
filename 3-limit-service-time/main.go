//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

var m sync.Map

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	var t <-chan time.Time
	done := make(chan struct{})

	if !u.IsPremium {
		v, _ := m.LoadOrStore(u.ID, time.After(10*time.Second))
		t = v.(<-chan time.Time)
	}

	go func() {
		defer close(done)
		process()
	}()

	select {
	case <-t:
		return false
	case <-done:
	}
	return true
}

func main() {
	RunMockServer()
}
