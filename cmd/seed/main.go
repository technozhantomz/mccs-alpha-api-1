package main

import (
	"github.com/technoshantoms/mccs-alpha-api/global"
	"github.com/technoshantoms/mccs-alpha-api/internal/seed"
)

func main() {
	global.Init()
	seed.LoadData()
	seed.Run()
}
