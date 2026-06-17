// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package runtime

import (
	"fmt"
	"sync"

	"ballerina-lang-go/platform/pal"
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/runtime/internal/exec"
	"ballerina-lang-go/values"
)

// State is the runtime lifecycle state.
type State uint32

const (
	StateUninitialized State = iota
	StateInitializing
	StateListening
	StateGracefulStopping
	StateImmediateStopping
	StateStopped
)

const numStates = int(StateStopped) + 1

func (s State) String() string {
	switch s {
	case StateUninitialized:
		return "uninitialized"
	case StateInitializing:
		return "initializing"
	case StateListening:
		return "listening"
	case StateGracefulStopping:
		return "gracefulStopping"
	case StateImmediateStopping:
		return "immediateStopping"
	case StateStopped:
		return "stopped"
	default:
		panic("unexpected")
	}
}

// lifeCycle bundles the lifecycle state private to a Runtime.
// Current spec don't fully describe the how the lifecycle management happens. Therefore this implementation makes following assumptions/design decisions
// 1. Signals are triggered by humans (we are not dealing with repeat signals coming within micro-seconds)
// 2. We will not handle signals until start has been completed (IMPORTANT: we are not dropping the signal we just don't execute them until start is over)
//   - this is to simplify the design
//
// 3. We listen to signal only after starting initializing phase. Also we don't support concurrent init so calls to init should be sequential
//   - Sending a stop signal during initialization will cause runtime to panic (we don't have listeners started so nothing to stop there, winding down any user strands started by init/main is ignored)
//   - In listening phase sending a kill signal will lead to a graceful shutdown from the runtime and will write the exit code to ExitStatus channel
//
// 4. We try to do best effort escalation on repeated graceful stop signals
//   - We'll finish whatever the module we are trying to gracefully stop and then escalate to immediate stop
//
// 5 . Sending signals to a stopped runtime could trigger a panic (this is treated as undefined behavior)
// IMPORTANT: "sender" for pal.Signal must close it after we write to ExitStatus channel
type lifeCycle struct {
	mu               sync.Mutex
	state            State
	startFns         []*exec.InvokableHandle
	gracefulStopFns  []*exec.InvokableHandle
	immediateStopFns []*exec.InvokableHandle

	exitCode     uint8
	exitCodeChan chan<- uint8
	listening    bool
}

type action func(rt *Runtime)

// transitionTable[from][to] holds the action invoked when state moves from
// `from` to `to`. A nil entry means the edge is illegal and any attempt to
// take it panics. Actions run with rt.mu released.
var transitionTable [numStates][numStates]action

func init() {
	transitionTable[StateUninitialized][StateInitializing] = initializingAction

	transitionTable[StateInitializing][StateInitializing] = initializingAction
	transitionTable[StateInitializing][StateGracefulStopping] = abortAction
	transitionTable[StateInitializing][StateImmediateStopping] = abortAction
	transitionTable[StateInitializing][StateListening] = listenAction
	transitionTable[StateInitializing][StateStopped] = stoppedAction

	transitionTable[StateListening][StateGracefulStopping] = gracefulStopAction
	transitionTable[StateListening][StateImmediateStopping] = immediateStopAction

	transitionTable[StateGracefulStopping][StateGracefulStopping] = immediateStopAction
	transitionTable[StateGracefulStopping][StateImmediateStopping] = immediateStopAction
	transitionTable[StateGracefulStopping][StateStopped] = stoppedAction

	transitionTable[StateImmediateStopping][StateStopped] = stoppedAction

	transitionTable[StateStopped][StateStopped] = stoppedAction
}

// transition is the only mutator of state. It validates the edge under
// rt.mu, swaps the state, releases the mutex, and then invokes the action
// with the mutex released so long-running stop sequences cannot block
// concurrent transitions (e.g. graceful → immediate escalation).
func (rt *Runtime) transition(target State) {
	rt.mu.Lock()
	from := rt.state
	action := transitionTable[from][target]
	rt.mu.Unlock()
	if action == nil {
		panic(fmt.Sprintf("invalid lifecycle transition from %s -> %s", from, target))
	}
	action(rt)
}

func initializingAction(rt *Runtime) {
	rt.mu.Lock()
	rt.state = StateInitializing
	rt.mu.Unlock()
	rt.setupSignalListeners()
}

func abortAction(_ *Runtime) {
	panic("ABORT: aborting module initialization due to stop signal")
}

func (rt *Runtime) stopAfterInitFailure() {
	rt.exitCode = 1
	rt.transition(StateStopped)
}

