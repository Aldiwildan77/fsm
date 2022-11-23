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

func main() {

	photoBoothState := stateless.NewStateMachine(OffHook)

	photoBoothState.
		Configure(OffHook).
		Permit(SessionStart, TakePicture)

	photoBoothState.
		Configure(TakePicture).
		OnEntryFrom(SessionStart, func(ctx context.Context, args ...interface{}) error {
			fmt.Println("SessionStart")
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

	fmt.Println(photoBoothState.ToGraph())

	ctx := context.Background()

	photoBoothState.FireCtx(ctx, SessionStart)
	photoBoothState.FireCtx(ctx, CountDown)
	photoBoothState.FireCtx(ctx, SharingScreen)
	photoBoothState.FireCtx(ctx, Printing)
	photoBoothState.FireCtx(ctx, SessionEnd)

	fmt.Println()

	lastState, _ := photoBoothState.State(ctx)
	fmt.Println("lastState:", lastState)

}
