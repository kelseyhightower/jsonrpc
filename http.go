// Copyright 2019 The jsonrpc Authors. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

// Package jsonrpc implements a HTTP handler for the net/rpc
// package using the jsonrpc codec.
//
// This package is largely a workaround to avoid using the HTTP
// CONNECT method required by the HTTP handlers provided by the
// net/rpc package, which are not supported by some HTTP proxies.
package jsonrpc

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type serverCodec struct {
	r io.ReadCloser // holds the JSON formated RPC request
	w io.Writer     // holds the JSON formated RPC response
}

// Read implements the io.ReadWriteCloser Read method.
func (s *serverCodec) Read(p []byte) (n int, err error) {
	return s.r.Read(p)
}

// Write implements the io.ReadWriteCloser Write method.
func (s *serverCodec) Write(p []byte) (n int, err error) {
	return s.w.Write(p)
}

// Close implements the io.ReadWriteCloser Close method.
func (s *serverCodec) Close() error {
	return s.r.Close()
}

// Handler returns a request handler that serves rpc request by
// decoding the request and invoking the registered method using
// the given rpc server.
func Handler(server *rpc.Server) http.Handler {
	return &httpHandler{server}
}

type httpHandler struct {
	server *rpc.Server
}

// ServeHTTP implements an http.Handler that answers RPC requests.
func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	h.server.ServeCodec(jsonrpc.NewServerCodec(&serverCodec{r.Body, w}))
}
