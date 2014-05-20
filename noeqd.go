/*
Copyright (C) 2011 by Blake Mizerany (@bmizerany)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

/* Much of this code was inspired by a heavily modified NOEQD from Blake Mizerany */

package noeqd

import (
	"fmt"
	"sync"
	"time"
)

const (
	workerIdBits       = uint64(5)
	datacenterIdBits   = uint64(5)
	maxWorkerId        = int64(-1) ^ (int64(-1) << workerIdBits)
	maxDatacenterId    = int64(-1) ^ (int64(-1) << datacenterIdBits)
	sequenceBits       = uint64(12)
	workerIdShift      = sequenceBits
	datacenterIdShift  = sequenceBits + workerIdBits
	timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
	sequenceMask       = int64(-1) ^ (int64(-1) << sequenceBits)

	// Tue, 21 Mar 2006 20:50:14.000 GMT
	twepoch = int64(1288834974657)
)

// Flags
var (
	wid, did, lts int64
)

var (
	mu  sync.Mutex
	seq int64
)

type Generator struct {
}

func NewGenerator(workerid, datacenterid int64) (*Generator, error) {

	if wid < 0 || wid > maxWorkerId {
		err := fmt.Errorf("worker id must be between 0 and %d", maxWorkerId)
		return nil, err
	}

	if did < 0 || did > maxDatacenterId {
		err := fmt.Errorf("datacenter id must be between 0 and %d", maxDatacenterId)
		return nil, err
	}
	wid = workerid
	did = datacenterid
	lts = -1
	return &Generator{}, nil
}

func (g *Generator) Get() (uint64, error) {

	id, err := nextId()
	return uint64(id), err
}

func milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func nextId() (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	ts := milliseconds()

	if ts < lts {
		return 0, fmt.Errorf("time is moving backwards, waiting until %d\n", lts)
	}

	if lts == ts {
		seq = (seq + 1) & sequenceMask
		if seq == 0 {
			for ts <= lts {
				ts = milliseconds()
			}
		}
	} else {
		seq = 0
	}

	lts = ts

	id := ((ts - twepoch) << timestampLeftShift) |
		(did << datacenterIdShift) |
		(wid << workerIdShift) |
		seq

	return id, nil
}
