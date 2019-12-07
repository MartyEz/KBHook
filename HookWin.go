package KBHook

// Reading docs is needed
// https://docs.microsoft.com/fr-fr/windows/win32/winmsg/about-hooks?redirectedfrom=MSDN#wh_msgfilter_wh_sysmsgfilterhooks
// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms644985(v=vs.85)?redirectedfrom=MSDN
// https://en.wikipedia.org/wiki/Hooking
// Some Codes sources :
// https://causeyourestuck.io/2016/04/20/keyboard-hook-win32api-2/
// https://github.com/moutend/go-hook

import (
	"syscall"
	"unsafe"
)

// Link between standard type in windows libraries and go types.
type HOOKPROC uintptr
type HINSTANCE uintptr
type HHOOK uintptr
type HWND uintptr
type LRESULT uintptr

// MSG Struct
// https://docs.microsoft.com/fr-fr/windows/win32/api/winuser/ns-winuser-msg?redirectedfrom=MSDN
type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uint32
	LParam  uint32
	Time    uint32
	POINT
}
type POINT struct {
	X int32
	Y int32
}

type BOOL int

var (
	user32Lib, _               = syscall.LoadDLL("user32.dll")
	procCallNextHookEx, _      = user32Lib.FindProc("CallNextHookEx")
	procSetWindowsHookExW, _   = user32Lib.FindProc("SetWindowsHookExW")
	procGetMessageW, _         = user32Lib.FindProc("GetMessageW")
	procUnhookWindowsHookEx, _ = user32Lib.FindProc("UnhookWindowsHookEx")
	procMapVirtualKey, _       = user32Lib.FindProc("MapVirtualKeyW")
	procTranslateMessage, _    = user32Lib.FindProc("TranslateMessage")
	procDispatchMessage, _     = user32Lib.FindProc("DispatchMessage")
	procGetKeyNameTextW, _     = user32Lib.FindProc("GetKeyNameTextW")
)

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-mapvirtualkeyexw
func MapVirtualKey(uCode, uMapType uint32) uint32 {
	r, _, _ := procMapVirtualKey.Call(uintptr(uCode), uintptr(uMapType))
	return uint32(r)
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getkeynametextw
func GetKeyNameText(lParam uint32, lPString *uint16, cchSize uint32) uint32 {
	r, _, _ := procGetKeyNameTextW.Call(
		uintptr(lParam),
		uintptr(unsafe.Pointer(lPString)),
		uintptr(cchSize))

	return uint32(r)
}

// https://docs.microsoft.com/fr-fr/windows/win32/api/winuser/nf-winuser-callnexthookex
func CallNextHookEx(opt, code, wParam, lParam uint64) LRESULT {
	r, _, _ := procCallNextHookEx.Call(
		uintptr(opt),
		uintptr(code),
		uintptr(wParam),
		uintptr(lParam))
	return LRESULT(r)
}

// https://docs.microsoft.com/fr-fr/windows/win32/api/winuser/nf-winuser-setwindowshookexa?redirectedfrom=MSDN
func SetWindowsHookExW(hookId int32, h HOOKPROC, module HINSTANCE, threadId uint32) HHOOK {
	r, _, _ := procSetWindowsHookExW.Call(
		uintptr(hookId),
		uintptr(h),
		uintptr(module),
		uintptr(threadId))
	return HHOOK(r)
}

// https://docs.microsoft.com/fr-fr/windows/win32/api/winuser/nf-winuser-getmessage?redirectedfrom=MSDN
func GetMessageW(message **MSG, hWindow uintptr, wMsgFilterMin, wMsgFilterMax uint32) bool {
	r, _, _ := procGetMessageW.Call(
		uintptr(unsafe.Pointer(message)),
		hWindow,
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMin))
	procTranslateMessage.Call(uintptr(unsafe.Pointer(message)))
	procDispatchMessage.Call(uintptr(unsafe.Pointer(message)))
	if r == 0 {
		return false
	} else {
		return true
	}
}

// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-unhookwindowshookex
func UnhookWindowsHookEx(h HHOOK) bool {
	r, _, _ := procUnhookWindowsHookEx.Call(uintptr(h))
	if r == 0 {
		return false
	} else {
		return true
	}
}
