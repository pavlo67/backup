package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Print(os.Stat("/home/pavlo"))

}
