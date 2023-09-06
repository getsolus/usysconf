// Copyright © Solus Project
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
	"io"
	"log/slog"
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
		slog.Warn("Failed to access, assuming chroot and continuing", "path", "/proc/1/root")
		return true
	}
	// Check /proc/mounts for overlayfs on "/"
	mounts, err := os.Open("/proc/mounts")
	if err != nil {
		slog.Warn("Failed to access", "path", "/proc/mounts", "reason", err)
		goto FALLBACK
	}
	raw, err = io.ReadAll(mounts)
	if err != nil {
		slog.Warn("Failed to read", "path", "/proc/mounts", "reason", err)
		goto FALLBACK
	}
	_ = mounts.Close()
	if strings.Contains(string(raw), "overlay / overlay") {
		slog.Debug("Overlayfs for '/' found, assuming chroot")
		return true
	}
FALLBACK:
	slog.Debug("Falling back to rigorous check for chroot")
	rootDir, err = os.Stat("/")
	if err != nil {
		slog.Error("Failed to access", "path", "/", "reason", err)
		// TODO: Return error instead of panicking.
		panic("Failed to access")
	}
	pid = os.Getpid()
	chrootDir, err = os.Stat(filepath.Join("/", "proc", strconv.Itoa(pid), "root"))
	if err != nil {
		slog.Error("Failed to access", "path", "/", "reason", err)
		// TODO: Return error instead of panicking.
		panic("Failed to access")
	}
	root = rootDir.Sys().(*syscall.Stat_t)
	chroot = chrootDir.Sys().(*syscall.Stat_t)
	return root.Dev != chroot.Dev || root.Ino != chroot.Ino
}
