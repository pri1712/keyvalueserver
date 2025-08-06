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
	ck       kvtest.IKVClerk
	key      string
	clientID string
	// You may add code here
}

// The tester calls MakeLock() and passes in a k/v clerk; your code can
// perform a Put or Get by calling lk.ck.Put() or lk.ck.Get().
//
// Use l as the key to store the "lock state" (you would have to decide
// precisely what the lock state is).
func MakeLock(ck kvtest.IKVClerk, l string) *Lock {
	lk := &Lock{ck: ck, key: l, clientID: kvtest.RandValue(8)}
	// You may add code here
	return lk
}

func (lk *Lock) Acquire() {
	for {
		value, version, err := lk.ck.Get(lk.key)
		if err == rpc.ErrNoKey {
			putErr := lk.ck.Put(lk.key, lk.clientID, 0)
			if putErr == rpc.OK {
				log.Printf("Lock acquired for key %s", lk.key)
				return
			} else {
				log.Printf("Failed to acquire lock for key %s", lk.key)
				time.Sleep(1 * time.Second)
				continue
			}
		} else if err != rpc.OK {
			log.Printf("Get error for key %s", lk.key)
			time.Sleep(1 * time.Second)
			continue
		} else {
			log.Printf("Lock acquired for key: %v, value: %v, version: %v", lk.key, value, version)
			if value == "unlocked" {
				putErr := lk.ck.Put(lk.key, lk.clientID, version)
				if putErr == rpc.OK {
					log.Printf("Lock acquired for key which was existing %s", lk.key)
					return
				} else {
					log.Printf("Failed to acquire lock for key which was existing %s", lk.key)
					time.Sleep(1 * time.Second)
					continue
				}
			} else if value == lk.clientID {
				log.Printf("This is our lock")
				return
			} else {
				log.Printf("lock held by other client...waiting")
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}
}

func (lk *Lock) Release() {

	for {
		value, version, err := lk.ck.Get(lk.key)
		if err == rpc.ErrNoKey {
			log.Printf("This key does not exist")
			return
		} else if err != rpc.OK {
			log.Printf("Got error during release: %v", err)
			time.Sleep(1 * time.Second)
			continue
		} else {
			if value == lk.clientID {
				log.Printf("Release lock")
				putErr := lk.ck.Put(lk.key, "unlocked", version)
				if putErr == rpc.OK {
					log.Printf("Lock released for key %s", lk.key)
				} else {
					log.Printf("Failed to release lock for key %s", lk.key)
					time.Sleep(1 * time.Second)
					continue
				}
			} else {
				log.Printf("Cannot release lock owned by other client")
				return
			}
		}
	}
	//lk.ck.Put(lk.key, "", version)
}
