package utils

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

// Get call method by runtime.Caller
func GetFuncCall(callDepth int) (string, string) {
	pc, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		return "???:0", "???"
	}

	f := runtime.FuncForPC(pc)
	_, function := path.Split(f.Name())
	_, filename := path.Split(file)

	fline := filename + ":" + strconv.Itoa(line)
	return fline, function
}

// 根据app_id生成asset_id
func GenAssetId(appId int64) int64 {
	return int64(GenIdHelp(uint64(appId), 0))
}

// 生成nonce值
func GenNonce() int64 {
	randId1 := GenRandId()
	randId2 := GenRandId()
	content := fmt.Sprintf("%d#%d#%d#%s", randId1, randId2, time.Now().UnixNano(), GetHostName())
	sign := StrSignToInt(content)
	return int64(sign & 0x7FFFFFFFFFFFFFFF)
}

// 生成伪唯一ID
func GenRandId() uint64 {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	randNum1 := rand.Int63()
	randNum2 := rand.Int63()
	shift1 := rand.Intn(16) + 2
	shift2 := rand.Intn(8) + 1

	randId := ((randNum1 >> uint(shift1)) + (randNum2 >> uint(shift2)) + (nano >> 1)) &
		0x7FFFFFFFFFFFFFFF
	return uint64(randId)
}

/**
 * | 0 - 19  	 | 20-31  | 32   | 33 - 40 | 41 - 56   | 57 - 60 | 61-63 |
 * | 20位    	 |  12位  | 1位  | 8位     | 16位      |  4位    |  3位  |
 * | baseId低20位| 随机值 | 标记 | 随机值  | 签名低16位|  随机值 |  0    |
 */
func GenIdHelp(baseId uint64, flag int) uint64 {
	var s, r1, r2, lk uint64
	content := fmt.Sprintf("%d#%d#%d", baseId, flag, time.Now().UnixNano())
	s = StrSignToInt(content)
	r1 = GenRandId()
	r2 = GenRandId()
	lk = baseId

	var id uint64
	id = (lk & 0x0000000000fffff)
	id += ((r2 & 0x000000000000fff0 >> 4) << 20)
	if flag == 1 {
		id += (0x0000000000000001 << 32)
	}
	id += ((r1 & 0x00000000000000ff) << 33)
	id += ((s & 0x000000000000ffff) << 41)
	id += ((r2 & 0x000000000000000f) << 57)

	return id
}

// 对字符串Hash后转化为整数
func StrSignToInt(content string) uint64 {
	h := md5.New()
	io.WriteString(h, content)
	digest := h.Sum(nil)

	var seg1, seg2, seg3, seg4 uint32
	seg1 = binary.LittleEndian.Uint32(digest[0:4])
	seg2 = binary.LittleEndian.Uint32(digest[4:8])
	seg3 = binary.LittleEndian.Uint32(digest[8:12])
	seg4 = binary.LittleEndian.Uint32(digest[12:16])

	var sign, sign1, sign2 uint64
	sign1 = uint64(seg1 + seg3)
	sign2 = uint64(seg2 + seg4)
	sign = (sign1 & 0x00000000ffffffff) | (sign2 << 32)

	return sign
}

func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "127.0.0.1"
	}

	return hostname
}