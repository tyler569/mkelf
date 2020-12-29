package main

import (
	"encoding/binary"
	"log"
	"os"
	"strings"
)

type code struct {
	code *strings.Builder
	data *strings.Builder
}

const (
	rax uint8 = iota
	rcx
	rdx
	rbx
	rsp
	rbp
	rsi
	rdi
)

const rexW uint8 = 0x48

func modRM(b, r1, r2 uint8) uint8 {
	return (b << 6) + (r2 << 3) + r1
}

func (c *code) emit(bytes ...byte) {
	c.code.Write(bytes)
}

func (c *code) emitI32(i int32) {
	binary.Write(c.code, binary.LittleEndian, i)
}

func (c *code) add(dst, src uint8) {
	c.emit(rexW, 1, modRM(3, dst, src))
}

func (c *code) ret() {
	c.emit(0xc3)
}

func (c *code) callR32(target int32) {
	c.emit(0xe8)
	c.emitI32(target - 5)
}

func (c *code) jmpR32(target int32) {
	c.emit(0xe9)
	c.emitI32(target - 5)
}

func (c *code) movI32(r byte, i int32) {
	c.emit(0xb8 + r)
	c.emitI32(i)
}

func (c *code) xor(r2, r1 byte) {
	c.emit(0x33, modRM(3, r2, r1))
}

func (c *code) zero(r byte) {
	c.xor(r, r)
}

func (c *code) syscall() {
	c.emit(0x0f, 0x05)
}

func (c *code) offset() int {
	return c.code.Len()
}

func (c *code) emitData(data interface{}) {
	if str, ok := data.(string); ok {
		binary.Write(c.data, binary.LittleEndian, []byte(str))
	} else {
		binary.Write(c.data, binary.LittleEndian, data)
	}
}

func newCode() *code {
	return &code{
		&strings.Builder{},
		&strings.Builder{},
	}
}

const sysExit = 60

func main() {
	c := newCode()

	// c.movI32(rax, sysWrite)
	// c.zero(rdi)
	// c.movI32(rsi, 0)
	// c.symbol(abs4, "hw") // -> c.relocations {
	// 					//      code + 9 - code + 13: absolute address of "hw"
	// 					//	  }
	// c.movI32(rdx, 10)
	// c.syscall()

	c.movI32(rax, sysExit)
	c.movI32(rdi, 100)
	c.syscall()

	// s := c.dataSymbol("hw")
	c.emitData("HelloWorld\n\x00")

	// b := new(strings.Builder)
	// c.elf(b)
	// fmt.Printf("% x\n", b.String())

	f, err := os.OpenFile("out.elf", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	c.elf(f)
}
