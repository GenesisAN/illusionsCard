package util

// 代码来源：github.com/qjfoidnh/BaiduPCS-Go/pcsutil/checksum

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

const (
	// B byte
	B = (int64)(1 << (10 * iota))
	// KB kilobyte
	KB
	// MB megabyte
	MB
	// GB gigabyte
	GB
	// TB terabyte
	TB
	// PB petabyte
	PB
)

const (
	// DefaultBufSize 默认的bufSize
	DefaultBufSize = int(256 * KB)
	// CHECKSUM_MD5 获取文件的 md5 值
	CHECKSUM_MD5 int = 1 << iota
	// CHECKSUM_SLICE_MD5 获取文件前 sliceSize 切片的 md5 值
	CHECKSUM_SLICE_MD5
	// CHECKSUM_CRC32 获取文件的 crc32 值
	CHECKSUM_CRC32
	SliceMD5Size = 256 * KB
)

type (
	ChecksumWriter interface {
		io.Writer
		Sum() interface{}
	}

	ChecksumWriteUnit struct {
		SliceEnd       int64
		End            int64
		SliceSum       interface{}
		Sum            interface{}
		OnlySliceSum   bool
		ChecksumWriter ChecksumWriter

		ptr int64
	}

	hashChecksumWriter struct {
		h hash.Hash
	}

	hash32ChecksumWriter struct {
		h hash.Hash32
	}

	// LocalFileMeta 本地文件元信息
	LocalFileMeta struct {
		Path     string `json:"path"`     // 本地路径
		Length   int64  `json:"length"`   // 文件大小
		SliceMD5 []byte `json:"slicemd5"` // 文件前 requiredSliceLen (256KB) 切片的 md5 值
		MD5      []byte `json:"md5"`      // 文件的 md5
		CRC32    uint32 `json:"crc32"`    // 文件的 crc32
		ModTime  int64  `json:"modtime"`  // 修改日期
	}

	// LocalFileChecksum 校验本地文件
	LocalFileChecksum struct {
		LocalFileMeta
		bufSize   int
		sliceSize int
		buf       []byte
		file      *os.File // 文件
	}
)

// GetFileBDMD5 获取文件的大小, md5, 前256KB切片的 md5, crc32
func GetFileBDMD5(localPath string) (bdmd5 string, err error) {
	lfc := NewLocalFileChecksum(localPath, int(SliceMD5Size))
	defer lfc.Close()

	err = lfc.OpenPath()
	if err != nil {
		return "", err
	}
	err = lfc.Sum(CHECKSUM_MD5 | CHECKSUM_SLICE_MD5 | CHECKSUM_CRC32)
	if err != nil {
		return "", err
	}
	strLength, strMd5, strSliceMd5 := strconv.FormatInt(lfc.Length, 10), hex.EncodeToString(lfc.MD5), hex.EncodeToString(lfc.SliceMD5)
	fileName := filepath.Base(localPath)
	regFileName := strings.Replace(fileName, " ", "_", -1)
	regFileName = strings.Replace(regFileName, "#", "_", -1)
	return strMd5 + "#" + strSliceMd5 + "#" + strLength + "#" + regFileName, nil
}

func NewLocalFileChecksum(localPath string, sliceSize int) *LocalFileChecksum {
	return NewLocalFileChecksumWithBufSize(localPath, DefaultBufSize, sliceSize)
}

func NewLocalFileChecksumWithBufSize(localPath string, bufSize, sliceSize int) *LocalFileChecksum {
	return &LocalFileChecksum{
		LocalFileMeta: LocalFileMeta{
			Path: localPath,
		},
		bufSize:   bufSize,
		sliceSize: sliceSize,
	}
}

