package fsmManager

import (
	"context"
	"fmt"
	"memetgbot/internal/core/logger"
	"sync"

	f "github.com/looplab/fsm"
)

const (
	StateInitial        = "initial"
	StateAwaitingKey    = "awaiting_key"
	StateProcessingLink = "processing_link"
)

const (
	InitialEvent        = "initial__event"
	AwaitingKeyEvent    = "awaiting_key__event"
	ProcessingLinkEvent = "processing_link__event"
)

var events = []f.EventDesc{
	{Name: InitialEvent, Src: []string{StateInitial, StateAwaitingKey, StateProcessingLink}, Dst: StateInitial},
	{Name: AwaitingKeyEvent, Src: []string{StateInitial}, Dst: StateAwaitingKey},
	{Name: ProcessingLinkEvent, Src: []string{StateInitial}, Dst: StateProcessingLink},
}

type FSMState struct {
	users map[int64]*f.FSM
	mu    *sync.Mutex
}

func initFSM() *FSMState {
	return &FSMState{
		users: make(map[int64]*f.FSM),
		mu:    &sync.Mutex{},
	}
}

func (fsm *FSMState) GetFSMForUser(userID int64) *f.FSM {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	if userFsm, exists := fsm.users[userID]; exists {
		return userFsm
	}

	userFsm := f.NewFSM(
		StateInitial,
		events,
		f.Callbacks{
			"before_event": func(_ context.Context, e *f.Event) { fsm.beforeEvent(userID, e) },
			"after_event":  func(_ context.Context, e *f.Event) { fsm.afterEvent(userID, e) },
		},
	)
	fsm.users[userID] = userFsm
	return userFsm
}

func (fsm *FSMState) beforeEvent(userID int64, e *f.Event) {
	logger.Logger.Debug(fmt.Sprintf("User %d - Before event '%s': State '%s' -> '%s'\n", userID, e.Event, e.Src, e.Dst))
}

func (fsm *FSMState) afterEvent(userID int64, e *f.Event) {
	logger.Logger.Debug(fmt.Sprintf("User %d - After event '%s': New state '%s'\n", userID, e.Event, e.Dst))
}

func (fsm *FSMState) UserEvent(ctx context.Context, chatId int64, event string, args ...interface{}) {
	userFSM := fsm.GetFSMForUser(chatId)
	err := userFSM.Event(ctx, event, args...)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}

var FSM = initFSM()
