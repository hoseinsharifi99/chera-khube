package helper

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"strconv"
	"time"
)

func GenerateTag() string {
	crc32q := crc32.MakeTable(0xD5828281)
	random := rand.Intn(100)
	s := time.Now().String() + strconv.Itoa(random)
	return fmt.Sprintf("%08x", crc32.Checksum([]byte(s), crc32q))
}
