package main

import (
	"context"
	"fmt"

	"github.com/qmuntal/stateless"
)

type PhotoBoothState string

const (
	SessionStart  PhotoBoothState = "SessionStart"
	CountDown     PhotoBoothState = "CountDown"
	SharingScreen PhotoBoothState = "SharingScreen"
	Printing      PhotoBoothState = "Printing"
	SessionEnd    PhotoBoothState = "SessionEnd"
)

func main() {

	photoBoothState := stateless.NewStateMachine(SessionStart)

	photoBoothState.Configure(SessionStart).Permit(CountDown, SharingScreen)
	photoBoothState.
		Configure(SharingScreen).
		OnEntryFrom(CountDown, func(ctx context.Context, args ...interface{}) error {
			fmt.Println("OnEntryFrom SharingScreen", args)
			return nil
		}).
		Permit(Printing, SessionEnd)

	photoBoothState.
		Configure(SessionEnd).
		OnEntry(func(ctx context.Context, args ...interface{}) error {
			fmt.Println("SessionEnd on entry", args)
			return nil
		}).
		OnExit(func(ctx context.Context, args ...interface{}) error {
			fmt.Println("SessionEnd on exit", args)
			return nil
		}).Permit(SessionStart, CountDown)

	photoBoothState.Fire(CountDown, "test")
	photoBoothState.Fire(Printing, "test")

}
