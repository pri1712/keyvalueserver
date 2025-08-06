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
			}
		} else {
			//log.Printf("Lock acquired for key: %v, value: %v, version: %v", lk.key, value, version)
			if value == "unlocked" {
				putErr := lk.ck.Put(lk.key, lk.clientID, version)
				if putErr == rpc.OK {
					log.Printf("Lock acquired for key which was existing %s", lk.key)
					return
				}
			} else {
				log.Printf("First unlock it for key %s", lk.key)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (lk *Lock) Release() {
	value, version, err := lk.ck.Get(lk.key)
	//log.Printf("error %v", err)
	//log.Printf("Got lock value: %v, version: %v", value, version)
	if err != rpc.OK {
		log.Printf("Got lock error: %v", err)
	} else if value == lk.clientID {
		log.Printf("Release lock")
		lk.ck.Put(lk.key, "unlocked", version)
	} else {
		log.Printf("Its not your lock to release it.")
	}
	return
	//lk.ck.Put(lk.key, "", version)
}
