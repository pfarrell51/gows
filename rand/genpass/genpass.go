package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"github.com/sethvargo/go-password/password"
)

func main() {
	gen, err := password.NewGenerator(nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := gen.Generate(64, 10, 10, false, false)
	if err != nil {
		log.Fatal(err)
	}
	res = res[20:]
	encoded := base64.StdEncoding.EncodeToString([]byte(res))
	fmt.Println(encoded)
}
