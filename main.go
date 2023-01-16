package main

import (
	"context"
	"fmt"

	"github.com/qmuntal/stateless"
)

type PhotoBoothTrigger string

const (
	SessionStart  PhotoBoothTrigger = "SessionStart"
	CountDown     PhotoBoothTrigger = "CountDown"
	SharingScreen PhotoBoothTrigger = "SharingScreen"
	Printing      PhotoBoothTrigger = "Printing"
	SessionEnd    PhotoBoothTrigger = "SessionEnd"
)

type PhotoBoothState string

const (
	OffHook      PhotoBoothState = "OffHook"
	TakePicture  PhotoBoothState = "TakePicture"
	ChromeWindow PhotoBoothState = "ChromeWindow"
)

type PhotoBoothMachine struct {
	photoBoothState *stateless.StateMachine
}

func NewPhotoBoothMachine(initialState any) PhotoBoothMachine {
	return PhotoBoothMachine{
		photoBoothState: stateless.NewStateMachine(initialState),
	}
}

func (p *PhotoBoothMachine) Register() error {
	photoBoothState := p.photoBoothState

	photoBoothState.
		Configure(OffHook).
		Permit(SessionStart, TakePicture)

	photoBoothState.
		Configure(TakePicture).
		OnEntryFrom(SessionStart, func(ctx context.Context, args ...interface{}) error {
			fmt.Println("OnEntryFrom, SessionStart")
			return nil
		}).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			fmt.Println("OnEntry, TakePicture")
			return nil
		}).
		InternalTransition(SharingScreen, func(ctx context.Context, args ...interface{}) error {
			fmt.Println("InternalTransition, SharingScreen")
			return nil
		}).
		InternalTransition(Printing, func(ctx context.Context, args ...interface{}) error {
			fmt.Println("InternalTransition, Printing")
			return nil
		}).
		PermitReentry(CountDown).
		Permit(SessionEnd, ChromeWindow)

	photoBoothState.
		Configure(ChromeWindow).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			fmt.Println("OnEntry, ChromeWindow")
			return nil
		}).
		Permit(SessionEnd, OffHook)

	return nil
}

func (p *PhotoBoothMachine) DrawGraph() string {
	return p.photoBoothState.ToGraph()
}

func main() {
	fireStates := make([][]PhotoBoothTrigger, 0)
	fireStates = append(fireStates, []PhotoBoothTrigger{SessionStart, CountDown, SharingScreen, Printing, SessionEnd})
	fireStates = append(fireStates, []PhotoBoothTrigger{SessionStart, CountDown, SharingScreen, Printing})
	fireStates = append(fireStates, []PhotoBoothTrigger{SessionStart, CountDown, Printing, SessionEnd})
	fireStates = append(fireStates, []PhotoBoothTrigger{CountDown, Printing, SessionEnd})
	fireStates = append(fireStates, []PhotoBoothTrigger{SessionStart, CountDown, SharingScreen, SessionEnd})

	for _, fireState := range fireStates {
		photoBoothMachine := NewPhotoBoothMachine(OffHook)
		photoBoothMachine.Register()
		photoBoothMachine.DrawGraph()

		ctx := context.Background()

		for _, trigger := range fireState {
			fmt.Printf("\ntrigger: %s\n\n", trigger)
			photoBoothMachine.photoBoothState.FireCtx(ctx, trigger)
			lastState, _ := photoBoothMachine.photoBoothState.State(ctx)
			fmt.Println("lastState:", lastState)
		}
		fmt.Println("--------------------")
	}
}

// FLOW
// 1 session start: do nothing
// 2 count down: do nothing
// 3 sharing screen:
//   - lock screen
//   - kill chrome
//   - start new chrome with KIOSK
// 4 printing:
//   - get current transaction
// 	 - update transaction status
//   - set webhook write access
// 5 session end:
//   - lock screen
//   - if on session start and on count down:
//    - unlock screen
//    - start photobooth session
//    - switch to photobooth application
//   - if first time init app:
//    - kill chrome
//    - start new chrome with KIOSK
//    - sleep 3s
//    - switch to chrome
//   - if sharing pass and printing pass and lock screen pass:
//    - switch to chrome

// 1. session_end
// 2. countdown_start
// 3. countdown
// 4. capture_start
// 5. file_download
// 6. loop 2-5 (max N times)
// 7. processing_start
// 8. sharing_screen
// 9. printing
// 10. file_upload
// 11. session_end
