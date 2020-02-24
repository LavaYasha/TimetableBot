package main

import (
	"github.com/YakovBudnikov/TimetableBot/config"
	"github.com/YakovBudnikov/TimetableBot/botman"
	"log"
)

func main() {
	mconfig, err := config.GetConfig("../config/config.json")
	if err != nil {
		log.Fatalln(err)
	}
	postgres := botman.NewPostgres()
	err = postgres.Connect(mconfig.DbArgs)
	if err != nil{
		log.Fatalln(err)
	}
	bot := botman.NewBotman(mconfig.Token, postgres)
	err = bot.Run()
	if err != nil{
		log.Fatalln(err)
	}
}