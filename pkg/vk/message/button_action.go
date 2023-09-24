package message

type Action map[string]map[string]any

func NewOpenAppAction(appId int, ownerId int, payload string, label string, hash string) Action {
	return map[string]map[string]any{
		"action": {
			"type":     "open_app",
			"owner_id": ownerId,
			"app_id":   appId,
			"payload":  payload,
			"label":    label,
			"hash":     hash,
		},
	}
}
