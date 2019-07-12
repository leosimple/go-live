package limiters

import "log"

type ConnectionLimiter struct {
	concurrentConnection int
	bucket               chan int
}

func NewConnectionLimiter(concurrentConn int) *ConnectionLimiter {
	return &ConnectionLimiter{
		concurrentConnection: concurrentConn,
		bucket:               make(chan int, concurrentConn),
	}
}

func (cl *ConnectionLimiter) GetConnection() bool {
	if len(cl.bucket) > cl.concurrentConnection {
		log.Println("Reached the rate limitation.")
		return false
	}
	cl.bucket <- 1
	log.Println("Successfully got connection.")
	return true
}

func (cl *ConnectionLimiter) FreeConnection() bool {
	c := <-cl.bucket
	log.Printf("One connection is free: %d", c)
	return true
}
