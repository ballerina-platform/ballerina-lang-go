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

import "testing"

// TestTransitionStoppedToStoppingIsNoOp exercises rt.transition directly
// (rather than through a racy stray-signal-over-channel path, as
// TestLifecycleStartFailureThenStraySignalIsNoOp in lifecycle_test.go does)
// to deterministically prove a stop signal arriving after Stopped is a no-op
// rather than a panic.
func TestTransitionStoppedToStoppingIsNoOp(t *testing.T) {
	for _, target := range []State{StateGracefulStopping, StateImmediateStopping} {
		rt := &Runtime{lifeCycle: lifeCycle{state: StateStopped}}
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("transition(%s) from Stopped panicked: %v", target, r)
				}
			}()
			rt.transition(target)
		}()
		if rt.state != StateStopped {
			t.Fatalf("transition(%s) from Stopped changed state to %s", target, rt.state)
		}
	}
}
