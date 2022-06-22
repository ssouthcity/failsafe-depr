package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	token := flag.String("token", "", "token used for authentication to discord")
	flag.Parse()

	cmds := []*discordgo.ApplicationCommand{}

	cmdsFile, err := os.Open("commands.json")
	if err != nil {
		log.Fatal().Err(err).Msg("command file does not exist")
	}
	defer cmdsFile.Close()

	if err := json.NewDecoder(cmdsFile).Decode(&cmds); err != nil {
		log.Fatal().Err(err).Msg("command file contains invalid json")
	}

	s, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatal().Err(err).Msg("session was incorrectly configured")
	}

	if err := s.Open(); err != nil {
		log.Fatal().Err(err).Msg("websocket connection could not be established")
	}
	defer s.Close()

	if _, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", cmds); err != nil {
		log.Fatal().Err(err).Msg("command creation failed")
	}

	log.Info().Msg("successfully synchronized application commands")
}
