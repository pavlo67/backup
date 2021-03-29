package thread

import (
	"fmt"
	"sync"

	"github.com/pavlo67/common/common/errors"

	"github.com/pavlo67/tools/common/kv"
)

const onNewFIFOKVStrings = "on thread.FIFOKVItems()"

func NewFIFOKVItems(process kv.ItemsProcess) (FIFOKVItems, error) {
	if process == nil {
		return nil, fmt.Errorf(onNewFIFOKVStrings + ": no processKVStrings interface")
	}

	return &fifoKVItems{
		mutex:   &sync.Mutex{},
		process: process,
	}, nil
}

type FIFOKVItems interface {
	FIFOKVItemsAdd
	FIFOKVItemsGetString
}

type FIFOKVItemsAdd interface {
	Add(kvItem kv.Item)
}

type FIFOKVItemsGetString interface {
	GetString() (string, error)
}

// implementation ------------------------------------------------

var _ FIFOKVItems = &fifoKVItems{}

type fifoKVItems struct {
	mutex   *sync.Mutex
	kvQueue []kv.Item
	process kv.ItemsProcess
}

func (fifo *fifoKVItems) Add(kvString kv.Item) {
	fifo.mutex.Lock()
	defer fifo.mutex.Unlock()

	fifo.kvQueue = append(fifo.kvQueue, kvString)
}

func (fifo *fifoKVItems) GetString() (string, error) {
	fifo.mutex.Lock()
	defer fifo.mutex.Unlock()

	var errs []interface{}

	for len(fifo.kvQueue) > 0 {
		value := fifo.kvQueue[0]
		fifo.kvQueue = fifo.kvQueue[1:]
		if err := fifo.process.ProcessOne(value); err != nil {
			errs = append(errs, err)
		}
	}

	resultString, err := fifo.process.Finish()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return resultString, errors.CommonError(errs...)
	}

	return resultString, nil
}
