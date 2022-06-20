package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ssouthcity/dgimux"
)

type Config struct {
	Token         string            `json:"token"`
	GuildID       string            `json:"guild"`
	ClassRoles    map[string]string `json:"classes"`
	ActivityRoles map[string]string `json:"activities"`
}

var (
	confPath = flag.String("config", "config.json", "path to config file")
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	flag.Parse()

	conf := &Config{}

	confFile, err := os.Open(*confPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", "").Msg("config file does not exist")
	}
	defer confFile.Close()

	if err := json.NewDecoder(confFile).Decode(&conf); err != nil {
		log.Fatal().Err(err).Msg("config file contains invalid json")
	}

	sess, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("session was incorrectly configured")
	}

	r := dgimux.NewRouter()

	r.ApplicationCommand("class", classCommand(conf))
	r.MessageComponent("class_select", classSelect(conf))

	r.ApplicationCommand("activities", activityCommand(conf))
	r.MessageComponent("activities_select", activitySelect(conf))

	sess.AddHandler(r.HandleInteraction)

	if err := sess.Open(); err != nil {
		log.Fatal().Err(err).Msg("websocket connection could not be established")
	}
	defer sess.Close()

	select {}
}
