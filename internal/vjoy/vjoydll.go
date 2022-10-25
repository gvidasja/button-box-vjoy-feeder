package vjoy

import (
	"fmt"
	"sync"
	"syscall"

	"github.com/gvidasja/button-box-vjoy-feeder/internal/device"
)

var lock sync.Once

var vjoyDll = syscall.NewLazyDLL("vJoyInterface.dll")

var (
	procSetButton     = vjoyDll.NewProc("SetBtn")
	procRelinquishVJD = vjoyDll.NewProc("RelinquishVJD")
	procAcquireVJD    = vjoyDll.NewProc("AcquireVJD")
	procVJoyEnabled   = vjoyDll.NewProc("vJoyEnabled")
	procGetVJDStatus  = vjoyDll.NewProc("GetVJDStatus")
)

const (
	VJD_STAT_OWN = iota
	VJD_STAT_FREE
	VJD_STAT_BUSY
	VJD_STAT_MISS
	VJD_STAT_UNKN
)

func loadVJoyDLL() error {
	var err error
	lock.Do(func() { err = vjoyDll.Load() })
	return err
}

func vJoyEnabled() bool {
	enabled, _, _ := procVJoyEnabled.Call()
	return enabled != 0
}

func acquireVJD(deviceID uint) error {
	acquired, _, _ := procAcquireVJD.Call(uintptr(deviceID))
	if acquired == 0 {
		return fmt.Errorf("could not acquire device %d", deviceID)
	}

	return nil
}

func relinquishVJD(deviceID uint) error {
	relinquished, _, _ := procRelinquishVJD.Call(uintptr(deviceID))
	if relinquished == 0 {
		return fmt.Errorf("could not dispose device %d", deviceID)
	}

	return nil
}

func getVJDStatus(deviceID uint) int {
	status, _, _ := procGetVJDStatus.Call(uintptr(deviceID))

	return int(status)
}

func setButton(deviceID uint, buttonID device.ButtonID, state bool) error {
	var stateInt uintptr

	if state {
		stateInt = 1
	} else {
		stateInt = 0
	}

	stateWasSet, _, _ := procSetButton.Call(uintptr(stateInt), uintptr(deviceID), uintptr(buttonID))

	if stateWasSet == 0 {
		return fmt.Errorf("could not set button %d state to %v on device %d", buttonID, state, deviceID)
	}

	return nil
}
