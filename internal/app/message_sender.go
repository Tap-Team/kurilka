package app

import (
	"net/http"

	"github.com/Tap-Team/kurilka/internal/config"
	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/internal/messagesender/vk"
)

func MessageSender(
	cnf config.VKConfig,
) messagesender.MessageSender {
	return vk.NewMessageSender(http.DefaultClient, cnf.ApiVersion, cnf.GroupAccessKey)
}
