package main

import (
	"io"
	"reflect"

	"bazil.org/fuse"
)

// type fuseError struct {
// 	str   string
// 	errno fuse.Errno
// }

// func (fE *fuseError) Errno() (errno fuse.Errno) {
// 	return fE.errno
// }

// func (fE *fuseError) Error() string {
// 	return fE.str
// }

// func newFuseError(err error) (fE *fuseError) {
// 	return &fuseError{
// 		str:   fmt.Sprintf("%v", err),
// 		errno: fuse.Errno(blunder.Errno(err)),
// 	}
// }

func serveFuse() {
	var (
		err error
		req fuse.Request
	)

	for {
		// Fetch next *fuse.Request... exiting on fuseConn error

		req, err = globals.fuseConn.ReadRequest()
		if nil != err {
			if io.EOF == err {
				logInfof("exiting serveFuse() due to io.EOF")
				return
			}
			logErrorf("serveFuse() exiting due to err: %v", err)
			return
		}
		logInfof("serveFuse() got %#v\n", req)
		logInfof("reflect.ValueOf(req).Type() == %v", reflect.ValueOf(req).Type())
		switch reflect.ValueOf(req).Type() {
		case reflect.ValueOf(&fuse.InterruptRequest{}).Type():
			handleInterruptRequest(req.(*fuse.InterruptRequest))
		case reflect.ValueOf(&fuse.StatfsRequest{}).Type():
			handleStatfsRequest(req.(*fuse.StatfsRequest))
		default:
			logWarnf("recieved unserviced %v", reflect.ValueOf(req).Type())
			req.RespondError(fuse.ENOTSUP)
		}
	}
}

func handleInterruptRequest(req *fuse.InterruptRequest) {
	req.RespondError(fuse.ENOTSUP)
}

func handleStatfsRequest(req *fuse.StatfsRequest) {
	req.RespondError(fuse.ENOTSUP)
}
