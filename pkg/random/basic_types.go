package random

import (
	"math/rand"
)

func String(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

type Range struct {
	Min int64
	Max int64
}

func (r Range) Int() int {
	return int(r.Min + rand.Int63n(r.Max-r.Min+1))
}

func (r Range) String() string {
	return String(int(r.Max))
}

func (r Range) Int8() int8 {
	return int8(r.Int())
}

func (r Range) Int16() int16 {
	return int16(r.Int())
}

func (r Range) Int32() int32 {
	return int32(r.Int())
}

func (r Range) Int64() int64 {
	return int64(r.Int())
}

func (r Range) Uint() uint {
	return uint(r.Int())
}

func (r Range) Uint8() uint8 {
	return uint8(r.Int())
}

func (r Range) Uint16() uint16 {
	return uint16(r.Int())
}

func (r Range) Uint32() uint32 {
	return uint32(r.Int())
}

func (r Range) Uint64() uint64 {
	return uint64(r.Int())
}

func (r Range) Float32() float32 {
	return float32(rand.Float64()*(float64(r.Max-r.Min)) + float64(r.Min))
}

func (r Range) Float64() float64 {
	return rand.Float64()*(float64(r.Max-r.Min)) + float64(r.Min)
}
