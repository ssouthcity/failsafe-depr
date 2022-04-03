package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ssouthcity/failsafe/pkg/dgmux"
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

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func init() {
	flag.Parse()
}

func main() {
	conf := &Config{}

	confFile, err := os.Open(*confPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", *confPath).Msg("config file does not exist")
	}
	defer confFile.Close()

	if err := json.NewDecoder(confFile).Decode(&conf); err != nil {
		log.Fatal().Err(err).Msg("config file contains invalid json")
	}

	sess, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("session was incorrectly configured")
	}

	sess.AddHandler(ready())
	sess.AddHandler(interactionCreate(conf))

	if err := sess.Open(); err != nil {
		log.Fatal().Err(err).Msg("websocket connection could not be established")
	}
	defer sess.Close()

	select {}
}

func ready() interface{} {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info().Str("name", r.User.Username).Msg("connected and listening for events")
	}
}

func interactionCreate(conf *Config) interface{} {
	r := dgmux.New()

	r.Command("class", classCommand(conf))
	r.Component("class_select", classSelect(conf))

	r.Command("activities", activityCommand(conf))
	r.Component("activities_select", activitySelect(conf))

	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Info().Str("interaction", i.Type.String()).Str("user", i.Member.User.Username).Msg("received interaction")

		r.HandleInteraction(s, i)
	}
}
