package main

import (
	"os"
	"os/exec"
	"path"
	"syscall"
	"time"

	"bazil.org/fuse"
)

const (
	maxRetryCount uint32 = 100
	retryGap             = 100 * time.Millisecond
)

func performMount() {
	var (
		curRetryCount                 uint32
		err                           error
		lazyUnmountCmd                *exec.Cmd
		mountPointContainingDirDevice int64
		mountPointDevice              int64
	)

	// err = fuse.Unmount(globals.config.FUSEMountPointPath)
	// if nil != err {
	// 	logTracef("pre-fuse.Unmount() in performMount() returned: %v", err)
	// }

	mountPointContainingDirDevice = fetchInodeDevice("path.Dir([Agent]FUSEMountPointPath", path.Dir(globals.config.FUSEMountPointPath))
	mountPointDevice = fetchInodeDevice("[Agent]FUSEMountPointPath", globals.config.FUSEMountPointPath)

	if mountPointDevice != mountPointContainingDirDevice {
		// Presumably, the mount point is (still) currently mounted, so attempt to unmount it first

		lazyUnmountCmd = exec.Command("fusermount", "-uz", globals.config.FUSEMountPointPath)
		err = lazyUnmountCmd.Run()
		if nil != err {
			logFatal(err)
		}

		curRetryCount = 0

		for mountPointDevice != mountPointContainingDirDevice {
			time.Sleep(retryGap) // Try again in a bit
			curRetryCount++
			if curRetryCount >= maxRetryCount {
				logFatalf("mountPointDevice != mountPointContainingDirDevice MaxRetryCount exceeded")
			}
			mountPointDevice = fetchInodeDevice("[Agent]FUSEMountPointPath", globals.config.FUSEMountPointPath)
		}
	}

	globals.fuseConn, err = fuse.Mount(
		globals.config.FUSEMountPointPath,
		fuse.AllowRoot(), // only one of fuse.AllowOther() or fuse.AllowRoot() allowed
		fuse.AsyncRead(),
		fuse.DefaultPermissions(),
		fuse.ExclCreate(),
		fuse.FSName(globals.config.FUSEVolumeName),
		fuse.NoAppleDouble(),
		fuse.NoAppleXattr(),
		fuse.ReadOnly(),
		fuse.Subtype("ProxyFS"),
		fuse.VolumeName(globals.config.FUSEVolumeName),
	)
	if nil != err {
		logFatal(err)
	}

	<-globals.fuseConn.Ready
	if nil != globals.fuseConn.MountError {
		logFatal(globals.fuseConn.MountError)
	}

	logInfof("%s mounted", globals.config.FUSEMountPointPath)

	go serveFuse()
}

func fetchInodeDevice(pathTitle string, path string) (inodeDevice int64) {
	var (
		err  error
		fi   os.FileInfo
		ok   bool
		stat *syscall.Stat_t
	)

	fi, err = os.Stat(path)
	if nil != err {
		if os.IsNotExist(err) {
			logFatalf("%s path (%s) not found", pathTitle, path)
		} else {
			logFatalf("%s path (%s) os.Stat() failed: %v", pathTitle, path, err)
		}
	}
	if nil == fi.Sys() {
		logFatalf("%s path (%s) had empty os.Stat()", pathTitle, path)
	}
	stat, ok = fi.Sys().(*syscall.Stat_t)
	if !ok {
		logFatalf("%s path (%s) fi.Sys().(*syscall.Stat_t) returned !ok", pathTitle, path)
	}

	inodeDevice = int64(stat.Dev)

	return
}

func performUnmount() {
	var (
		err error
	)

	err = fuse.Unmount(globals.config.FUSEMountPointPath)
	if nil != err {
		logFatal(err)
	}

	logInfof("%s unmounted", globals.config.FUSEMountPointPath)
}
