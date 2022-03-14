package dgmux

import "github.com/bwmarrin/discordgo"

type ResponseWriter struct {
	res *discordgo.InteractionResponse
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		res: &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{},
		},
	}
}

func (r *ResponseWriter) Type(i discordgo.InteractionResponseType) {
	r.res.Type = i
}

func (r *ResponseWriter) Content(c string) {
	r.res.Data.Content = c
}

func (r *ResponseWriter) Ephemral() {
	r.res.Data.Flags = 1 << 6
}

func (r *ResponseWriter) ClearRows() {
	r.res.Data.Components = []discordgo.MessageComponent{}
}

func (r *ResponseWriter) Components(comps ...discordgo.MessageComponent) {
	r.res.Data.Components = append(r.res.Data.Components, &discordgo.ActionsRow{
		Components: comps,
	})
}

func (r *ResponseWriter) Response() *discordgo.InteractionResponse {
	return r.res
}
