package vote

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Tap-Team/kurilka/vote/handler"
	"github.com/Tap-Team/kurilka/vote/vkgetsubscripitonhandler"

	"github.com/Tap-Team/kurilka/vote/datamanager/subscriptiondatamanager"
	subscriptionitemstorage "github.com/Tap-Team/kurilka/vote/storage/local/subscriptionstorage"
	"github.com/Tap-Team/kurilka/vote/storage/postgres/subscriptionstorage"
	"github.com/Tap-Team/kurilka/vote/storage/postgres/votesubscriptionstorage"
	subscriptioncache "github.com/Tap-Team/kurilka/vote/storage/redis/subscriptionstorage"
	"github.com/Tap-Team/kurilka/vote/usecase/changesubscriptionstatususecase"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Redis     *redis.Client
	DB        *sql.DB
	ApiRouter *mux.Router
	Mux       *mux.Router

	SubscriptionConfig struct {
		Expiration time.Duration
	}
	VKConfig struct {
		VKAppSecret     string
		VKAppServiceKey string
		Version         string
	}
}

func SetUp(cnf *Config) {
	subscriptionStorage := subscriptionstorage.New(cnf.DB)
	subscriptionCache := subscriptioncache.New(cnf.Redis, cnf.SubscriptionConfig.Expiration)
	voteSubscriptionStorage := votesubscriptionstorage.New(cnf.DB)
	subscriptionItemStorage := subscriptionitemstorage.New()

	subscriptionDataManager := subscriptiondatamanager.New(subscriptionStorage, subscriptionCache)

	changeSubscriptionStatusUseCase := changesubscriptionstatususecase.New(voteSubscriptionStorage, subscriptionDataManager, subscriptionItemStorage)

	h := handler.New(cnf.VKConfig.VKAppSecret, subscriptionItemStorage, changeSubscriptionStatusUseCase, voteSubscriptionStorage)

	cnf.Mux.Handle("/vk/payments", h.HandleNotificationHandler()).Methods(http.MethodPost)
	cnf.ApiRouter.Handle("/vote-subscription/user", h.UserSubscriptionIdHandler()).Methods(http.MethodGet)
	cnf.ApiRouter.Handle("/vote-subscription/getUserSubscriptionById", vkgetsubscripitonhandler.New(http.DefaultClient, cnf.VKConfig.Version, cnf.VKConfig.VKAppServiceKey))
}
