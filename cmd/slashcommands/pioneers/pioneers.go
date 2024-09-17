package pioneers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/thenanox/ficsit-discord-bot/internal/satisfactory"
	"github.com/zekrotja/ken"
)

type PioneersCommand struct{}

var _ ken.SlashCommand = (*PioneersCommand)(nil)

func (c *PioneersCommand) Name() string {
	return "pioneers"
}

func (c *PioneersCommand) Description() string {
	return "Query the pioneers in the server"
}

func (c *PioneersCommand) Version() string {
	return "1.0.0"
}

func (c *PioneersCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{}
}

func (c *PioneersCommand) Run(ctx ken.Context) (err error) {
	state, err := satisfactory.QueryServerState()
	if err != nil {
		return err
	}
	err = ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%d/%d pioneers in MASSAGE-2(A-B)b", state.Data.ServerGameState.NumConnectedPlayers, state.Data.ServerGameState.PlayerLimit),
		},
	})
	return
}
