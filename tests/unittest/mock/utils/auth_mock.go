package mockutils

import "scheduler/internal/app/models"

type MockTokenGenerator struct {
	GenerateTokenFunc func(uint, string, string) (*models.JWT, error)
}

func (m *MockTokenGenerator) GenerateToken(id uint, email, role string) (*models.JWT, error) {
	return m.GenerateTokenFunc(id, email, role)
}