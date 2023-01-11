package id

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	machineID = genMachineID()
	counter   = genCounter()
)

func genCounter() uint32 {
	return uint32(time.Now().Unix())
}

func genMachineID() [3]byte {
	var sum [3]byte
	if hostname, err := os.Hostname(); err != nil {
		n := uint32(time.Now().Unix())
		sum[0] = byte(n >> 0)
		sum[1] = byte(n >> 8)
		sum[2] = byte(n >> 16)
	} else {
		hw := md5.New()
		if _, err = hw.Write([]byte(hostname)); err != nil {
			panic(err)
		}
		copy(sum[:], hw.Sum(nil))
	}

	offset := time.Now().UnixNano() % 3
	sum[0] = sum[0] | (0b1001 << offset)
	sum[0] = sum[0] | (1 << offset)
	sum[0] = sum[0] & 0b00111111
	fmt.Println("machine id", "machineID", fmt.Sprintf("%.8b %.8b %.8b", sum[0], sum[1], sum[2]))
	return sum
}

func high() []byte {
	var a [8]byte
	a[0] = 0
	a[1] = 0b1100 // a base value

	now := uint32(time.Now().Unix())
	a[2] = (byte)(now >> 24 & 0xff)
	a[3] = (byte)(now >> 16 & 0xff)
	a[4] = (byte)(now >> 8 & 0xff)
	a[5] = (byte)(now & 0xff)

	i := atomic.AddUint32(&counter, 1)
	a[6] = (byte)((i >> 8) & 0xff)
	a[7] = (byte)(i & 0xff)
	return a[:]
}

func low() []byte {
	var a [8]byte
	a[0] = 0
	a[1] = machineID[0] & 0b111
	a[2] = machineID[1]
	a[3] = machineID[2]

	jitter := rand.Int()
	a[4] = byte(jitter >> 24)
	a[5] = byte(jitter >> 16)
	a[6] = byte(jitter >> 8)
	a[7] = byte(jitter)
	return a[:]
}

// generate a 20 bytes uuid ,you can add a prefix if necessary
func NewWithPrefix(prefix string) string {
	b := strings.Builder{}
	b.Grow(len(prefix) + 20)
	high := strconv.FormatUint(binary.BigEndian.Uint64(high()), 36)
	low := strconv.FormatUint(binary.BigEndian.Uint64(low()), 36)

	if len(high)+len(low) != 20 {
		fmt.Println("id length != 20", "low", low, "high", high, "id", prefix+high+low)
	}

	b.WriteString(prefix)
	b.WriteString(high)
	b.WriteString(low)
	return b.String()
}
