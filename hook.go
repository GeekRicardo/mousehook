package mousehook

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")

	hook HHOOK
)

const (
	WH_MOUSE_LL    = 14
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208
	WM_MOUSEWHEEL  = 0x020A
	WM_XBUTTONDOWN = 0x020B
	WM_XBUTTONUP   = 0x020C
	XBUTTON1       = 0x0001
	XBUTTON2       = 0x0002
)

type POINT struct {
	X, Y int32
}

type MSLLHOOKSTRUCT struct {
	Pt          POINT
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type HHOOK uintptr

type MSG struct {
	HWnd    uintptr
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// MouseButton represents a mouse button
type MouseButton string

// MouseEvent represents a mouse event
type MouseEvent struct {
	Button MouseButton
	X, Y   int32
	Delta  int16 // For wheel events
}

// Mouse buttons
const (
	LeftButton   MouseButton = "LeftButton"
	RightButton  MouseButton = "RightButton"
	MiddleButton MouseButton = "MiddleButton"
	XButton1     MouseButton = "XButton1"
	XButton2     MouseButton = "XButton2"
	Wheel        MouseButton = "Wheel"
)

// Callback functions
var (
	onMouseDown  func(MouseEvent)
	onMouseUp    func(MouseEvent)
	onMouseWheel func(MouseEvent)
)

// SetMouseDownCallback sets the callback function for mouse button down events
func SetMouseDownCallback(callback func(MouseEvent)) {
	onMouseDown = callback
}

// SetMouseUpCallback sets the callback function for mouse button up events
func SetMouseUpCallback(callback func(MouseEvent)) {
	onMouseUp = callback
}

// SetMouseWheelCallback sets the callback function for mouse wheel events
func SetMouseWheelCallback(callback func(MouseEvent)) {
	onMouseWheel = callback
}

// Start starts the mouse hook
func Start() {
	hook = setHook()
	defer unhook()

	var msg MSG
	for {
		ret, _, err := procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(ret) == -1 {
			fmt.Println("GetMessage 错误:", err)
			break
		}
	}
}

func setHook() HHOOK {
	h, _, err := procSetWindowsHookEx.Call(
		uintptr(WH_MOUSE_LL),
		syscall.NewCallback(mouseProc),
		0,
		0,
	)
	if h == 0 {
		fmt.Println("设置钩子失败:", err)
	}
	return HHOOK(h)
}

func unhook() {
	procUnhookWindowsHookEx.Call(uintptr(hook))
}

func mouseProc(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode == 0 {
		mouse := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		event := MouseEvent{
			X: mouse.Pt.X,
			Y: mouse.Pt.Y,
		}

		switch wParam {
		case WM_LBUTTONDOWN:
			event.Button = LeftButton
			if onMouseDown != nil {
				onMouseDown(event)
			}
		case WM_LBUTTONUP:
			event.Button = LeftButton
			if onMouseUp != nil {
				onMouseUp(event)
			}
		case WM_RBUTTONDOWN:
			event.Button = RightButton
			if onMouseDown != nil {
				onMouseDown(event)
			}
		case WM_RBUTTONUP:
			event.Button = RightButton
			if onMouseUp != nil {
				onMouseUp(event)
			}
		case WM_MBUTTONDOWN:
			event.Button = MiddleButton
			if onMouseDown != nil {
				onMouseDown(event)
			}
		case WM_MBUTTONUP:
			event.Button = MiddleButton
			if onMouseUp != nil {
				onMouseUp(event)
			}
		case WM_XBUTTONDOWN:
			if (mouse.MouseData>>16)&XBUTTON1 != 0 {
				event.Button = XButton1
			} else if (mouse.MouseData>>16)&XBUTTON2 != 0 {
				event.Button = XButton2
			}
			if onMouseDown != nil {
				onMouseDown(event)
			}
		case WM_XBUTTONUP:
			if (mouse.MouseData>>16)&XBUTTON1 != 0 {
				event.Button = XButton1
			} else if (mouse.MouseData>>16)&XBUTTON2 != 0 {
				event.Button = XButton2
			}
			if onMouseUp != nil {
				onMouseUp(event)
			}
		case WM_MOUSEWHEEL:
			event.Button = Wheel
			event.Delta = int16(mouse.MouseData >> 16) // High word contains the wheel delta
			if onMouseWheel != nil {
				onMouseWheel(event)
			}
		}
	}

	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}
