package ping

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type PingCommand struct{}

var _ ken.SlashCommand = (*PingCommand)(nil)

func (c *PingCommand) Name() string {
	return "ping"
}

func (c *PingCommand) Description() string {
	return "Basic Ping Command"
}

func (c *PingCommand) Version() string {
	return "1.0.0"
}

func (c *PingCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *PingCommand) Run(ctx ken.Context) (err error) {
	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
		},
	})
	return
}
