package lock

import (
	"kvserver/src/kvtest1"
	"kvserver/src/rpc"
	"log"
	"time"
)

type Lock struct {
	// IKVClerk is a go interface for k/v clerks: the interface hides
	// the specific Clerk type of ck but promises that ck supports
	// Put and Get.  The tester passes the clerk in when calling
	// MakeLock().
	ck  kvtest.IKVClerk
	key string
	// You may add code here
}

// The tester calls MakeLock() and passes in a k/v clerk; your code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
//
// Use l as the key to store the "lock state" (you would have to decide
// precisely what the lock state is).
func MakeLock(ck kvtest.IKVClerk, l string) *Lock {
	lk := &Lock{ck: ck, key: l}
	// You may add code here
	return lk
}

func (lk *Lock) Acquire() {
	for {
		err := lk.ck.Put(lk.key, "locked", 0)
		if err == rpc.OK {
			log.Printf("Acquired lock")
			return
		} else {
			time.Sleep(100 * time.Millisecond)
			log.Printf("Sleeping for 100ms")
		}
	}
}

func (lk *Lock) Release() {
	value, version, err := lk.ck.Get(lk.key)
	log.Printf("Got lock value: %v, version: %v", value, version)
	if err != rpc.OK {
		log.Printf("Got lock error: %v", err)
		return
	}
	//lk.ck.Put(lk.key, "", version)
}
