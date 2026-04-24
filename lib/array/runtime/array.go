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

package array

import (
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/values"
	"fmt"
	"math"
	"sort"
)

const (
	orgName    = "ballerina"
	moduleName = "lang.array"
)

func initArrayModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "push", func(args []values.BalValue) (values.BalValue, error) {
		if list, ok := args[0].(*values.List); ok {
			list.Append(args[1:]...)
			return nil, nil
		}
		return nil, fmt.Errorf("first argument must be an array")
	})
	runtime.RegisterExternFunction(rt, orgName, moduleName, "length", func(args []values.BalValue) (values.BalValue, error) {
		if list, ok := args[0].(*values.List); ok {
			return int64(list.Len()), nil
		}
		return nil, fmt.Errorf("first argument must be an array")
	})
	runtime.RegisterExternFunction(rt, orgName, moduleName, "querySort", func(args []values.BalValue) (values.BalValue, error) {
		keys, ok := args[0].(*values.List)
		if !ok {
			return nil, fmt.Errorf("first argument must be a list")
		}
		directions, ok := args[1].(*values.List)
		if !ok {
			return nil, fmt.Errorf("second argument must be a list")
		}
		indices, ok := args[2].(*values.List)
		if !ok {
			return nil, fmt.Errorf("third argument must be a list")
		}
		payloads, ok := args[3].(*values.List)
		if !ok {
			return nil, fmt.Errorf("fourth argument must be a list")
		}

		rowCount := indices.Len()
		if keys.Len() != rowCount {
			return nil, fmt.Errorf("keys and indices length mismatch: keys=%d indices=%d", keys.Len(), rowCount)
		}
		keyCount := directions.Len()

		directionFlags := make([]bool, keyCount)
		for keyIndex := 0; keyIndex < keyCount; keyIndex++ {
			isAscending, ok := directions.Get(keyIndex).(bool)
			if !ok {
				return nil, fmt.Errorf("directions entry %d must be a bool", keyIndex)
			}
			directionFlags[keyIndex] = isAscending
		}

		keyRows := make([]*values.List, rowCount)
		for row := 0; row < rowCount; row++ {
			rowKeys, ok := keys.Get(row).(*values.List)
			if !ok {
				return nil, fmt.Errorf("keys entry %d must be a list", row)
			}
			if rowKeys.Len() != keyCount {
				return nil, fmt.Errorf("keys entry %d length mismatch: got %d, expected %d", row, rowKeys.Len(), keyCount)
			}
			keyRows[row] = rowKeys
		}

		payloadCount := payloads.Len()
		payloadLists := make([]*values.List, payloadCount)
		for i := 0; i < payloadCount; i++ {
			payloadList, ok := payloads.Get(i).(*values.List)
			if !ok {
				return nil, fmt.Errorf("payload entry %d must be a list", i)
			}
			if payloadList.Len() != rowCount {
				return nil, fmt.Errorf("payload entry %d length mismatch: got %d, expected %d", i, payloadList.Len(), rowCount)
			}
			payloadLists[i] = payloadList
		}

		order := make([]int, rowCount)
		for i := 0; i < rowCount; i++ {
			order[i] = i
		}

		sort.SliceStable(order, func(i, j int) bool {
			leftRow := order[i]
			rightRow := order[j]
			leftKeys := keyRows[leftRow]
			rightKeys := keyRows[rightRow]
			for keyIndex := 0; keyIndex < keyCount; keyIndex++ {
				isAscending := directionFlags[keyIndex]
				cmp := compareQuerySortValues(leftKeys.Get(keyIndex), rightKeys.Get(keyIndex), isAscending)
				switch {
				case cmp < 0:
					return true
				case cmp > 0:
					return false
				}
			}
			return false
		})

		reorderListInPlace(indices, order)
		reorderListInPlace(keys, order)
		for i := 0; i < payloadCount; i++ {
			reorderListInPlace(payloadLists[i], order)
		}
		return nil, nil
	})
}

func reorderListInPlace(list *values.List, order []int) {
	old := make([]values.BalValue, list.Len())
	for i := range list.Len() {
		old[i] = list.Get(i)
	}
	for i, sourceIndex := range order {
		list.FillingSet(i, old[sourceIndex])
	}
}

func compareQuerySortValues(left values.BalValue, right values.BalValue, isAscending bool) int {
	if left == nil {
		if right == nil {
			return 0
		}
		return 1
	}
	if right == nil {
		return -1
	}

	if leftList, ok := left.(*values.List); ok {
		if rightList, ok := right.(*values.List); ok {
			return compareQuerySortLists(leftList, rightList, isAscending)
		}
	}

	if leftFloat, ok := left.(float64); ok {
		if rightFloat, ok := right.(float64); ok {
			leftIsNaN := math.IsNaN(leftFloat)
			rightIsNaN := math.IsNaN(rightFloat)
			switch {
			case leftIsNaN && rightIsNaN:
				return 0
			case leftIsNaN:
				return 1
			case rightIsNaN:
				return -1
			}
		}
	}

	switch values.Compare(left, right) {
	case values.CmpLT:
		if isAscending {
			return -1
		}
		return 1
	case values.CmpGT:
		if isAscending {
			return 1
		}
		return -1
	default:
		return 0
	}
}

func compareQuerySortLists(left *values.List, right *values.List, isAscending bool) int {
	minLength := left.Len()
	if right.Len() < minLength {
		minLength = right.Len()
	}
	for i := range minLength {
		cmp := compareQuerySortValues(left.Get(i), right.Get(i), isAscending)
		if cmp != 0 {
			return cmp
		}
	}
	switch {
	case left.Len() == right.Len():
		return 0
	case isAscending:
		if left.Len() < right.Len() {
			return -1
		}
		return 1
	default:
		if left.Len() > right.Len() {
			return -1
		}
		return 1
	}
}

func init() {
	runtime.RegisterModuleInitializer(initArrayModule)
}
