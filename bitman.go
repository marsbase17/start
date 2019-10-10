package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	data:=[]byte{0x01,0x0b,0x00}
	r:=mooo(data, []int{6,10,8})
	fmt.Println(r)
}

func mooo(data []byte, bits []int) []uint64 {
	var nd []byte
	var res []uint64

	for _,x:=range data {
		s:=reverseBits(x,8)
		nd= append(nd, s)
	}

	br:=New(bytes.NewReader(nd))

	for _, k:= range bits {
		m,_:=br.ReadReverseUint(k)
		res=append(res, m)
	}
	return res
}

type BitReader struct {
	reader io.ByteReader
	byte   byte
	offset byte
}

func New(r io.ByteReader) *BitReader {
	return &BitReader{r, 0, 0}
}

func (r *BitReader) ReadBit() (bool, error) {
	if r.offset == 8 {
		r.offset = 0
	}
	if r.offset == 0 {
		var err error
		if r.byte, err = r.reader.ReadByte(); err != nil {
			return false, err
		}
	}
	bit := (r.byte & (0x80 >> r.offset)) != 0
	r.offset++
	return bit, nil
}

func (r *BitReader) ReadUint(nbits uint) (uint64, error) {
	var result uint64
	for i := nbits - 1; i >= 0; i-- {
		bit, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		if bit {
			result |= 1 << uint(i)
		}
	}
	return result, nil
}

func (r *BitReader) ReadReverseUint(nbits int) (uint64, error) {
	var result uint64
	for i := 0; i < nbits; i++ {
		bit, err := r.ReadBit()
		if err != nil {
			return 0, err
		}
		if bit {
			result |= 1 << uint(i)
		}
	}
	return result, nil
}

func reverseBits(b byte, len int) byte {
	var d byte
	for i:= 0; i < len; i++ {
		d <<= 1
		d |= b & 1
		b >>= 1
	}
	return d
}
