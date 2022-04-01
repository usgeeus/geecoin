package main

import (
	"github.com/usgeeus/geecoin.git/cli"
	"github.com/usgeeus/geecoin.git/db"
)

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()
}
