package winsleepd

import (
	"fmt"
	"os/exec"
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

var (
	user32, _      = syscall.LoadLibrary("user32.dll")
	sendMessage, _ = syscall.GetProcAddress(user32, "SendMessageW")
)

func Sleep() {
	// 0x0112 is the WM_SYSCOMMAND message
	// 0xF170 is the SC_MONITORPOWER message
	// 2 is the power-off parameter
	ret, _, _ := syscall.Syscall6(sendMessage, 4, 0xffff, 0x0112, 0xF170, 2, 0, 0)
	fmt.Printf("ret: %v\n", ret)
}

func LockScreen() {
	err := exec.Command("rundll32.exe", "user32.dll,LockWorkStation").Run()
	if err != nil {
		fmt.Println(err)
	}
}

func Hibernate() {

}

func ScreenOff() {
	err := exec.Command("powershell.exe", "(Add-Type '[DllImport(\"user32.dll\")]public static extern int SendMessage(int hWnd, int hMsg, int wParam, int lParam);' -Name a -Pas)::SendMessage(-1,0x0112,0xF170,2)").Run()
	if err != nil {
		fmt.Println(err)
	}
}
