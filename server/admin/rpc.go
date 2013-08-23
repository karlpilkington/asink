package main

import (
	"net"
	"net/rpc"
	"log"
	"syscall"
)

func RPCCall(method string, args interface{}, reply interface{}) error {
	socket := "/tmp/asink.sock"
	client, err := rpc.DialHTTP("unix", socket)
	if err != nil {
		if err2, ok := err.(*net.OpError); ok {
			if err2.Err == syscall.ENOENT {
				log.Fatal("The socket ("+socket+") was not found")
			} else if err2.Err == syscall.ECONNREFUSED {
				log.Fatal("A connection was refused to "+socket+". Please check the permissions and ensure the server is running.")
			}
		}
		return err
	}
	defer client.Close()

	err = client.Call(method, args, reply)
	return err
}
