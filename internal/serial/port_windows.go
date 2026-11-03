package serial

import (
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// port is a minimal synchronous Win32 COM port, used instead of github.com/tarm/serial.
// tarm wraps its overlapped handle in os.NewFile, which since Go 1.26 registers the handle with the
// runtime's IOCP poller. Every tarm read then strands an IRP in non-paged kernel pool until the handle
// closes: ~400 B per read, ~20 MB/min with the handbrake streaming. That froze the PC after ~10 h uptime.
type port struct {
	h windows.Handle
}

func openPort(name string, baud uint32, readTimeout time.Duration) (*port, error) {
	path, err := windows.UTF16PtrFromString(`\\.\` + name)
	if err != nil {
		return nil, err
	}
	// no FILE_FLAG_OVERLAPPED: one reader per port, and reads return within readTimeout
	h, err := windows.CreateFile(path, windows.GENERIC_READ|windows.GENERIC_WRITE, 0, nil, windows.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, err
	}

	dcb := windows.DCB{
		DCBlength: uint32(unsafe.Sizeof(windows.DCB{})),
		BaudRate:  baud,
		Flags:     0x01 | 0x10, // fBinary, DTR_CONTROL_ENABLE (same as tarm)
		ByteSize:  8,
		Parity:    windows.NOPARITY,
		StopBits:  windows.ONESTOPBIT,
	}
	const maxDword = 0xFFFFFFFF
	// return as soon as any byte is available, or with 0 bytes after readTimeout
	timeouts := windows.CommTimeouts{
		ReadIntervalTimeout:        maxDword,
		ReadTotalTimeoutMultiplier: maxDword,
		ReadTotalTimeoutConstant:   uint32(readTimeout.Milliseconds()),
	}
	if err := windows.SetCommState(h, &dcb); err != nil {
		windows.CloseHandle(h)
		return nil, err
	}
	if err := windows.SetCommTimeouts(h, &timeouts); err != nil {
		windows.CloseHandle(h)
		return nil, err
	}
	return &port{h: h}, nil
}

func (p *port) Read(b []byte) (int, error) {
	var n uint32
	err := windows.ReadFile(p.h, b, &n, nil)
	return int(n), err
}

func (p *port) Close() error {
	return windows.CloseHandle(p.h)
}
