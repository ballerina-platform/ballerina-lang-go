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

package native

import (
	"strings"

	"ballerina-lang-go/runtime"
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/values"
)

const (
	orgName    = "ballerina"
	moduleName = "log"
)

// logLevelWeight maps level name → numeric weight for filtering.
// Higher weight = higher severity (matches jBallerina's logLevelWeight map).
var logLevelWeight = map[string]int{
	"ERROR": 1000,
	"WARN":  900,
	"INFO":  800,
	"DEBUG": 700,
}

// defaultLevelWeight corresponds to INFO — the jBallerina default.
const defaultLevelWeight = 800

func isLevelEnabled(level string) bool {
	w, ok := logLevelWeight[level]
	if !ok {
		return true
	}
	return w >= defaultLevelWeight
}

// logfmtEscaper replaces special characters for LOGFMT string values.
// Matches the escape() function in jBallerina's natives.bal.
var logfmtEscaper = strings.NewReplacer(
	`\`, `\\`,
	"\t", `\t`,
	"\n", `\n`,
	"\r", `\r`,
	`'`, `\'`,
	`"`, `\"`,
)

func formatLogfmt(level, timeStr, msg string, errVal *values.Error) string {
	var b strings.Builder
	b.WriteString("time=")
	b.WriteString(timeStr)
	b.WriteString(" level=")
	b.WriteString(level)
	b.WriteString(` module=""`)
	b.WriteString(` message="`)
	b.WriteString(logfmtEscaper.Replace(msg))
	b.WriteByte('"')

	if errVal != nil {
		b.WriteString(" error=")
		b.WriteString(errVal.String(nil))
	}

	return b.String()
}

func initLogModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "externPrintLog",
		func(_ *extern.Context, args []values.BalValue) (values.BalValue, error) {
			level, _ := args[0].(string)
			msg, _ := args[1].(string)

			if !isLevelEnabled(level) {
				return nil, nil
			}

			var errVal *values.Error
			if args[2] != nil {
				errVal, _ = args[2].(*values.Error)
			}

			now := rt.Platform().Time.Now()
			timeStr := now.Format("2006-01-02T15:04:05.000Z07:00")

			output := formatLogfmt(level, timeStr, msg, errVal)
			_, _ = rt.Platform().IO.Stderr([]byte(output + "\n"))

			return nil, nil
		})
}

func init() {
	runtime.RegisterModuleInitializer(initLogModule)
}
