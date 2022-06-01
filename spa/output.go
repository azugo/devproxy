package spa

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync/atomic"
)

var ansiRegexp = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

func cleanOutput(line string, build bool) string {
	line = ansiRegexp.ReplaceAllString(line, "")
	if !build && strings.HasPrefix(line, "<s> [webpack.Progress]") {
		return ""
	}
	return line
}

func forwardOutput(stdout io.ReadCloser, stderr io.ReadCloser, startRegexp *regexp.Regexp, buildInfo bool) chan struct{} {
	done := make(chan struct{})

	var c int32

	stdoutScanner := bufio.NewScanner(stdout)
	go func() {
		for stdoutScanner.Scan() {
			line := cleanOutput(stdoutScanner.Text(), buildInfo)
			if len(line) == 0 {
				continue
			}
			if c == 0 && startRegexp != nil && startRegexp.MatchString(line) {
				if atomic.CompareAndSwapInt32(&c, 0, 1) {
					done <- struct{}{}
					close(done)
				}
			}
			fmt.Printf("%s\n", line)
		}

		if atomic.CompareAndSwapInt32(&c, 0, 1) {
			done <- struct{}{}
			close(done)
		}
	}()

	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			line := cleanOutput(stderrScanner.Text(), buildInfo)
			if len(line) == 0 {
				continue
			}
			fmt.Printf("ERR: %s\n", line)
		}

		if atomic.CompareAndSwapInt32(&c, 0, 1) {
			done <- struct{}{}
			close(done)
		}
	}()

	return done
}
