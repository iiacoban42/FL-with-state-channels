// Copyright (c) 2019 Chair of Applied Cryptography, Technische Universität
// Darmstadt, Germany. All rights reserved. This file is part of
// perun-eth-demo. Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package client

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	// "github.com/spf13/viper"
	"github.com/spf13/viper"
	"perun.network/go-perun/log"
)

type testAPI struct {
	listener net.Listener
	closed   chan interface{}
}

var api *testAPI

// StartTestAPI sets the package variable `api` to a new `testAPI`
// listening at 0.0.0.0:8080. Should be called after setting up the node.
func StartTestAPI() {
	api = newTestAPI("0.0.0.0:"+ viper.GetString("node.apiport"))
}

func newTestAPI(url string) *testAPI {
	l, err := net.Listen("tcp", url)
	if err != nil {
		log.Panic("TCP listening: ", err.Error())
	}
	log.Info("Listening on: " + url)
	ret := &testAPI{listener: l, closed: make(chan interface{})}
	go ret.startHandling()
	return ret
}

func (a *testAPI) close() {
	close(a.closed)
}

func (a *testAPI) startHandling() {
	log.Trace("Started handling API requests")
	for {
		select {
		case <-a.closed:
			log.Debug("Stopped TestAPI request handling")
			return
		default:
		}

		conn, err := a.listener.Accept()
		if err != nil {
			log.Error("Accepting connection: ", err.Error())
		}
		go a.handleConnection(conn)
	}
}

func (a *testAPI) handleConnection(conn net.Conn) {
	for {
		select {
		case <-a.closed:
			return
		default:
		}

		buf := make([]byte, 1024)
		l, err := conn.Read(buf)
		if err != nil {
			log.Error("Socket reading: ", err.Error())
			return
		}
		log.Tracef("API request: '%s'", string(buf[0:l]))
		response := a.execRequest(string(buf[0:l]), conn)
		conn.Write([]byte(response))
	}
}

func (a *testAPI) execRequest(req string, conn net.Conn) string {

	// Use strings.Split to parse the string
	// The second argument is the delimiter (comma in this case)
	args := strings.Split(req, ",")
	fmt.Println(args)
	if args[0] == "getbals" {
		data, err := json.Marshal(backend.GetBals())
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		return string(data)

	} else if args[0] == "config" {
		config, err := backend.PrintConfig()
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		data, err := json.Marshal(config)
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		log.Debug("data: ", string(data))
		return string(data)

	} else if args[0] == "info" {
		info, err := backend.Info(args[1:])
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		data, err := json.Marshal(info)
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		log.Debug("data: ", string(data))
		return string(data)

	} else if args[0] == "open" {

		info, err := backend.Open(args[1:])
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		data, err := json.Marshal(info)
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		log.Debug("data: ", string(data))
		return string(data)

	} else if args[0] == "settle" {
		info, err := backend.Settle(args[1:])
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		data, err := json.Marshal(info)
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		log.Debug("data: ", string(data))
		return string(data)

	} else if args[0] == "set" {
		info, err := backend.Set(args[1:])
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		data, err := json.Marshal(info)
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		log.Debug("data: ", string(data))
		return string(data)

	} else if args[0] == "forceset" {
		info, err := backend.ForceSet(args[1:])
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		data, err := json.Marshal(info)
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		log.Debug("data: ", string(data))
		return string(data)

	} else if args[0] == "exit" {
		err := backend.Exit(args[1:])
		if err != nil {
			log.Error(err)
			return err.Error()
		}
		log.Debug("Exiting cli...")
		return "Exiting cli..."

	} else if args[0] == "help" {
		return "Commands: config, info, open, settle, exit"
	}else if err := Execute(req); err != nil {
		return err.Error()
	}
	return "OK"
}
