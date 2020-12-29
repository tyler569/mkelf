package main

import (
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"unsafe"
)

func (c *code) elf(w io.Writer) error {
	const loadBase = 0x100000

	h := elf.Header64{}
	dh := elf.Prog64{}
	ch := elf.Prog64{}
	cb := []byte(c.code.String())
	db := []byte(c.data.String())
	fmt.Printf("% x\n", cb)
	fmt.Printf("% x\n", db)

	const entryOff = uint64(unsafe.Sizeof(h) + 2*unsafe.Sizeof(ch))
	const entry = uint64(loadBase + entryOff)

	h.Ident = [elf.EI_NIDENT]byte{
		0x7f, 'E', 'L', 'F',
		byte(elf.ELFCLASS64),
		byte(elf.ELFDATA2LSB),
		1,
	}
	h.Type = uint16(elf.ET_EXEC)
	h.Machine = uint16(elf.EM_X86_64)
	h.Version = uint32(elf.EV_CURRENT)
	h.Entry = entry
	h.Phoff = uint64(unsafe.Sizeof(h))
	h.Shoff = 0
	h.Flags = 0
	h.Ehsize = uint16(unsafe.Sizeof(h))
	h.Phentsize = uint16(unsafe.Sizeof(ch))
	h.Phnum = 2
	h.Shentsize = uint16(unsafe.Sizeof(elf.Section64{}))
	h.Shnum = 0
	h.Shstrndx = 0

	ch.Type = uint32(elf.PT_LOAD)
	ch.Flags = uint32(elf.PF_R | elf.PF_X)
	ch.Off = entryOff
	ch.Vaddr = entry
	ch.Paddr = 0
	ch.Filesz = uint64(len(cb))
	ch.Memsz = uint64(len(cb))
	ch.Align = 0x1000

	dh.Type = uint32(elf.PT_LOAD)
	dh.Flags = uint32(elf.PF_R | elf.PF_W)
	dh.Off = entryOff + uint64(len(cb))
	dh.Vaddr = (entry+uint64(len(cb))+0xFFF)&^0xFFF + (dh.Off & 0xFFF)
	dh.Paddr = 0
	dh.Filesz = uint64(len(db))
	dh.Memsz = uint64(len(db))
	dh.Align = 0x1000

	err := binary.Write(w, binary.LittleEndian, h)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(w, binary.LittleEndian, ch)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(w, binary.LittleEndian, dh)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(w, binary.LittleEndian, cb)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(w, binary.LittleEndian, db)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
