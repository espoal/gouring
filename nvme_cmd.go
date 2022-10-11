package gouring

/* struct nvme_uring_cmd {
__u8	opcode;
__u8	flags;
__u16	rsvd1;
__u32	nsid;
__u32	cdw2;
__u32	cdw3;
__u64	metadata;
__u64	addr;
__u32	metadata_len;
__u32	data_len;
__u32	cdw10;
__u32	cdw11;
__u32	cdw12;
__u32	cdw13;
__u32	cdw14;
__u32	cdw15;
__u32	timeout_ms;
__u32   rsvd2;
}*/

const (
	NVME_CMD_WRITE = 0x01
	NVME_CMD_READ  = 0x02
)

type NvmeUringCmd struct {
	Opcode      uint8
	Flags       uint8
	Rsvd1       uint16
	Nsid        uint32
	Cdw2        uint32
	Cdw3        uint32
	Metadata    uint64
	Addr        uint64
	MetadataLen uint32
	DataLen     uint32
	Cdw10       uint32
	Cdw11       uint32
	Cdw12       uint32
	Cdw13       uint32
	Cdw14       uint32
	Cdw15       uint32
	Timeout     uint32
	Rsvd2       uint32
}
