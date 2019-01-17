package main

import (
	"io"
	"os"
	"reflect"
	"time"

	"bazil.org/fuse"
	"github.com/swiftstack/ProxyFS/utils"
)

func serveFuse() {
	var (
		err     error
		request fuse.Request
	)

	for {
		// Fetch next *fuse.Request... exiting on fuseConn error

		request, err = globals.fuseConn.ReadRequest()
		if nil != err {
			if io.EOF == err {
				logTracef("exiting serveFuse() due to io.EOF")
				return
			}
			logErrorf("serveFuse() exiting due to err: %v", err)
			return
		}
		logTracef("serveFuse() got %v", reflect.ValueOf(request).Type())
		switch reflect.ValueOf(request).Type() {
		case reflect.ValueOf(&fuse.DestroyRequest{}).Type():
			handleDestroyRequest(request.(*fuse.DestroyRequest))
		case reflect.ValueOf(&fuse.GetattrRequest{}).Type():
			handleGetattrRequest(request.(*fuse.GetattrRequest))
		case reflect.ValueOf(&fuse.InterruptRequest{}).Type():
			handleInterruptRequest(request.(*fuse.InterruptRequest))
		case reflect.ValueOf(&fuse.StatfsRequest{}).Type():
			handleStatfsRequest(request.(*fuse.StatfsRequest))
		default:
			logWarnf("recieved unserviced %v", reflect.ValueOf(request).Type())
			request.RespondError(fuse.ENOTSUP)
		}
	}
}

func handleDestroyRequest(request *fuse.DestroyRequest) {
	logInfof("TODO: handleDestroyRequest()")
	logInfof("Header:\n%s", utils.JSONify(request.Header, true))
	logInfof("Responding with fuse.ENOTSUP")
	request.RespondError(fuse.ENOTSUP)
}

func handleGetattrRequest(request *fuse.GetattrRequest) {
	var (
		response *fuse.GetattrResponse
	)

	logInfof("TODO: handleGetattrRequest()")
	logInfof("Header:\n%s", utils.JSONify(request.Header, true))
	logInfof("Payload\n%s", utils.JSONify(request, true))
	if fuse.RootID == request.Header.Node {
		response = &fuse.GetattrResponse{
			Attr: fuse.Attr{
				Valid:     time.Duration(0),
				Inode:     uint64(request.Header.Node),
				Size:      uint64(0),
				Blocks:    uint64(0),
				Atime:     time.Now(),
				Mtime:     time.Now(),
				Ctime:     time.Now(),
				Crtime:    time.Now(),
				Mode:      os.FileMode(0777),
				Nlink:     uint32(2),
				Uid:       uint32(0),
				Gid:       uint32(0),
				Rdev:      uint32(0),
				Flags:     uint32(0),
				BlockSize: uint32(0),
			},
		}
		logInfof("resonding with:\n%s", utils.JSONify(response, true))
		request.Respond(response)
	} else {
		logInfof("Responding with fuse.ENOTSUP")
		request.RespondError(fuse.ENOTSUP)
	}
}

func handleInterruptRequest(request *fuse.InterruptRequest) {
	logInfof("TODO: handleInterruptRequest()")
	logInfof("Header:\n%s", utils.JSONify(request.Header, true))
	logInfof("Payload\n%s", utils.JSONify(request, true))
	logInfof("Responding with fuse.ENOTSUP")
	request.RespondError(fuse.ENOTSUP)
}

func handleStatfsRequest(request *fuse.StatfsRequest) {
	var (
		response *fuse.StatfsResponse
	)

	logInfof("TODO: handleStatfsRequest()")
	logInfof("Header:\n%s", utils.JSONify(request.Header, true))
	response = &fuse.StatfsResponse{
		Blocks:  uint64(0),
		Bfree:   uint64(0),
		Bavail:  uint64(0),
		Files:   uint64(0),
		Ffree:   uint64(0),
		Bsize:   uint32(0),
		Namelen: uint32(0),
		Frsize:  uint32(1),
	}
	logInfof("resonding with:\n%s", utils.JSONify(response, true))
	request.Respond(response)
}
