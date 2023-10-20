package subscriptionerror

import "net/http"

var SubscriptionNotFound subscriptionNotFound

type subscriptionNotFound struct{}

func (s subscriptionNotFound) Error() string {
	return "Подписки не существует"
}

func (s subscriptionNotFound) HttpCode() int {
	return http.StatusNotFound
}

func (s subscriptionNotFound) Code() int {
	return 20
}

func (s subscriptionNotFound) Critical() bool {
	return true
}

var SubscriptionIdExists subscriptionIdExists

type subscriptionIdExists struct{}

func (s subscriptionIdExists) Error() string {
	return "Подписка пользователя уже существует"
}

func (s subscriptionIdExists) HttpCode() int {
	return http.StatusBadRequest
}

func (s subscriptionIdExists) Code() int {
	return 2
}
func (s subscriptionIdExists) Critical() bool {
	return false
}
