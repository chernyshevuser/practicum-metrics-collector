package semaphore

type Semaphore interface {
	Acquire()
	Release()
	Close()
}
