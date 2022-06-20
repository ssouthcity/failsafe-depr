package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/ssouthcity/dgimux"
)

func activityCommand(config *Config) dgimux.InteractionHandlerFunc {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		rw := dgimux.NewResponseWriter()
		rw.Text("[cheerful] Captain, I've readied the notifictation system! [dreary] Is this all you intend to use golden technology for?")
		rw.Ephemral()

		menu := discordgo.SelectMenu{
			CustomID:    "activities_select",
			Placeholder: "Choose the activities you would like to be notified of",
			MinValues:   NewIntPtr(0),
			MaxValues:   len(config.ActivityRoles),
		}

		for name, id := range config.ActivityRoles {
			menu.Options = append(menu.Options, discordgo.SelectMenuOption{
				Label:   strings.Title(name),
				Value:   name,
				Default: ListContainsStr(i.Member.Roles, id),
			})
		}

		rw.ComponentRow(menu)

		s.InteractionRespond(i.Interaction, rw.Response())
	}
}

func activitySelect(config *Config) dgimux.InteractionHandlerFunc {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		picks := i.MessageComponentData().Values

		rw := dgimux.NewResponseWriter()
		rw.Type(discordgo.InteractionResponseUpdateMessage)
		rw.Text("[happy] I've set you up for notifications from friends [murky] Since you're lucky enough to have them")
		rw.ClearComponentRows()
		rw.Ephemral()

		for name, id := range config.ActivityRoles {
			var err error

			wantsRole := ListContainsStr(picks, name)
			hasRole := ListContainsStr(i.Member.Roles, id)

			if wantsRole && !hasRole {
				err = s.GuildMemberRoleAdd(config.GuildID, i.Member.User.ID, id)
			} else if hasRole && !wantsRole {
				err = s.GuildMemberRoleRemove(config.GuildID, i.Member.User.ID, id)
			}

			if err != nil {
				log.Err(err).Str("user", i.Member.User.Username).Bool("add", wantsRole).Msg("role was not managed")
				rw.Type(discordgo.InteractionResponseChannelMessageWithSource)
				rw.Text("[happy] Operation failed, [pessimistic] what am I even good for...")
				break
			}
		}

		s.InteractionRespond(i.Interaction, rw.Response())
	}
}
