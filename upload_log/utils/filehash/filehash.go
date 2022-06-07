package filehash

// ————————————————
// 版权声明：本文为CSDN博主「Grassto」的原创文章，遵循CC 4.0 BY-SA版权协议，转载请附上原文出处链接及本声明。
// 原文链接：https://blog.csdn.net/DisMisPres/article/details/120545195

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"sync"
)

//#region HashWriter 用 channel 做各哈希的计算
type HashWriter struct {
	md5W    hash.Hash
	md5Chan chan []byte

	sha1W    hash.Hash
	sha1Chan chan []byte

	sha256W    hash.Hash
	sha256Chan chan []byte

	onceClose *sync.Once // 只关闭一次

	wg *sync.WaitGroup
}

func NewHashWriter(useMd5, useSha1, useSha256 bool) *HashWriter {
	writer := new(HashWriter)
	writer.onceClose = new(sync.Once)
	writer.wg = new(sync.WaitGroup)

	var hashCount int
	chanCount := 30

	if useMd5 {
		hashCount++
		writer.md5W = md5.New()
		writer.md5Chan = make(chan []byte, chanCount)
		go writer.writeMd5()
	}

	if useSha1 {
		hashCount++
		writer.sha1W = sha1.New()
		writer.sha1Chan = make(chan []byte, chanCount)
		go writer.writeSha1()
	}

	if useSha256 {
		hashCount++
		writer.sha256W = sha256.New()
		writer.sha256Chan = make(chan []byte, chanCount)
		go writer.writeSha256()
	}

	writer.wg.Add(hashCount)

	return writer
}

func (this *HashWriter) writeMd5() {
	doWriteHash(this.md5W, this.md5Chan)
	// fmt.Println("md5 exit")
	this.wg.Done()
}
func (this *HashWriter) writeSha1() {
	doWriteHash(this.sha1W, this.sha1Chan)
	// fmt.Println("sha1 exit")
	this.wg.Done()
}
func (this *HashWriter) writeSha256() {
	doWriteHash(this.sha256W, this.sha256Chan)
	// fmt.Println("sha256 exit")
	this.wg.Done()
}
func doWriteHash(writer hash.Hash, in chan []byte) {
	for {
		select {
		case buf, open := <-in:
			if !open {
				return
			}
			writer.Write(buf)
		}
	}
}

func (this *HashWriter) Write(buf []byte) {
	// 注意这里的channel传递，[]byte的引用传值的问题
	tmpBuf := make([]byte, len(buf))
	copy(tmpBuf, buf)
	if this.md5Chan != nil {
		this.md5Chan <- tmpBuf
	}
	if this.sha1Chan != nil {
		this.sha1Chan <- tmpBuf
	}
	if this.sha256Chan != nil {
		this.sha256Chan <- tmpBuf
	}
}

func (this *HashWriter) Close() {
	this.onceClose.Do(func() {
		if this.md5Chan != nil {
			close(this.md5Chan)
		}
		if this.sha1Chan != nil {
			close(this.sha1Chan)
		}
		if this.sha256Chan != nil {
			close(this.sha256Chan)
		}
	})
}

func (this *HashWriter) Sum(b []byte) (md5Sum []byte, sha1Sum []byte, sha256Sum []byte) {
	this.Close()

	// 需确保数据都已写完
	this.wg.Wait()

	if this.md5W != nil {
		md5Sum = this.md5W.Sum(b)
	}
	if this.sha1W != nil {
		sha1Sum = this.sha1W.Sum(b)
	}
	if this.sha256W != nil {
		sha256Sum = this.sha256W.Sum(b)
	}

	return
}

//#endregion
