package goroutineinfo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const _defaultBufferSize = 64 * 1024 // 64 KiB

// Stack is the default
// Goroutine information
type Stack struct {
	id            int
	state         string
	firstFunction string
	fullStack     *bytes.Buffer
	timeStamp     int64
}

func (s Stack) ID() int {
	return s.id
}

func (s Stack) State() string {
	return s.state
}
func (s Stack) FirstFunction() string {
	return s.state
}

func (s Stack) Full() string {
	return s.fullStack.String()
}
func (s Stack) TimeStamp() int64 {
	return s.timeStamp
}

func (s Stack) String() string {
	return fmt.Sprintf(
		"[%d]Goroutine %d in state %s, with %s on top of the stack:\n%s",
		s.timeStamp, s.id, s.state, s.firstFunction, s.Full())
}

func getStackBuffer(all bool) []byte {
	for i := _defaultBufferSize; ; i *= 2 {
		buf := make([]byte, i)
		if n := runtime.Stack(buf, all); n < i {
			return buf[:n]
		}
	}
}

// GetInfo returns goroutine stack information.
// Stack information includes id, state, firstFunction,fullStack and timestamp.
func GetInfo(all bool) []Stack {
	var stacks []Stack

	var curStack *Stack
	stackReader := bufio.NewReader(bytes.NewReader(getStackBuffer(all)))
	for {
		line, err := stackReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic("bufio.NewReader failed on a fixed string")
		}
		// If we see the goroutine header, start a new stack.
		isFirstLine := false
		if strings.HasPrefix(line, "goroutine ") {
			// flush any previous stack
			if curStack != nil {
				stacks = append(stacks, *curStack)
			}
			id, goState := parseGoStackHeader(line)
			curStack = &Stack{
				id:        id,
				state:     goState,
				fullStack: &bytes.Buffer{},
				timeStamp: time.Now().UnixNano(),
			}
			isFirstLine = true
		}
		curStack.fullStack.WriteString(line)
		if !isFirstLine && curStack.firstFunction == "" {
			curStack.firstFunction = parseFirstFunc(line)
		}
	}
	if curStack != nil {
		stacks = append(stacks, *curStack)
	}
	return stacks
}

func parseFirstFunc(line string) string {
	line = strings.TrimSpace(line)
	if idx := strings.LastIndex(line, "("); idx > 0 {
		return line[:idx]
	}
	panic(fmt.Sprintf("function calls missing parents: %q", line))
}

// parseGoStackHeader parses a stack header that looks like:
// goroutine 643 [runnable]:\n
// And returns the goroutine ID, and the state.
func parseGoStackHeader(line string) (goroutineID int, state string) {
	line = strings.TrimSuffix(line, ":\n")
	parts := strings.SplitN(line, " ", 3)
	if len(parts) != 3 {
		panic(fmt.Sprintf("unexpected stack header format: %q", line))
	}

	id, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(fmt.Sprintf("failed to parse goroutine ID: %v in line %q", parts[1], line))
	}

	state = strings.TrimSuffix(strings.TrimPrefix(parts[2], "["), "]")
	return id, state
}
