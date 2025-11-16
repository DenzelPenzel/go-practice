package main

type ChannelWaitGroup struct {
	addCh  chan int
	doneCh chan struct{}
	waitCh chan struct{}
}

func NewChannelWaitGroup() *ChannelWaitGroup {
	wg := &ChannelWaitGroup{
		addCh:  make(chan int),
		doneCh: make(chan struct{}),
		waitCh: make(chan struct{}),
	}

	return wg
}

func main() {

}
