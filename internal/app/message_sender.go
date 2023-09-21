package app

import (
	"net/http"

	"github.com/Tap-Team/kurilka/internal/messagesender"
	"github.com/Tap-Team/kurilka/internal/messagesender/vk"
)

func MessageSender(
	apiVersion string,
	token string,
) messagesender.MessageSender {
	return vk.NewMessageSender(http.DefaultClient, apiVersion, token)
}
