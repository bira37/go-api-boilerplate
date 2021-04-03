package user

import (
	"reflect"
	"testing"
	"time"

	"github.com/bira37/go-rest-api/api/config"
	"github.com/bira37/go-rest-api/api/internal/errs"
	"github.com/bira37/go-rest-api/pkg/cockroach"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GenRandomUser() Model {
	faker := gofakeit.NewCrypto()

	return Model{
		Id:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Name:         faker.Name(),
		Username:     faker.Username(),
		PasswordHash: faker.Password(true, true, true, false, false, 20),
		Email:        faker.Email(),
	}
}

func checkUserModelEquality(t *testing.T, a Model, b Model) {
	deltaCreated := a.CreatedAt.Sub(b.CreatedAt)

	if deltaCreated.Seconds() >= 1 {
		t.Errorf("dates do not match: %v and %v", a.CreatedAt, b.CreatedAt)
	}

	deltaUpdated := a.UpdatedAt.Sub(b.UpdatedAt)

	if deltaUpdated.Seconds() >= 1 {
		t.Errorf("dates do not match: %v and %v (%v seconds)", a.UpdatedAt, b.UpdatedAt, deltaUpdated.Seconds())
	}

	// ignore date in comparison due to db precision
	fakeTime := time.Now().UTC()
	a.CreatedAt = fakeTime
	a.UpdatedAt = fakeTime
	b.CreatedAt = fakeTime
	b.UpdatedAt = fakeTime

	if !reflect.DeepEqual(b, a) {
		t.Errorf("expected %v, found %v", a, b)
	}
}

func TestInsert(t *testing.T) {
	store := NewStore()

	mockUser := GenRandomUser()

	connection := cockroach.NewCockroachDB(config.SQLDBConnectionString)

	_, err := store.Insert(connection.GetConnection(), mockUser)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = store.Insert(connection.GetConnection(), mockUser)

	if err == nil {
		t.Errorf("expected error inserting user with same username and id")
	}
}

func TestFindByUsername(t *testing.T) {
	store := NewStore()

	mockUser := GenRandomUser()

	connection := cockroach.NewCockroachDB(config.SQLDBConnectionString)

	_, err := store.Insert(connection.GetConnection(), mockUser)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	dbUser, err := store.FindByUsername(connection.GetConnection(), mockUser.Username)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	checkUserModelEquality(t, mockUser, dbUser)

	_, err = store.FindByUsername(connection.GetConnection(), mockUser.Username+"_suffix")

	if err == nil {
		t.Errorf("expected error, got success")
	}

	storeErr, ok := err.(*errs.StoreError)

	if !ok {
		t.Errorf("expected store error")
	}

	expectedCode := errs.StoreNotFound("").Code

	if storeErr.Code != expectedCode {
		t.Errorf("expected code %v, got %v", storeErr.Code, expectedCode)
	}
}

func TestStoreUnderTransaction(t *testing.T) {
	store := NewStore()

	mockUser := GenRandomUser()

	connection := cockroach.NewCockroachDB(config.SQLDBConnectionString)

	err := connection.Transaction(func(tx *sqlx.Tx) error {
		_, err := store.Insert(tx, mockUser)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	dbUser, err := store.FindByUsername(connection.GetConnection(), mockUser.Username)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	checkUserModelEquality(t, mockUser, dbUser)
}
