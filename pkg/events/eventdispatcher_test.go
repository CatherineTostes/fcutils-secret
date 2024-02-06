package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name     string
	Datetime time.Time
	Payload  interface{}
}

func (t TestEvent) GetName() string {
	return t.Name
}

func (t TestEvent) GetDateTime() time.Time {
	return t.Datetime
}

func (t TestEvent) GetPayload() interface{} {
	return t.Payload
}

type TestEventHandler struct {
	ID int
}

func (th *TestEventHandler) Handler(event EventInterface, wg *sync.WaitGroup) {}

type EventTestDispatcherSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

type MockHandler struct {
	mock.Mock
}

func (suite *EventTestDispatcherSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.handler = TestEventHandler{
		ID: 1,
	}
	suite.handler2 = TestEventHandler{
		ID: 2,
	}
	suite.handler3 = TestEventHandler{
		ID: 3,
	}
	suite.event = TestEvent{
		Name:     "test",
		Datetime: time.Now(),
		Payload:  "test",
	}
}

func TestSuite(testing *testing.T) {
	suite.Run(testing, new(EventTestDispatcherSuite))
}

func (m *MockHandler) Handler(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suite *EventTestDispatcherSuite) TestEventDispatcherRegister() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.Handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.Handlers[suite.event.GetName()]))

	assert.Equal(suite.T(), &suite.handler, suite.eventDispatcher.Handlers[suite.event.GetName()][0])
	assert.Equal(suite.T(), &suite.handler2, suite.eventDispatcher.Handlers[suite.event.GetName()][1])
}

func (suite *EventTestDispatcherSuite) TestEventDispatcherRegisterError() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.Handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Equal(errorEventAlreadyRegistered, err)
	suite.Equal(1, len(suite.eventDispatcher.Handlers[suite.event.GetName()]))
}

func (suite *EventTestDispatcherSuite) TestEventDispatcherClear() {
	// event 1
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.Handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.Handlers[suite.event.GetName()]))

	// event 2
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.Handlers[suite.event2.GetName()]))

	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.Handlers[suite.event.GetName()]))
}

func (suite *EventTestDispatcherSuite) TestEventDispatcherHas() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.Handlers[suite.event.GetName()]))

	assert.True(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler))
	assert.False(suite.T(), suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2))
}

func (suite *EventTestDispatcherSuite) TestEventDispatcherDispatch() {
	mh := &MockHandler{}
	mh.On("Handler", suite.event)

	mh2 := &MockHandler{}
	mh2.On("Handler", suite.event)

	suite.eventDispatcher.Register(suite.event.GetName(), mh)
	suite.eventDispatcher.Register(suite.event.GetName(), mh2)

	suite.eventDispatcher.Dispatch(suite.event)

	mh.AssertExpectations(suite.T())
	mh2.AssertExpectations(suite.T())

	mh.AssertNumberOfCalls(suite.T(), "Handler", 1)
	mh2.AssertNumberOfCalls(suite.T(), "Handler", 1)
}

func (suite *EventTestDispatcherSuite) TestEventDispatcherRemove() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.Handlers[suite.event.GetName()]))

	suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	suite.Equal(0, len(suite.eventDispatcher.Handlers[suite.event.GetName()]))
}
