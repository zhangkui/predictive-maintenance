package util

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"
)

var sequence uint64

func NextID() uint64 {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return uint64(time.Now().UnixNano()) ^ binary.LittleEndian.Uint64(b[:]) ^ atomic.AddUint64(&sequence, 1)
}
func TicketCode(id uint64) string { return fmt.Sprintf("PM-%d-%06d", time.Now().Year(), id%1000000) }
func TraceID() string             { return fmt.Sprintf("%d-%06d", time.Now().UnixNano(), NextID()%1000000) }
