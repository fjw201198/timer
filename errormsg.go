package timer

import "errors"

const (
	E_QFull  = 0;
	E_QEmpty = 1;
)

var ErrInfo []error = []error{
	errors.New("timer queue full"),
	errors.New("timer queue empty"),
}