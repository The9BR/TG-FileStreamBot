package commands

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/utils"
	"strings"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
)

func (m *command) LoadJoin(dispatcher dispatcher.Dispatcher) {
	log := m.log.Named("join")
	defer log.Sugar().Info("Loaded")
	dispatcher.AddHandler(handlers.NewCommand("join", join))
}

func join(ctx *ext.Context, u *ext.Update) error {
	chatId := u.EffectiveChat().GetID()
	peerChatId := ctx.PeerStorage.GetPeerById(chatId)
	if peerChatId.Type != int(storage.TypeUser) {
		return dispatcher.EndGroups
	}
	if len(config.ValueOf.AllowedUsers) != 0 && !utils.Contains(config.ValueOf.AllowedUsers, chatId) {
		ctx.Reply(u, ext.ReplyTextString("You are not allowed to use this bot."), nil)
		return dispatcher.EndGroups
	}

	// Get the channel link from the message
	args := strings.Fields(u.EffectiveMessage.GetText())
	if len(args) < 2 {
		ctx.Reply(u, ext.ReplyTextString("Usage: /join <channel_link>\nExample: /join https://t.me/+BMLpmVmqbx5lNzg0"), nil)
		return dispatcher.EndGroups
	}

	channelLink := args[1]

	// Try to join the channel
	_, err := ctx.Client.JoinChatWithInviteLink(ctx.Context, &channelLink)
	if err != nil {
		ctx.Reply(u, ext.ReplyTextString("❌ Failed to join channel: "+err.Error()), nil)
		return dispatcher.EndGroups
	}

	ctx.Reply(u, ext.ReplyTextString("✅ Successfully joined the channel!"), nil)
	return dispatcher.EndGroups
}
