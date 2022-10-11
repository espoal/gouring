package main

import (
	"github.com/ii64/gouring"
	"golang.org/x/sys/unix"
	"log"
	"unsafe"
)

func main() {

	h, err := gouring.New(256, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	fd, err := unix.Open("/dev/nvme0n1", unix.O_RDONLY, 0677)

	sqe := h.GetSqe()
	b := make([]byte, 20)
	gouring.PrepRW(gouring.IORING_OP_URING_CMD, sqe, fd, unsafe.Pointer(&b[0]), len(b), 0)

	submitted, err := h.SubmitAndWait(1)
	if err != nil {
		log.Fatal(err)
	}
	println(submitted) // 1

	var cqe *gouring.IoUringCqe
	err = h.WaitCqe(&cqe)
	if err != nil {
		log.Fatal(err)
	} // check also EINTR

	log.Println("CQE: ", cqe)
	log.Println("Buffer: ", b)
	log.Println("Buffer: ", string(b))

	_ = cqe.UserData
	_ = cqe.Res
	_ = cqe.Flags
}
