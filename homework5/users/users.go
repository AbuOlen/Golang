package users

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	store "hw5/documentstore"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	coll store.Collection
}

func NewService() *Service {
	ok, coll := store.NewStore().CreateCollection("users", &store.CollectionConfig{PrimaryKey: "ID"})
	if !ok {
		panic("failed to create collection")
	}
	return &Service{
		coll: *coll,
	}
}

func (s *Service) CreateUser(name string) (*User, error) {
	_, err := s.GetUser(name)
	if err == nil {
		return nil, ErrCollectionAlreadyExists
	}
	user := User{ID: uuid.New().String(), Name: name}
	doc, err := MarshalDocument(user)
	if err != nil {
		return nil, err
	}
	s.coll.Put(*doc)
	return &user, nil
}

func (s *Service) ListUsers() ([]User, error) {
	var users []User
	var errs []error
	for _, doc := range s.coll.List() {
		var user = User{}
		err := UnmarshalDocument(&doc, &user)
		if err != nil {
			errs = append(errs, fmt.Errorf("unmarshaling error: %w", ErrWrongDataType))
		} else {
			users = append(users, user)
		}
	}
	return users, errors.Join(errs...)
}

func (s *Service) GetUser(userID string) (*User, error) {
	doc, ok := s.coll.Get(userID)
	if !ok {
		return nil, ErrUserNotFound
	}
	var user = User{}
	err := UnmarshalDocument(doc, &user)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling error: %w", err)
	}
	return &user, nil
}

func (s *Service) DeleteUser(userID string) error {
	ok := s.coll.Delete(userID)
	if !ok {
		return ErrUserNotFound
	}
	return nil
}
