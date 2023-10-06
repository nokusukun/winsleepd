package winsleepd

import (
	"fmt"
	"syscall"
	"unsafe"
)

func GetMousePos() (x, y int) {
	userDll := syscall.NewLazyDLL("user32.dll")
	getWindowRectProc := userDll.NewProc("GetCursorPos")
	type POINT struct {
		X, Y int32
	}
	var pt POINT
	_, _, eno := syscall.SyscallN(getWindowRectProc.Addr(), uintptr(unsafe.Pointer(&pt)))
	if eno != 0 {
		fmt.Println(eno)
	}
	return int(pt.X), int(pt.Y)
}

func Sleep() {
	var (
		powrProfDLL     = syscall.NewLazyDLL("powrprof.dll")
		setSuspendState = powrProfDLL.NewProc("SetSuspendState")
	)
	ret, _, _ := setSuspendState.Call(0, 0, 0)
	if ret != 0 {
		fmt.Println("Computer is entering sleep mode")
	} else {
		fmt.Println("SetSuspendState failed")
	}
}

func LockScreen() {
	var (
		user32          = syscall.NewLazyDLL("user32.dll")
		lockWorkStation = user32.NewProc("LockWorkStation")
	)

	ret, _, _ := lockWorkStation.Call()
	if ret == 0 {
		fmt.Println("LockWorkStation failed")
	} else {
		fmt.Println("PC locked")
	}
}

func Hibernate() {

}

func ScreenOff() {
	var (
		user32                 = syscall.NewLazyDLL("user32.dll")
		sendMessageTimeout     = user32.NewProc("SendMessageTimeoutW")
		hwndBroadcast          = uintptr(0xffff)
		wmSysCommand           = uintptr(0x0112)
		scMonitorPower         = uintptr(0xF170)
		monitorOff             = uintptr(2)
		smtoNotimeoutifnothung = uintptr(0x0002)
	)

	ret, _, _ := sendMessageTimeout.Call(
		hwndBroadcast,
		wmSysCommand,
		scMonitorPower,
		monitorOff,
		smtoNotimeoutifnothung,
		1000,
		0,
	)

	if ret == 0 {
		fmt.Println("SendMessageTimeout failed")
	} else {
		fmt.Println("Monitor power state changed")
	}
}
