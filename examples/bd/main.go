package main

import (
	"fmt"

	avrophonetic "github.com/mhshajib/avro-phonetic-go"
)

func main() {
	fmt.Println(avrophonetic.ToBD("tmi valo"))
}
