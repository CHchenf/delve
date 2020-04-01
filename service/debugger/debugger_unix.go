// +build !windows

package debugger

import (
	"debug/elf"
	"debug/macho"
	"os"
	"runtime"

	"github.com/go-delve/delve/pkg/proc"
)

func verifyBinaryFormat(exePath string) error {
	f, err := os.Open(exePath)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	if (fi.Mode() & 0111) == 0 {
		return proc.ErrNotExecutable
	}

	// check that the binary format is what we expect for the host system
	switch runtime.GOOS {
	case "darwin":
		_, err = macho.NewFile(f)
	case "linux", "freebsd":
		_, err = elf.NewFile(f)
	default:
		panic("attempting to open file Delve cannot parse")
	}
	if err != nil {
		return proc.ErrNotExecutable
	}
	return nil
}
