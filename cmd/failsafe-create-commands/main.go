package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Token   string `json:"token"`
	GuildID string `json:"guild"`
}

var (
	confPath = flag.String("config", "config.json", "path to config file")
	cmdsPath = flag.String("commands", "commands.json", "path to command spec file")
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func init() {
	flag.Parse()
}

func main() {
	cmds := []*discordgo.ApplicationCommand{}

	cmdsFile, err := os.Open(*cmdsPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", *cmdsPath).Msg("command file does not exist")
	}
	defer cmdsFile.Close()

	if err := json.NewDecoder(cmdsFile).Decode(&cmds); err != nil {
		log.Fatal().Err(err).Msg("command file contains invalid json")
	}

	conf := &Config{}

	confFile, err := os.Open(*confPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", *confPath).Msg("config file does not exist")
	}
	defer confFile.Close()

	if err := json.NewDecoder(confFile).Decode(&conf); err != nil {
		log.Fatal().Err(err).Msg("config file contains invalid json")
	}

	s, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		log.Fatal().Err(err).Msg("session was incorrectly configured")
	}

	if err := s.Open(); err != nil {
		log.Fatal().Err(err).Msg("websocket connection could not be established")
	}
	defer s.Close()

	if _, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, conf.GuildID, cmds); err != nil {
		log.Fatal().Err(err).Msg("command creation failed")
	}

	log.Info().Msg("successfully synchronized application commands")
}