func (lfc *LocalFileChecksum) createChecksumWriteUnit(cw ChecksumWriter, isAll, isSlice bool, getSumFunc func(sliceSum interface{}, sum interface{})) (wu *ChecksumWriteUnit, deferFunc func(err error)) {
	wu = &ChecksumWriteUnit{
		ChecksumWriter: cw,
		End:            lfc.LocalFileMeta.Length,
		OnlySliceSum:   !isAll,
	}

	if isSlice {
		wu.SliceEnd = int64(lfc.sliceSize)
	}

	return wu, func(err error) {
		if err != nil {
			return
		}
		getSumFunc(wu.SliceSum, wu.Sum)
	}
}

// Sum 计算文件摘要值
func (lfc *LocalFileChecksum) Sum(checkSumFlag int) (err error) {
	lfc.fix()
	wus := make([]*ChecksumWriteUnit, 0, 2)
	if (checkSumFlag & (CHECKSUM_MD5 | CHECKSUM_SLICE_MD5)) != 0 {
		md5w := md5.New()
		wu, d := lfc.createChecksumWriteUnit(
			NewHashChecksumWriter(md5w),
			(checkSumFlag&CHECKSUM_MD5) != 0,
			(checkSumFlag&CHECKSUM_SLICE_MD5) != 0,
			func(sliceSum interface{}, sum interface{}) {
				if sliceSum != nil {
					lfc.SliceMD5 = sliceSum.([]byte)
				}
				if sum != nil {
					lfc.MD5 = sum.([]byte)
				}
			},
		)

		wus = append(wus, wu)
		defer d(err)
	}
	if (checkSumFlag & CHECKSUM_CRC32) != 0 {
		crc32w := crc32.NewIEEE()
		wu, d := lfc.createChecksumWriteUnit(
			NewHash32ChecksumWriter(crc32w),
			true,
			false,
			func(sliceSum interface{}, sum interface{}) {
				if sum != nil {
					lfc.CRC32 = sum.(uint32)
				}
			},
		)

		wus = append(wus, wu)
		defer d(err)
	}

	err = lfc.repeatRead(wus...)
	return
}

func NewHashChecksumWriter(h hash.Hash) ChecksumWriter {
	return &hashChecksumWriter{
		h: h,
	}
}

func (hc *hashChecksumWriter) Write(p []byte) (n int, err error) {
	return hc.h.Write(p)
}

func (hc *hashChecksumWriter) Sum() interface{} {
	return hc.h.Sum(nil)
}
func NewHash32ChecksumWriter(h32 hash.Hash32) ChecksumWriter {
	return &hash32ChecksumWriter{
		h: h32,
	}
}
func (hc *hash32ChecksumWriter) Write(p []byte) (n int, err error) {
	return hc.h.Write(p)
}

func (hc *hash32ChecksumWriter) Sum() interface{} {
	return hc.h.Sum32()
}

func (lfc *LocalFileChecksum) initBuf() {
	if lfc.buf == nil {
		lfc.buf = RawMallocByteSlice(lfc.bufSize)
	}
}

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ uintptr, needzero bool) unsafe.Pointer

func RawMallocByteSlice(size int) []byte {
	p := mallocgc(uintptr(size), 0, false)
	b := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(p),
		Len:  size,
		Cap:  size,
	}))
	return b
}
func (lfc *LocalFileChecksum) repeatRead(wus ...*ChecksumWriteUnit) (err error) {
	if lfc.file == nil {
		return errors.New("file is nil")
	}

	lfc.initBuf()

	defer func() {
		_, err = lfc.file.Seek(0, os.SEEK_SET) // 恢复文件指针
		if err != nil {
			return
		}
	}()

	// 读文件
	var (
		n int
	)
read:
	for {
		n, err = lfc.file.Read(lfc.buf)
		switch err {
		case io.EOF:
			err = lfc.writeChecksum(lfc.buf[:n], wus...)
			break read
		case nil:
			err = lfc.writeChecksum(lfc.buf[:n], wus...)
		default:
			return
		}
	}
	switch err {
	case ErrChecksumWriteAllStop: // 全部结束
		err = nil
	}
	return
}

