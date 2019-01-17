package main

import (
	"io"
	"os"
	"reflect"
	"syscall"
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
		switch request.(type) {
		case *fuse.DestroyRequest:
			handleDestroyRequest(request.(*fuse.DestroyRequest))
		case *fuse.GetattrRequest:
			handleGetattrRequest(request.(*fuse.GetattrRequest))
		case *fuse.InitRequest:
			handleInitRequest(request.(*fuse.InitRequest))
		case *fuse.InterruptRequest:
			handleInterruptRequest(request.(*fuse.InterruptRequest))
		case *fuse.LookupRequest:
			handleLookupRequest(request.(*fuse.LookupRequest))
		case *fuse.OpenRequest:
			handleOpenRequest(request.(*fuse.OpenRequest))
		case *fuse.StatfsRequest:
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
				Valid:     time.Duration(10 * time.Second),
				Inode:     uint64(request.Header.Node),
				Size:      uint64(0),
				Blocks:    uint64(0),
				Atime:     time.Now(),
				Mtime:     time.Now(),
				Ctime:     time.Now(),
				Crtime:    time.Now(),
				Mode:      os.ModeDir | syscall.S_IRWXU | syscall.S_IRWXG | syscall.S_IRWXO,
				Nlink:     uint32(2),
				Uid:       uint32(0),
				Gid:       uint32(0),
				Rdev:      uint32(0),
				Flags:     uint32(0),
				BlockSize: uint32(4096),
			},
		}
		logInfof("resonding with:\n%s", utils.JSONify(response, true))
		request.Respond(response)
	} else {
		logInfof("Responding with fuse.ENOTSUP")
		request.RespondError(fuse.ENOTSUP)
	}
}

func handleInitRequest(request *fuse.InitRequest) {
	var (
		response *fuse.InitResponse
	)

	logWarnf("handleInitRequest() should not have been called... fuse.Mount() supposedly took care of it")

	logInfof("TODO: handleInitRequest()")
	logInfof("Header:\n%s", utils.JSONify(request.Header, true))
	logInfof("Payload\n%s", utils.JSONify(request, true))
	response = &fuse.InitResponse{
		Library:      request.Kernel,
		MaxReadahead: request.MaxReadahead,
		Flags:        request.Flags,
		MaxWrite:     uint32(1024 * 1024 * 1024),
	}
	logInfof("resonding with:\n%s", utils.JSONify(response, true))
	request.Respond(response)
}

func handleInterruptRequest(request *fuse.InterruptRequest) {
	logInfof("TODO: handleInterruptRequest()")
	logInfof("Header:\n%s", utils.JSONify(request.Header, true))
	logInfof("Payload\n%s", utils.JSONify(request, true))
	logInfof("Responding with fuse.ENOTSUP")
	request.RespondError(fuse.ENOTSUP)
}

func handleLookupRequest(request *fuse.LookupRequest) {
	logInfof("TODO: handleLookupRequest()")
	logInfof("Header:\n%s", utils.JSONify(request.Header, true))
	logInfof("Payload\n%s", utils.JSONify(request, true))
	logInfof("Responding with fuse.ENOTSUP")
	request.RespondError(fuse.ENOTSUP)
}

func handleOpenRequest(request *fuse.OpenRequest) {
	logInfof("TODO: handleOpenRequest()")
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
		Blocks:  uint64(2 * 1024 * 1024 * 1024),
		Bfree:   uint64(1024 * 1024 * 1024),
		Bavail:  uint64(1024 * 1024 * 1024),
		Files:   uint64(2 * 1024 * 1024),
		Ffree:   uint64(1024 * 1024),
		Bsize:   uint32(4096),
		Namelen: uint32(4096),
		Frsize:  uint32(4096),
	}
	logInfof("resonding with:\n%s", utils.JSONify(response, true))
	request.Respond(response)
}
