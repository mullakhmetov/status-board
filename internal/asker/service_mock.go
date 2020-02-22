package asker

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockedService struct {
	mock.Mock
}

func (m *MockedService) Run(ctx context.Context) {
	_ = m.Called(ctx)
	return
}

func (m *MockedService) CheckAll(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(1)
}

func (m *MockedService) Get(ctx context.Context, name string) (Response, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(Response), args.Error(1)
}

func (m *MockedService) GetMin(ctx context.Context) (Response, error) {
	args := m.Called(ctx)
	return args.Get(0).(Response), args.Error(1)
}

func (m *MockedService) GetMax(ctx context.Context) (Response, error) {
	args := m.Called(ctx)
	return args.Get(0).(Response), args.Error(1)
}

func (m *MockedService) GetRandom(ctx context.Context) (Response, error) {
	args := m.Called(ctx)
	return args.Get(0).(Response), args.Error(1)
}

func (m *MockedService) Close() {}
