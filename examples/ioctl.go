package main

import (
	"github.com/ii64/gouring"
	"golang.org/x/sys/unix"
	"log"
	"unsafe"
)

// See ioctl.h

const (
	IocNrBits           = 8
	IocTypeBits         = 8
	IocSizeBits         = 14
	IocNrShift          = 0
	IocRead     uintptr = 2
	IocWrite    uintptr = 2
)

const (
	IocTypeShift = IocNrShift + IocNrBits
	IocSizeShift = IocTypeShift + IocTypeBits
	IocDirshift  = IocSizeShift + IocSizeBits
)

func IOC(dir, t, nr, size uintptr) uintptr {
	return (dir << IocDirshift) |
		(t << IocTypeShift) |
		(nr << IocNrShift) |
		(size << IocSizeShift)
}

func IOWR(t, nr, size uintptr) uintptr {
	return IOC(IocRead|IocWrite, t, nr, size)
}

func NVME_URING_CMD_IO() uintptr {
	return IOWR('N', 0x80, 32)
}

func main() {

	h, err := gouring.New(256, gouring.IORING_SETUP_IOPOLL|
		gouring.IORING_SETUP_SQE128|gouring.IORING_SETUP_CQE32)
	if err != nil {
		log.Fatal(err)
	}
	defer h.Close()

	fd, err := unix.Open("/dev/nvme0n1", unix.O_RDONLY, 0677)

	sqe := h.GetSqe()
	b := make([]byte, 4096)
	gouring.PrepRead(sqe, fd, &b[0], len(b), 0)
	log.Println("Buffer: ", b)

	sqe.IoUringSqe_Union1.SetCmdOp(NVME_URING_CMD_IO())
	sqe.Opcode = gouring.IORING_OP_URING_CMD

	sqe.UserData.SetUint64(117)

	var cmd gouring.NvmeUringCmd
	cmd.Opcode = gouring.NVME_CMD_READ
	cmd.DataLen = 4096
	cmd.Nsid = 1 // TODO: find nsid
	sqe.Cmd = unsafe.Pointer(&cmd)

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

}
