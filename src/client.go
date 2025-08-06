package kvsrv

import (
	"kvserver/src/kvtest1"
	"kvserver/src/rpc"
	"kvserver/src/tester1"
)

var tries int = 2

type Clerk struct {
	clnt   *tester.Clnt
	server string
}

func MakeClerk(clnt *tester.Clnt, server string) kvtest.IKVClerk {
	ck := &Clerk{clnt: clnt, server: server}
	// You may add code here.
	return ck
}

type PutClientArgs struct {
	ClientKey         string
	ClientDataValue   string
	ClientDataVersion rpc.Tversion
}

type PutClientReply struct {
	ClientRpcError rpc.Err
}

type GetClientArgs struct {
	ClientKey string
}

type GetClientReply struct {
	ClientKey          string
	ClientRpcError     rpc.Err
	ClientValueVersion rpc.Tversion
}

// Get fetches the current value and version for a key.  It returns
// ErrNoKey if the key does not exist. It keeps trying forever in the
// face of all other errors.
//
// You can send an RPC with code like this:
// ok := ck.clnt.Call(ck.server, "KVServer.Get", &args, &reply)
//
// The types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. Additionally, reply must be passed as a pointer.
func (ck *Clerk) Get(key string) (string, rpc.Tversion, rpc.Err) {
	// You will have to modify this function.
	for {
		request := &rpc.GetArgs{Key: key}
		reply := &rpc.GetReply{}
		ok := ck.clnt.Call(ck.server, "KVServer.Get", request, reply)
		if ok {
			return reply.Value, reply.Version, reply.Err
		}
	}

}

// Put updates key with value only if the version in the
// request matches the version of the key at the server.  If the
// versions numbers don't match, the server should return
// ErrVersion.  If Put receives an ErrVersion on its first RPC, Put
// should return ErrVersion, since the Put was definitely not
// performed at the server. If the server returns ErrVersion on a
// resend RPC, then Put must return ErrMaybe to the application, since
// its earlier RPC might have been processed by the server successfully
// but the response was lost, and the Clerk doesn't know if
// the Put was performed or not.
//
// You can send an RPC with code like this:
// ok := ck.clnt.Call(ck.server, "KVServer.Put", &args, &reply)
//
// The types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. Additionally, reply must be passed as a pointer.
func (ck *Clerk) Put(key string, value string, version rpc.Tversion) rpc.Err {
	// You will have to modify this function.
	retriedRequest := false
	for {
		request := &rpc.PutArgs{Key: key, Value: value, Version: version}
		reply := &rpc.PutReply{}
		ok := ck.clnt.Call(ck.server, "KVServer.Put", request, reply)
		if ok {
			switch reply.Err {
			case rpc.ErrNoKey:
				return rpc.ErrNoKey
			case rpc.ErrVersion:
				if retriedRequest {
					return rpc.ErrMaybe
				} else {
					return rpc.ErrVersion
				}
			default:
				return rpc.OK
			}
		} else {
			retriedRequest = true
		}
	}

}
