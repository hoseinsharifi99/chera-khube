package helper

import (
	"log"
	"testing"
)

func TestHasher(t *testing.T) {
	pass := GetPassword("121", "_E[*.61P.wA;-M>sS>Gb?@ZVd&>zp\"C=A=R=p+.=K0d7CQ1R\\WX;froI?wY![io`7Hk`<p=:xlK{nA<T[22H,ziR0Nv-BqweuMca+cZ1gR&D<J[6_s<,ko[Yh!aooK~C")
	log.Println(pass)
}
