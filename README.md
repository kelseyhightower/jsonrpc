# jsonrpc

[![GoDoc](https://godoc.org/github.com/kelseyhightower/jsonrpc?status.svg)](https://godoc.org/github.com/kelseyhightower/jsonrpc)

## Example Usage

```Go
package main

import (
    "log"
    "net/http"
    "net/rpc"

    "github.com/kelseyhightower/jsonrpc"
)

type Args struct {
    A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}

func main() {
    arith := new(Arith)
    rpc.Register(arith)
    http.Handle("/", jsonrpc.Handler(rpc.DefaultServer))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
``` 

Build and start the server, then invoke the `Arith.Multiply` method with curl:

```
curl http://127.0.0.1:8080 \
  -d '{"method":"Arith.Multiply","params":[{"A": 10, "B":2}], "id": 0}'
```

```
{"id":0,"result":20,"error":null}
```
