package subscriptionstorage_test

import (
	"context"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Tap-Team/kurilka/internal/model/usermodel"
	"github.com/Tap-Team/kurilka/pkg/rediscontainer"
	"github.com/Tap-Team/kurilka/user/database/redis/subscriptionstorage"
	"github.com/redis/go-redis/v9"
	"gotest.tools/v3/assert"
)

var (
	rc      *redis.Client
	storage *subscriptionstorage.Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	rd, term, err := rediscontainer.New(ctx)
	if err != nil {
		log.Fatalf("failed create redis container, %s", err)
	}
	defer term(ctx)
	rc = rd
	storage = subscriptionstorage.New(rc, time.Hour*10)
	os.Exit(m.Run())
}

func Test_Subscription_Cycle(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		subscription       usermodel.Subscription
		updateSubscription usermodel.Subscription
	}{
		{
			subscription:       usermodel.NewSubscription(usermodel.BASIC, time.Now().Add(time.Hour)),
			updateSubscription: usermodel.NewSubscription(usermodel.NONE, time.Time{}),
		},
		{
			subscription:       usermodel.NewSubscription(usermodel.NONE, time.Time{}),
			updateSubscription: usermodel.NewSubscription(usermodel.BASIC, time.Now().Add(time.Hour)),
		},
		{
			subscription:       usermodel.NewSubscription(usermodel.TRIAL, time.Now().Add(time.Hour*24*5)),
			updateSubscription: usermodel.NewSubscription(usermodel.NONE, time.Time{}),
		},
	}

	for _, cs := range cases {
		userId := rand.Int63()

		err := storage.UpdateUserSubscription(ctx, userId, cs.subscription)
		assert.NilError(t, err, "non nil error")

		subscription, err := storage.UserSubscription(ctx, userId)
		subscriptionEqual(t, cs.subscription, subscription)

		err = storage.UpdateUserSubscription(ctx, userId, cs.updateSubscription)
		assert.NilError(t, err, "non nil error")

		subscription, err = storage.UserSubscription(ctx, userId)
		subscriptionEqual(t, cs.updateSubscription, subscription)
	}

}

func subscriptionEqual(t *testing.T, sub1, sub2 usermodel.Subscription) {
	assert.Equal(t, sub1.Expired.Unix(), sub2.Expired.Unix(), "subscription time not equal")
	assert.Equal(t, sub1.Type, sub2.Type, "subscription type not equal")
}
