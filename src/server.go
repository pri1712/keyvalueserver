package kvsrv

import (
	"log"
	"sync"

	"kvserver/src/labrpc"
	"kvserver/src/rpc"
	"kvserver/src/tester1"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type ValueTuple struct {
	Value   string
	Version rpc.Tversion
}
type KVServer struct {
	mu          sync.Mutex
	keyValueMap map[string]ValueTuple
	// Your definitions here.
}

func MakeKVServer() *KVServer {
	kv := &KVServer{keyValueMap: make(map[string]ValueTuple)}
	// Your code here.
	return kv
}

// Get returns the value and version for args.Key, if args.Key
// exists. Otherwise, Get returns ErrNoKey.
func (kv *KVServer) Get(args *rpc.GetArgs, reply *rpc.GetReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	valueTuple, exists := kv.keyValueMap[args.Key]
	if !exists {
		reply.Err = rpc.ErrNoKey
	} else {
		reply.Value = valueTuple.Value
		reply.Version = valueTuple.Version
		reply.Err = rpc.OK
	}
}

// Update the value for a key if args.Version matches the version of
// the key on the server. If versions don't match, return ErrVersion.
// If the key doesn't exist, Put installs the value if the
// args.Version is 0, and returns ErrNoKey otherwise.
func (kv *KVServer) Put(args *rpc.PutArgs, reply *rpc.PutReply) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	valueTuple, exists := kv.keyValueMap[args.Key]
	if !exists {
		if args.Version == 0 {
			kv.keyValueMap[args.Key] = ValueTuple{args.Value, 1}
			reply.Err = rpc.OK
		} else {
			reply.Err = rpc.ErrNoKey
		}
	} else if args.Version != valueTuple.Version {
		reply.Err = rpc.ErrVersion
	} else {
		kv.keyValueMap[args.Key] = ValueTuple{args.Value, args.Version + 1}
		reply.Err = rpc.OK
	}
}

// You can ignore Kill() for this lab
func (kv *KVServer) Kill() {
}

// You can ignore all arguments; they are for replicated KVservers
func StartKVServer(ends []*labrpc.ClientEnd, gid tester.Tgid, srv int, persister *tester.Persister) []tester.IService {
	kv := MakeKVServer()
	return []tester.IService{kv}
}
