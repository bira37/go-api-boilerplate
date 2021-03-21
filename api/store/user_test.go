package store

import (
	"reflect"
	"testing"
	"time"

	"github.com/bira37/go-rest-api/api/domain/user"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GenUser() user.Model {
	faker := gofakeit.NewCrypto()

	return user.Model{
		Id:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Name:         faker.Name(),
		Username:     faker.Username(),
		PasswordHash: faker.Password(true, true, true, false, false, 20),
		Email:        faker.Email(),
	}
}

func checkUserModelEquality(t *testing.T, a user.Model, b user.Model) {
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
	store := User{}

	mockUser := GenUser()

	connection := Connection.GetConnection()

	_, err := store.Insert(mockUser, connection)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = store.Insert(mockUser, connection)

	if err == nil {
		t.Errorf("expected error inserting user with same username and id")
	}
}

func TestFindByUsername(t *testing.T) {
	store := User{}

	mockUser := GenUser()

	connection := Connection.GetConnection()

	_, err := store.Insert(mockUser, connection)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	dbUser, err := store.FindByUsername(mockUser.Username, connection)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	checkUserModelEquality(t, mockUser, dbUser)
}

func TestStoreUnderTransaction(t *testing.T) {
	store := User{}

	mockUser := GenUser()

	connection := Connection.GetConnection()

	err := Connection.Transaction(func(tx *sqlx.Tx) error {
		_, err := store.Insert(mockUser, tx)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		return nil
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	dbUser, err := store.FindByUsername(mockUser.Username, connection)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	checkUserModelEquality(t, mockUser, dbUser)
}
