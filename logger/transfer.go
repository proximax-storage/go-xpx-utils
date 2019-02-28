package logger

import (
	"time"
)

type Sender interface {
	Send(logs []string) error
}

type SenderFunc func(logs []string) error

func (ref SenderFunc) Send(logs []string) error {
	return ref(logs)
}

type Transfer struct {
	receiver chan string
	logPool  logPool
	stop     chan struct{}
	sender   Sender
}

func NewLogsTransfer(chanBufSize int, sender Sender) *Transfer {
	return &Transfer{
		logPool:  *newLogPool(chanBufSize),
		sender:   sender,
		receiver: make(chan string, chanBufSize),
	}
}

func (ref *Transfer) Pool() chan<- string {
	return ref.receiver
}

func (ref *Transfer) Start(delay time.Duration) {
	go func() {
		ref.start(delay)
	}()
}

func (ref *Transfer) Close() error {
	if ref.stop != nil {
		ref.stop <- struct{}{}
		close(ref.stop)
	}

	return nil
}

func (ref *Transfer) start(delay time.Duration) {
	ref.stop = make(chan struct{})

	ticker := time.Tick(delay)

	for {
		select {
		case log := <-ref.receiver:
			ref.logPool.add(log)
		case <-ticker:
			logs := ref.logPool.getAndReset()

			if len(logs) == 0 {
				continue
			}

			if err := ref.sender.Send(logs); err != nil {
				ref.logPool.add(logs...)
			}
		case <-ref.stop:
			var notSentLogs []string

			close(ref.receiver)

			for log := range ref.receiver {
				notSentLogs = append(notSentLogs, log)
			}

			notSentLogs = append(notSentLogs, ref.logPool.getAndReset()...)

			_ = ref.sender.Send(notSentLogs)
			return
		}
	}
}

type logPool struct {
	pool []string
	size int
}

func newLogPool(size int) *logPool {
	return &logPool{
		pool: make([]string, 0),
		size: size,
	}
}

func (ref *logPool) add(logs ...string) {
	if len(ref.pool)+len(logs) <= ref.size {
		ref.pool = append(ref.pool, logs...)
	}
}

func (ref *logPool) getAndReset() []string {
	strings := ref.pool
	ref.pool = make([]string, 0)

	return strings
}
