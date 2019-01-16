package main

import (
	"io"
	"reflect"

	"bazil.org/fuse"
)

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
				logTracef("exiting serveFuse() due to io.EOF")
				return
			}
			logErrorf("serveFuse() exiting due to err: %v", err)
			return
		}
		logTracef("serveFuse() got %v", reflect.ValueOf(req).Type())
		switch reflect.ValueOf(req).Type() {
		case reflect.ValueOf(&fuse.DestroyRequest{}).Type():
			handleDestroyRequest(req.(*fuse.DestroyRequest))
		case reflect.ValueOf(&fuse.GetattrRequest{}).Type():
			handleGetattrRequest(req.(*fuse.GetattrRequest))
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

func handleDestroyRequest(req *fuse.DestroyRequest) {
	logWarnf("TODO: handleDestroyRequest()")
	req.RespondError(fuse.ENOTSUP)
}

func handleGetattrRequest(req *fuse.GetattrRequest) {
	logWarnf("TODO: handleGetattrRequest()")
	req.RespondError(fuse.ENOTSUP)
}

func handleInterruptRequest(req *fuse.InterruptRequest) {
	logWarnf("TODO: handleInterruptRequest()")
	req.RespondError(fuse.ENOTSUP)
}

func handleStatfsRequest(req *fuse.StatfsRequest) {
	logWarnf("TODO: handleStatfsRequest()")
	req.RespondError(fuse.ENOTSUP)
}
