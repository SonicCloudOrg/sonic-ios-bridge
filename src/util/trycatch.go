package util

import (
	"reflect"
)

type TryCatch struct {
	errChan      chan interface{}
	catches      map[reflect.Type]func(err error)
	defaultCatch func(err error)
} // end type

func (t TryCatch) Try(block func()) TryCatch {
	t.errChan = make(chan interface{})
	t.catches = map[reflect.Type]func(err error){}
	t.defaultCatch = func(err error) {}
	go func() {
		defer func() {
			t.errChan <- recover()
		}()
		block()
	}()
	return t
} // end Try()

func (t TryCatch) CatchAll(block func(err error)) TryCatch {
	t.defaultCatch = block
	return t
} // end CatchAll()

func (t TryCatch) Catch(e error, block func(err error)) TryCatch {
	errorType := reflect.TypeOf(e)
	t.catches[errorType] = block
	return t
} // end Catch()

func (t TryCatch) Finally(block func()) TryCatch {
	err := <-t.errChan
	if err != nil {
		catch := t.catches[reflect.TypeOf(err)]
		if catch != nil {
			catch(err.(error))
		} else {
			t.defaultCatch(err.(error))
		} // end if
	} // end if
	block()
	return t
} // end Finally()

/*
References:
https://xiaorui.cc/archives/4674
*/
