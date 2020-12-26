// Copyright © 2019-2020 Solus Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	log "github.com/DataDrake/waterlog"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// IsChroot detects if the current process is running in a chroot environment
func IsChroot() bool {
	var raw []byte
	var root, chroot *syscall.Stat_t
	var rootDir, chrootDir os.FileInfo
	var pid int
	// Try to check for access to the root partition of PID1 (shell?)
	if _, err := os.Stat("/proc/1/root"); err != nil {
		log.Warnln("Failed to access '/proc/1/root', assuming chroot and continuing.")
		return true
	}
	// Check /proc/mounts for overlayfs on "/"
	mounts, err := os.Open("/proc/mounts")
	if err != nil {
		log.Warnf("Failed to access '/proc/mounts', reason: %s\n", err)
		goto FALLBACK
	}
	raw, err = ioutil.ReadAll(mounts)
	if err != nil {
		log.Warnf("Failed to read '/proc/mounts', reason: %s\n", err)
		goto FALLBACK
	}
	_ = mounts.Close()
	if strings.Contains(string(raw), "overlay / overlay") {
		log.Debugln("Overlayfs for '/' found, assuming chroot.\n")
		return true
	}
FALLBACK:
	log.Debugln("Falling back to rigorous check for chroot")
	rootDir, err = os.Stat("/")
	if err != nil {
		log.Fatalf("Failed to access '/', reason: %s\n", err)
	}
	pid = os.Getpid()
	chrootDir, err = os.Stat(filepath.Join("/", "proc", strconv.Itoa(pid), "root"))
	if err != nil {
		log.Fatalf("Failed to access '/', reason: %s\n", err)
	}
	root = rootDir.Sys().(*syscall.Stat_t)
	chroot = chrootDir.Sys().(*syscall.Stat_t)
	return root.Dev == chroot.Dev && root.Ino == chroot.Ino
}
