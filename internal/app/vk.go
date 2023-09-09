package app

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/Tap-Team/kurilka/internal/config"
)

func VK(cnf config.VKConfig) *api.VK {
	return api.NewVK(cnf.AppAccessKey, cnf.GroupAccessKey)
}
