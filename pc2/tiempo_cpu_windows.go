//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

var processTimes = syscall.NewLazyDLL("kernel32.dll").NewProc("GetProcessTimes")

func processCPUSeconds() float64 {
	var created, exited, kernel, user syscall.Filetime
	handle, _ := syscall.GetCurrentProcess()
	result, _, _ := processTimes.Call(uintptr(handle), uintptr(unsafe.Pointer(&created)), uintptr(unsafe.Pointer(&exited)), uintptr(unsafe.Pointer(&kernel)), uintptr(unsafe.Pointer(&user)))
	if result == 0 {
		return 0
	}
	return float64(filetimeTicks(kernel)+filetimeTicks(user)) / 10_000_000
}

func filetimeTicks(ft syscall.Filetime) uint64 {
	return uint64(ft.HighDateTime)<<32 | uint64(ft.LowDateTime)
}
