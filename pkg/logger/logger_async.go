// logger write log file
package logger

import (
	"bytes"
	"io"
	"time"

	"go.uber.org/zap/buffer"
)

type writeAsyncer struct {
	p        buffer.Pool
	writer   io.Writer
	ch       chan *buffer.Buffer
	syncChan chan struct{}
}

func newWriteAsyncer(writer io.Writer) *writeAsyncer {
	const logDataChanLen = 20480

	wa := &writeAsyncer{}
	wa.writer = writer
	wa.ch = make(chan *buffer.Buffer, logDataChanLen)
	wa.syncChan = make(chan struct{})
	wa.p = buffer.NewPool()
	go batchWriteLog(wa)
	return wa
}

func (wa *writeAsyncer) Write(data []byte) (int, error) {
	buf := wa.p.Get()
	// 不需要处理返回值
	_, _ = buf.Write(data)
	wa.ch <- buf
	return len(data), nil
}
func (wa *writeAsyncer) Sync() error {
	wa.syncChan <- struct{}{}
	return nil
}

func batchWriteLog(wa *writeAsyncer) {
	buf := bytes.NewBuffer(make([]byte, 0, 10240))
	ticker := time.NewTicker(time.Millisecond * 200)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if len(buf.Bytes()) > 0 {
				_, _ = wa.writer.Write(buf.Bytes())
				buf.Reset()
			}
		case record := <-wa.ch:
			buf.Write(record.Bytes())
			record.Free()
			if len(buf.Bytes()) >= 1024*4 {
				_, _ = wa.writer.Write(buf.Bytes())
				buf.Reset()
			}
		case <-wa.syncChan:
			if len(buf.Bytes()) > 0 {
				_, _ = wa.writer.Write(buf.Bytes())
				buf.Reset()
			}
			break
		}
	}
}
