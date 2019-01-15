package main

import (
	"fmt"
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

	fE := fuse.Unmount(globals.config.FUSEMountPointPath)
	fmt.Printf("UNDO: In performMount(), fuse.Unmount(\"%s\") returned %v\n", globals.config.FUSEMountPointPath, fE)

	mountPointContainingDirDevice = fetchInodeDevice("path.Dir([Agent]FUSEMountPointPath", path.Dir(globals.config.FUSEMountPointPath))
	mountPointDevice = fetchInodeDevice("[Agent]FUSEMountPointPath", globals.config.FUSEMountPointPath)

	if mountPointDevice != mountPointContainingDirDevice {
		// Presumably, the mount point is (still) currently mounted, so attempt to unmount it first

		lazyUnmountCmd = exec.Command("fusermount", "-uz", globals.config.FUSEMountPointPath)
		err = lazyUnmountCmd.Run()
		if nil != err {
			return
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

	logInfof("TODO: serve the mount point")

	globals.fuseConn, err = fuse.Mount(
		globals.config.FUSEMountPointPath,
		fuse.AllowOther(),
		fuse.AsyncRead(),
		fuse.NoAppleDouble(),
		fuse.NoAppleXattr(),
		fuse.ReadOnly(),
	)

	fmt.Printf("UNDO: fuse.Mount() returned %#v (err: %v)\n", globals.fuseConn, err)

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
	logInfof("TODO: unserve the mount point")
	fE := fuse.Unmount(globals.config.FUSEMountPointPath)
	fmt.Printf("UNDO: In performUnmount(), fuse.Unmount(\"%s\") returned %v\n", globals.config.FUSEMountPointPath, fE)
}

/*
func performMount(volume *volumeStruct) (err error) {
	var (
		conn                          *fuselib.Conn
		curRetryCount                 uint32
		lazyUnmountCmd                *exec.Cmd
		missing                       bool
		mountHandle                   fs.MountHandle
		mountPointContainingDirDevice int64
		mountPointDevice              int64
	)

	volume.mounted = false

	missing, mountPointContainingDirDevice, err = fetchInodeDevice(path.Dir(volume.mountPointName))
	if nil != err {
		return
	}
	if missing {
		logger.Infof("Unable to serve %s.FUSEMountPoint == %s (mount point dir's parent does not exist)", volume.volumeName, volume.mountPointName)
		return
	}

	missing, mountPointDevice, err = fetchInodeDevice(volume.mountPointName)
	if nil == err {
		if missing {
			logger.Infof("Unable to serve %s.FUSEMountPoint == %s (mount point dir does not exist)", volume.volumeName, volume.mountPointName)
			return
		}
	}

	if (nil != err) || (mountPointDevice != mountPointContainingDirDevice) {
		// Presumably, the mount point is (still) currently mounted, so attempt to unmount it first

		lazyUnmountCmd = exec.Command("fusermount", "-uz", volume.mountPointName)
		err = lazyUnmountCmd.Run()
		if nil != err {
			return
		}

		curRetryCount = 0

		for {
			time.Sleep(retryGap) // Try again in a bit
			missing, mountPointDevice, err = fetchInodeDevice(volume.mountPointName)
			if nil == err {
				if missing {
					err = fmt.Errorf("Race condition: %s.FUSEMountPoint == %s disappeared [case 1]", volume.volumeName, volume.mountPointName)
					return
				}
				if mountPointDevice == mountPointContainingDirDevice {
					break
				}
			}
			curRetryCount++
			if curRetryCount >= maxRetryCount {
				err = fmt.Errorf("MaxRetryCount exceeded for %s.FUSEMountPoint == %s [case 1]", volume.volumeName, volume.mountPointName)
				return
			}
		}
	}

	conn, err = fuselib.Mount(
		volume.mountPointName,
		fuselib.FSName(volume.mountPointName),
		fuselib.AllowOther(),
		// OS X specific—
		fuselib.LocalVolume(),
		fuselib.VolumeName(volume.mountPointName),
	)

	if nil != err {
		logger.WarnfWithError(err, "Couldn't mount %s.FUSEMountPoint == %s", volume.volumeName, volume.mountPointName)
		err = nil
		return
	}

	mountHandle, err = fs.Mount(volume.volumeName, fs.MountOptions(0))
	if nil != err {
		return
	}

	fs := &ProxyFUSE{mountHandle: mountHandle}

	// We synchronize the mounting of the mount point to make sure our FUSE goroutine
	// has reached the point that it can service requests.
	//
	// Otherwise, if proxyfsd is killed after we block on a FUSE request but before our
	// FUSE goroutine has had a chance to run we end up with an unkillable proxyfsd process.
	//
	// This would result in a "proxyfsd <defunct>" process that is only cleared by rebooting
	// the system.
	fs.wg.Add(1)

	go func(mountPointName string, conn *fuselib.Conn) {
		defer conn.Close()
		fusefslib.Serve(conn, fs)
	}(volume.mountPointName, conn)

	// Wait for FUSE to mount the file system.   The "fs.wg.Done()" is in the
	// Root() routine.
	fs.wg.Wait()

	// If we made it to here, all was ok

	logger.Infof("Now serving %s.FUSEMountPoint == %s", volume.volumeName, volume.mountPointName)

	volume.mounted = true

	err = nil
	return
}

func (dummy *globalsStruct) UnserveVolume(confMap conf.ConfMap, volumeName string) (err error) {
	var (
		lazyUnmountCmd *exec.Cmd
		ok             bool
		volume         *volumeStruct
	)

	err = nil // default return

	volume, ok = globals.volumeMap[volumeName]

	if ok {
		if volume.mounted {
			err = fuselib.Unmount(volume.mountPointName)
			if nil == err {
				logger.Infof("Unmounted %v", volume.mountPointName)
			} else {
				lazyUnmountCmd = exec.Command("fusermount", "-uz", volume.mountPointName)
				err = lazyUnmountCmd.Run()
				if nil == err {
					logger.Infof("Lazily unmounted %v", volume.mountPointName)
				} else {
					logger.Infof("Unable to lazily unmount %v - got err == %v", volume.mountPointName, err)
				}
			}
		}

		delete(globals.volumeMap, volume.volumeName)
		delete(globals.mountPointMap, volume.mountPointName)
	}

	return // return err as set from above
}
*/
