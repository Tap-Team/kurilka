package userworker_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/Tap-Team/kurilka/internal/userworker"
	"github.com/Tap-Team/kurilka/pkg/amidsql"
	"github.com/Tap-Team/kurilka/workers"
	"gotest.tools/v3/assert"
)

var (
	db *sql.DB
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	database, term, err := amidsql.NewContainer(ctx, amidsql.DEFAULT_MIGRATION_PATH)
	if err != nil {
		log.Fatalf("failed create postgres container")
	}
	defer term(ctx)
	db = database
	os.Exit(m.Run())
}

type fakeUserWorker struct {
	t               *testing.T
	mu              *sync.Mutex
	registeredUsers map[int64]struct{}
}

func NewFakeUserWorker(t *testing.T) *fakeUserWorker {
	var mu sync.Mutex
	users := make(map[int64]struct{})
	t.Cleanup(func() {
		mu.Lock()
		defer mu.Unlock()
		if len(users) != 0 {
			t.Fatalf("no calls")
		}
	})
	return &fakeUserWorker{
		t:               t,
		mu:              &mu,
		registeredUsers: users,
	}
}

func (f *fakeUserWorker) RegisterUser(userId int64) {
	f.mu.Lock()
	f.registeredUsers[userId] = struct{}{}
	f.mu.Unlock()
}

func (f *fakeUserWorker) AddAllUsers(ctx context.Context, users []*workers.User) {}
func (f *fakeUserWorker) RemoveUser(ctx context.Context, userId int64)           {}

func (f *fakeUserWorker) AddUser(ctx context.Context, user workers.User) {
	f.mu.Lock()
	if _, ok := f.registeredUsers[user.ID]; !ok {
		f.t.Fatalf("unexpected call, undefined user id: %d", user.ID)
	}
	delete(f.registeredUsers, user.ID)
	f.mu.Unlock()
}

func TestInitWorker(t *testing.T) {
	worker := NewFakeUserWorker(t)

	usersCount := 24321

	err := insertUsers(db, 24321)
	assert.NilError(t, err, "failed insert users")
	for i := 0; i < usersCount; i++ {
		worker.RegisterUser(int64(i))
	}
	userworker.InitUserWorkerWorker(db, worker)
}