var (
	ErrFileIsNil            = errors.New("file is nil")
	ErrChecksumWriteStop    = errors.New("checksum write stop")
	ErrChecksumWriteAllStop = errors.New("checksum write all stop")
)

func (wi *ChecksumWriteUnit) handleEnd() error {
	if wi.ptr >= wi.End {
		// 已写完
		if !wi.OnlySliceSum {
			wi.Sum = wi.ChecksumWriter.Sum()
		}
		return ErrChecksumWriteStop
	}
	return nil
}
func (wi *ChecksumWriteUnit) write(p []byte) (n int, err error) {
	if wi.End <= 0 {
		// do nothing
		err = ErrChecksumWriteStop
		return
	}
	err = wi.handleEnd()
	if err != nil {
		return
	}

	var (
		i    int
		left = wi.End - wi.ptr
		lenP = len(p)
	)
	if left < int64(lenP) {
		// 读取即将完毕
		i = int(left)
	} else {
		i = lenP
	}
	n, err = wi.ChecksumWriter.Write(p[:i])
	if err != nil {
		return
	}
	wi.ptr += int64(n)
	if left < int64(lenP) {
		err = wi.handleEnd()
		return
	}
	return
}
func (wi *ChecksumWriteUnit) Write(p []byte) (n int, err error) {
	if wi.SliceEnd <= 0 { // 忽略Slice
		// 读取全部
		n, err = wi.write(p)
		return
	}

	// 要计算Slice的情况
	// 调整slice
	if wi.SliceEnd > wi.End {
		wi.SliceEnd = wi.End
	}

	// 计算剩余Slice
	var (
		sliceLeft = wi.SliceEnd - wi.ptr
	)
	if sliceLeft <= 0 {
		// 已处理完Slice
		if wi.OnlySliceSum {
			err = ErrChecksumWriteStop
			return
		}

		// 继续处理
		n, err = wi.write(p)
		return
	}

	var (
		lenP = len(p)
	)
	if sliceLeft <= int64(lenP) {
		var n1, n2 int
		n1, err = wi.write(p[:sliceLeft])
		n += n1
		if err != nil {
			return
		}
		wi.SliceSum = wi.ChecksumWriter.Sum().([]byte)
		n2, err = wi.write(p[sliceLeft:])
		n += n2
		if err != nil {
			return
		}
		return
	}
	n, err = wi.write(p)
	return
}
func (lfc *LocalFileChecksum) writeChecksum(data []byte, wus ...*ChecksumWriteUnit) (err error) {
	doneCount := 0
	for _, wu := range wus {
		_, err := wu.Write(data)
		switch err {
		case ErrChecksumWriteStop:
			doneCount++
			continue
		case nil:
		default:
			return err
		}
	}
	if doneCount == len(wus) {
		return ErrChecksumWriteAllStop
	}
	return nil
}
func (lfc *LocalFileChecksum) fix() {
	if lfc.sliceSize <= 0 {
		lfc.sliceSize = DefaultBufSize
	}
	if lfc.bufSize < DefaultBufSize {
		lfc.bufSize = DefaultBufSize
	}
}

// OpenPath 检查文件状态并获取文件的大小 (Length)
func (lfc *LocalFileChecksum) OpenPath() error {
	if lfc.file != nil {
		lfc.file.Close()
	}

	var err error
	lfc.file, err = os.Open(lfc.Path)
	if err != nil {
		return err
	}

	info, err := lfc.file.Stat()
	if err != nil {
		return err
	}

	lfc.Length = info.Size()
	lfc.ModTime = info.ModTime().Unix()
	return nil
}

// GetFile 获取文件
func (lfc *LocalFileChecksum) GetFile() *os.File {
	return lfc.file
}

// Close 关闭文件
func (lfc *LocalFileChecksum) Close() error {
	if lfc.file == nil {
		return errors.New("file is nil")
	}

	return lfc.file.Close()
}
