package main

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"

// 	"github.com/lsowen/hoover-dam/pkg/api"
// )

// func main() {
// 	ctx := context.Background()
// 	r, err := api.Serve(ctx)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	http.ListenAndServe(":8080", r)
// }

import "github.com/lsowen/hoover-dam/cmd/hooverdam/cmd"

func main() {
	cmd.Execute()
}
