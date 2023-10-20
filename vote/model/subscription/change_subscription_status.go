package subscription

type SubscriptionStatus string

const (
	CHARGEABLE SubscriptionStatus = "chargeable"
	ACTIVE     SubscriptionStatus = "active"
	CANCELLED  SubscriptionStatus = "cancelled"
)

type CancelReason string

const (
	USER_DECISION CancelReason = "user_decision"
	APP_DECISION  CancelReason = "app_decision"
	PAYMENT_FAIL  CancelReason = "payment_fail"
	UNKNOWN       CancelReason = "unknown"
)

type ChangeSubscriptionStatusResponse struct {
	SubscriptionId int64 `json:"subscription_id"`
	AppOrderId     int64 `json:"app_order_id,omitempty"`
}

type ChangeSubscriptionStatus struct {
	SubscriptionId int64
	UserId         int64
	ItemId         string
	Status         SubscriptionStatus
	CancelReason   CancelReason
}

func NewChangeSubscriptionStatus(
	subscriptionId int64,
	userId int64,
	itemId string,
	status SubscriptionStatus,
	cancelReason CancelReason,
) ChangeSubscriptionStatus {
	return ChangeSubscriptionStatus{
		SubscriptionId: subscriptionId,
		UserId:         userId,
		ItemId:         itemId,
		Status:         status,
		CancelReason:   cancelReason,
	}
}
