package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"time"
)

// Flip this when we're ready
const PRODUCTION = true

func dummy(s string) *regexp.Regexp {
	ret, _ := regexp.Compile(s)
	return ret
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, copy the file contents from src to
// dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)",
			sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)",
				dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func moveToChecked(prefix string) string {
	dname := time.Now().Format("2006Jan02_150405")
	src := fmt.Sprintf("/home/ed/spam/%sspammail", prefix)
	dst := fmt.Sprintf("/home/ed/spam/checked/%s%s", prefix, dname)

	// copy to the checked directory
	CopyFile(src, dst)
	if PRODUCTION {
		// empty the source file
		empty := []byte{}
		ioutil.WriteFile(src, empty, 0660)
	}
	return dst
}

func main() {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	spams := []string{"", "list"}
	for _, prefix := range spams {
		loc := moveToChecked(prefix)
		analyzed := Analyze(loc)
		MailOut(analyzed, loc, prefix)
	}
}
