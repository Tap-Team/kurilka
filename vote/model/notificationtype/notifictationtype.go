package notificationtype

import "strings"

const (
	GET_SUBSCRIPTION           NotificationType = "get_subscription"
	SUBSCRIPTION_STATUS_CHANGE NotificationType = "subscription_status_change"
)

type NotificationType string

func (n NotificationType) Test() NotificationType {
	return n + "_test"
}

func (n NotificationType) Is(s string) bool {
	s = strings.TrimSuffix(s, "_test")
	return string(n) == s
}
