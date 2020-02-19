package main

import (
	config "github.com/YakovBudnikov/TimetableBot/config"
	module "github.com/YakovBudnikov/TimetableBot/modules"
	"log"
)

func main() {
	config, err := config.GetConfig("../config/config.json")
	if err != nil {
		log.Fatalln(err)
	}
	db,_ := module.ConnectDataBase(config.DbArgs)
	module.StartBot(config.Token, db)
}