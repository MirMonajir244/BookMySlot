package tests

import (
	"errors"
	"testing"

	"github.com/MirMonajir244/BookMySlot/internal/models"
	"github.com/MirMonajir244/BookMySlot/internal/service"
)

// Mock repos for AuthService unit test (using simple closure/map)
type mockUserRepo struct {
	users map[string]*models.User
}

func (m *mockUserRepo) Create(u *models.User) error {
	if _, ok := m.users[u.Email]; ok {
		return errors.New("duplicate email")
	}
	m.users[u.Email] = u
	return nil
}
func (m *mockUserRepo) FindByEmail(e string) (*models.User, error) {
	if u, ok := m.users[e]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}
func (m *mockUserRepo) FindByID(id uint) (*models.User, error) { return nil, nil }

func TestAuthService_RegisterUnit(t *testing.T) {
	userRepo := &mockUserRepo{users: make(map[string]*models.User)}
	// We only need userRepo for this specific test
	s := service.NewAuthService(userRepo, nil, nil)

	// 1. Valid Registration
	_, err := s.RegisterUser("John", "john@test.com", "pass123")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// 2. Duplicate Email Registration
	_, err = s.RegisterUser("John2", "john@test.com", "pass456")
	if err == nil || err.Error() != "email already registered" {
		t.Errorf("expected 'email already registered' error, got %v", err)
	}
}
