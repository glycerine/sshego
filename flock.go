package gosshtun

import (
	"log"
	"os"
	"syscall"
)

func flock(file *os.File) {
	log.Println("getting flock on ", file.Name())
	syscall.Flock(int(file.Fd()), syscall.LOCK_SH)
}

func funlock(file *os.File) {
	log.Println("releasing flock on ", file.Name())
	syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}
