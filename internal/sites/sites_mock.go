package sites

import "github.com/stretchr/testify/mock"

type MockedService struct {
	mock.Mock
}

func (m *MockedService) Warmup() error {
	args := m.Called()
	return args.Error(1)
}

func (m *MockedService) GetAll() []*Site {
	args := m.Called()
	return args.Get(0).([]*Site)
}
func (m *MockedService) GetAvailable() []*Site {
	args := m.Called()
	return args.Get(0).([]*Site)
}
func (m *MockedService) GetSortedByLatency() []*Site {
	args := m.Called()
	return args.Get(0).([]*Site)
}

func (m *MockedService) Close() {}
