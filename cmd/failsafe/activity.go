package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/ssouthcity/failsafe/pkg/dgmux"
)

func activityCommand(config *Config) dgmux.Handler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		rw := dgmux.NewResponseWriter()
		rw.Content("[cheerful] Captain, I've readied the notifictation system! [dreary] Is this all you intend to use golden technology for?")
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

		rw.Components(menu)

		s.InteractionRespond(i.Interaction, rw.Response())
	}
}

func activitySelect(config *Config) dgmux.Handler {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		picks := i.MessageComponentData().Values

		rw := dgmux.NewResponseWriter()
		rw.Type(discordgo.InteractionResponseUpdateMessage)
		rw.Content("[happy] I've set you up for notifications from friends [murky] Since you're lucky enough to have them")
		rw.ClearRows()
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
				rw.Content("[happy] Operation failed, [pessimistic] what am I even good for...")
				break
			}
		}

		s.InteractionRespond(i.Interaction, rw.Response())
	}
}
