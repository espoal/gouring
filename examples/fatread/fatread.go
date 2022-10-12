package main

import (
	"github.com/ii64/gouring"
	"golang.org/x/sys/unix"
	"log"
)

func main() {

	h, err := gouring.New(256, gouring.IORING_SETUP_CQE32)
	if err != nil {
		log.Fatal("error creating:", err)
	}
	defer h.Close()

	fd, err := unix.Open("/tmp/test", unix.O_RDWR, 0677)

	sqe := h.GetSqe()
	b := make([]byte, 20)
	gouring.PrepRead(sqe, fd, &b[0], len(b), 0)
	log.Println("Buffer: ", b)

	submitted, err := h.SubmitAndWait(1)
	if err != nil {
		log.Fatal("Error submitting: ", err)
	}
	println(submitted) // 1

	var cqe *gouring.IoUringCqe
	err = h.WaitCqe(&cqe)
	if err != nil {
		log.Fatal("Error waiting:", err)
	} // check also EINTR

	log.Println("CQE: ", cqe)
	log.Println("Buffer: ", b)
	log.Println("Buffer: ", string(b))

	_ = cqe.UserData
	_ = cqe.Res
	_ = cqe.Flags
}
