# KBHook

This repo provide low level keylogger for windows systems.
it use hooking on user32.dll

HookWin.go provides main functions to general hooking on user32.dll

KBHook.go is the keylogger. You can use different commands through a chan string to manage it :
- startLog. Start the keylogger
- getLog. Get the log of the keylogger
- stopLog. Stop the keylogger





# Sources
- https://docs.microsoft.com/fr-fr/windows/win32/winmsg/about-hooks?redirectedfrom=MSDN#wh_msgfilter_wh_sysmsgfilterhooks
- https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms644985(v=vs.85)?redirectedfrom=MSDN
- https://en.wikipedia.org/wiki/Hooking
- https://causeyourestuck.io/2016/04/20/keyboard-hook-win32api-2/
- https://github.com/moutend/go-hook