// listenAction runs $start for every registered module on the caller's
// goroutine. On the first $start failure it cascades into a graceful
// stop. Spawns the signal-watcher goroutine exactly once before any
// $start runs so a signal that arrives mid-startup is not lost.
func listenAction(rt *Runtime) {
	rt.mu.Lock()
	rt.state = StateListening
	onError := func(message string) {
		writeStderr(rt.env, message)
		rt.exitCode = 1
		rt.mu.Unlock()
		rt.transition(StateGracefulStopping)
	}
	for _, fn := range rt.startFns {
		cx := exec.CreateContext(rt.env)
		res, err := exec.Invoke(cx, fn, nil)
		if err != nil {
			onError(err.Error())
			return
		}
		if errVal, isErr := res.(*values.Error); isErr {
			onError("error: " + errVal.Message + "\n")
			return
		}
	}
	rt.mu.Unlock()
}

// gracefulStopAction walks gracefulStopFns from the end, popping one
// handle per iteration under rt.mu and then dispatching it with the
// mutex released. Both gracefulStopFns and immediateStopFns are popped
// in lockstep so that if a concurrent ImmediateStop signal flips state
// to StateImmediateStopping mid-loop, the next iteration breaks and the
// immediateStopAction picks up from the same module index.
func gracefulStopAction(rt *Runtime) {
	rt.mu.Lock()
	rt.state = StateGracefulStopping
	rt.mu.Unlock()
	if rt.exitCode == 0 {
		rt.exitCode = 130 // 128 + SIGINT
	}
	onError := func(message string) {
		writeStderr(rt.env, message)
		rt.transition(StateImmediateStopping)
	}
	for len(rt.gracefulStopFns) != 0 {
		rt.mu.Lock()
		if rt.state != StateGracefulStopping {
			rt.mu.Unlock()
			return
		}
		fn := rt.gracefulStopFns[len(rt.gracefulStopFns)-1]
		rt.gracefulStopFns = rt.gracefulStopFns[:len(rt.gracefulStopFns)-1]
		rt.immediateStopFns = rt.immediateStopFns[:len(rt.immediateStopFns)-1]
		rt.mu.Unlock()

		cx := exec.CreateContext(rt.env)
		res, err := exec.Invoke(cx, fn, nil)
		if err != nil {
			onError(err.Error())
			return
		}
		if errVal, isErr := res.(*values.Error); isErr {
			onError("error: " + errVal.Message + "\n")
			return
		}
	}
	rt.transition(StateStopped)
}

// immediateStopAction is the graceful action's sibling: post an immediate
// command and let the lifecycle goroutine pick it up. If the goroutine
// is already mid-graceful, escalateAction will have flipped state to
// ImmediateStopping; the goroutine notices on its next module-boundary
// poll.
func immediateStopAction(rt *Runtime) {
	if rt.exitCode == 0 {
		rt.exitCode = 131 // 128  + SIGQUIT
	}
	onError := func(reason string) {
		rt.mu.Unlock()
		writeStderr(rt.env, fmt.Sprintf("panic: immediate stop failed due to %s\n", reason))
		rt.transition(StateStopped)
	}
	rt.mu.Lock()
	rt.state = StateImmediateStopping
	for i := len(rt.immediateStopFns) - 1; i >= 0; i-- {
		fn := rt.immediateStopFns[i]
		cx := exec.CreateContext(rt.env)
		res, err := exec.Invoke(cx, fn, nil)
		if err != nil {
			onError(err.Error())
			return
		}
		if errVal, isErr := res.(*values.Error); isErr {
			onError(errVal.Message)
			return
		}
	}
	rt.mu.Unlock()
	rt.transition(StateStopped)
}

func stoppedAction(rt *Runtime) {
	if rt.state == StateStopped {
		return
	}
	rt.mu.Lock()
	rt.state = StateStopped
	rt.mu.Unlock()
	rt.exitCodeChan <- rt.exitCode
	close(rt.exitCodeChan)
}

// this assume no concurrent calls (fine given triggered from init)
func (rt *Runtime) setupSignalListeners() {
	if rt.listening {
		return
	}
	rt.listening = true
	ch := rt.env.Platform.Signals.Signals
	if ch == nil {
		panic("no signal channel in PAL")
	}
	go func(ch <-chan pal.Signal) {
		for rt.state != StateStopped {
			sig, ok := <-ch
			if !ok {
				return
			}
			switch sig {
			case pal.GracefulStop:
				go rt.transition(StateGracefulStopping)
			case pal.ImmediateStop:
				go rt.transition(StateImmediateStopping)
			default:
				panic("unknown signal")
			}
		}
	}(ch)
}

func writeStderr(env *extern.Env, s string) {
	if env.Platform.IO.Stderr == nil {
		panic("no stderr in PAL")
	}
	_, _ = env.Platform.IO.Stderr([]byte(s))
}
