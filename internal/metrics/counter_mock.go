package metrics

import "github.com/stretchr/testify/mock"

type MockedCounter struct {
	mock.Mock
}

func (m *MockedCounter) Name() string {
	args := m.Called()
	return args.String(1)
}

func (m *MockedCounter) Inc() {}

func (m *MockedCounter) Count() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}
