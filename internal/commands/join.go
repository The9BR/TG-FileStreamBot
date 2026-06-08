package commands

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/utils"
	"strings"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/tg"
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
	args := strings.Fields(u.EffectiveMessage.Message)
	if len(args) < 2 {
		ctx.Reply(u, ext.ReplyTextString("Usage: /join <channel_link>\nExample: /join https://t.me/+BMLpmVmqbx5lNzg0"), nil)
		return dispatcher.EndGroups
	}

	channelLink := args[1]

	// Extract the invite hash from the link (the part after t.me/+ or t.me/joinchat/)
	inviteHash := channelLink
	if idx := strings.LastIndex(channelLink, "/+"); idx != -1 {
		inviteHash = channelLink[idx+2:]
	} else if idx := strings.LastIndex(channelLink, "/joinchat/"); idx != -1 {
		inviteHash = channelLink[idx+10:]
	}

	// Try to join the channel using the raw Telegram API
	_, err := ctx.Raw.MessagesImportChatInvite(ctx, &tg.MessagesImportChatInviteRequest{
		Hash: inviteHash,
	})
	if err != nil {
		ctx.Reply(u, ext.ReplyTextString("❌ Failed to join channel: "+err.Error()), nil)
		return dispatcher.EndGroups
	}

	ctx.Reply(u, ext.ReplyTextString("✅ Successfully joined the channel!"), nil)
	return dispatcher.EndGroups
}
