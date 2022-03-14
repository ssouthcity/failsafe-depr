package dgmux

import "github.com/bwmarrin/discordgo"

type key struct {
	discordgo.InteractionType
	string
}

type Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type Mux struct {
	handlers map[key]Handler
}

func New() *Mux {
	return &Mux{make(map[key]Handler)}
}

func (m *Mux) handle(i discordgo.InteractionType, k string, h Handler) {
	m.handlers[key{i, k}] = h
}

func (m *Mux) Command(name string, h Handler) {
	m.handle(discordgo.InteractionApplicationCommand, name, h)
}

func (m *Mux) Component(customID string, h Handler) {
	m.handle(discordgo.InteractionMessageComponent, customID, h)
}

func (m *Mux) Autocomplete(name string, h Handler) {
	m.handle(discordgo.InteractionApplicationCommandAutocomplete, name, h)
}

func (m *Mux) Modal(name string, h Handler) {
	m.handle(discordgo.InteractionModalSubmit, name, h)
}

func (m *Mux) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := m.handlers[getKey(i)]; ok {
		h(s, i)
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "I am unable to handle this command at the moment",
				Flags:   1 << 6,
			},
		})
	}
}

func getKey(i *discordgo.InteractionCreate) key {
	var k string
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		k = i.ApplicationCommandData().Name
	case discordgo.InteractionMessageComponent:
		k = i.MessageComponentData().CustomID
	case discordgo.InteractionApplicationCommandAutocomplete:
		k = i.ApplicationCommandData().Name
	case discordgo.InteractionModalSubmit:
		k = i.ModalSubmitData().CustomID
	}
	return key{i.Type, k}
}
