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

package langinternalruntime

import (
	"ballerina-lang-go/decimal"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/values"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

const (
	orgName    = "ballerina"
	moduleName = "lang.__internal"
)

func initInternalModule(rt *runtime.Runtime) {
	runtime.RegisterExternFunction(rt, orgName, moduleName, "querySort", func(args []values.BalValue) (values.BalValue, error) {
		sortKeyRows, ok := args[0].(*values.List)
		if !ok {
			return nil, fmt.Errorf("first argument must be a list")
		}
		sortDirections, ok := args[1].(*values.List)
		if !ok {
			return nil, fmt.Errorf("second argument must be a list")
		}
		rowIndices, ok := args[2].(*values.List)
		if !ok {
			return nil, fmt.Errorf("third argument must be a list")
		}
		payloadRows, ok := args[3].(*values.List)
		if !ok {
			return nil, fmt.Errorf("fourth argument must be a list")
		}

		rowCount := rowIndices.Len()
		if sortKeyRows.Len() != rowCount {
			return nil, fmt.Errorf("sort keys and row indices length mismatch: keys=%d indices=%d", sortKeyRows.Len(), rowCount)
		}
		keyCount := sortDirections.Len()

		directionFlags := make([]bool, keyCount)
		for keyIndex := 0; keyIndex < keyCount; keyIndex++ {
			isAscending, ok := sortDirections.Get(keyIndex).(bool)
			if !ok {
				return nil, fmt.Errorf("sort direction %d must be a bool", keyIndex)
			}
			directionFlags[keyIndex] = isAscending
		}

		keyRows := make([]*values.List, rowCount)
		for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
			rowKeys, ok := sortKeyRows.Get(rowIndex).(*values.List)
			if !ok {
				return nil, fmt.Errorf("sort key row %d must be a list", rowIndex)
			}
			if rowKeys.Len() != keyCount {
				return nil, fmt.Errorf("sort key row %d length mismatch: got %d, expected %d", rowIndex, rowKeys.Len(), keyCount)
			}
			keyRows[rowIndex] = rowKeys
		}

		payloadCount := payloadRows.Len()
		payloadLists := make([]*values.List, payloadCount)
		for payloadIndex := 0; payloadIndex < payloadCount; payloadIndex++ {
			payloadList, ok := payloadRows.Get(payloadIndex).(*values.List)
			if !ok {
				return nil, fmt.Errorf("payload row %d must be a list", payloadIndex)
			}
			if payloadList.Len() != rowCount {
				return nil, fmt.Errorf("payload row %d length mismatch: got %d, expected %d", payloadIndex, payloadList.Len(), rowCount)
			}
			payloadLists[payloadIndex] = payloadList
		}

		order := make([]int, rowCount)
		for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
			order[rowIndex] = rowIndex
		}

		sort.SliceStable(order, func(i, j int) bool {
			leftRow := order[i]
			rightRow := order[j]
			leftKeys := keyRows[leftRow]
			rightKeys := keyRows[rightRow]
			for keyIndex := 0; keyIndex < keyCount; keyIndex++ {
				cmp := compareQuerySortValues(leftKeys.Get(keyIndex), rightKeys.Get(keyIndex), directionFlags[keyIndex])
				switch cmp {
				case values.CmpLT:
					return true
				case values.CmpGT:
					return false
				}
			}
			return false
		})

		reorderListInPlace(rowIndices, order)
		reorderListInPlace(sortKeyRows, order)
		for payloadIndex := 0; payloadIndex < payloadCount; payloadIndex++ {
			reorderListInPlace(payloadLists[payloadIndex], order)
		}
		return nil, nil
	})
	runtime.RegisterExternFunction(rt, orgName, moduleName, "queryGroup", func(args []values.BalValue) (values.BalValue, error) {
		if len(args) != 3 {
			return nil, fmt.Errorf("queryGroup expects 3 arguments, got %d", len(args))
		}
		rows, ok := args[0].(*values.List)
		if !ok {
			return nil, fmt.Errorf("first argument must be a list")
		}
		keyRows, ok := args[1].(*values.List)
		if !ok {
			return nil, fmt.Errorf("second argument must be a list")
		}
		scalarFlags, ok := args[2].(*values.List)
		if !ok {
			return nil, fmt.Errorf("third argument must be a list")
		}
		return queryGroup(rows, keyRows, scalarFlags)
	})
	runtime.RegisterExternFunction(rt, orgName, moduleName, "queryCollect", func(args []values.BalValue) (values.BalValue, error) {
		if len(args) != 3 {
			return nil, fmt.Errorf("queryCollect expects 3 arguments, got %d", len(args))
		}
		rows, ok := args[0].(*values.List)
		if !ok {
			return nil, fmt.Errorf("first argument must be a list")
		}
		slotCount, ok := args[1].(int64)
		if !ok {
			return nil, fmt.Errorf("second argument must be an int")
		}
		flattenFlags, ok := args[2].(*values.List)
		if !ok {
			return nil, fmt.Errorf("third argument must be a list")
		}
		return queryCollect(rows, int(slotCount), flattenFlags)
	})
}

type queryGroupState struct {
	row *values.List
}

type queryGroupIndex map[string]int

func queryGroup(rows *values.List, keyRows *values.List, scalarFlags *values.List) (*values.List, error) {
	rowCount := rows.Len()
	if keyRows.Len() != rowCount {
		return nil, fmt.Errorf("group keys and rows length mismatch: keys=%d rows=%d", keyRows.Len(), rowCount)
	}
	slotCount := scalarFlags.Len()
	scalarSlots := make([]bool, slotCount)
	for slot := 0; slot < slotCount; slot++ {
		isScalar, ok := scalarFlags.Get(slot).(bool)
		if !ok {
			return nil, fmt.Errorf("group scalar flag %d must be a bool", slot)
		}
		scalarSlots[slot] = isScalar
	}

	result := values.NewList(0, semtypes.LIST, nil)
	groups := make([]queryGroupState, 0)
	groupIndices := make(queryGroupIndex)
	expectedKeyArity := -1
	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		sourceRow, ok := rows.Get(rowIndex).(*values.List)
		if !ok {
			return nil, fmt.Errorf("query row %d must be a list", rowIndex)
		}
		if sourceRow.Len() != slotCount {
			return nil, fmt.Errorf("query row %d length mismatch: got %d, expected %d", rowIndex, sourceRow.Len(), slotCount)
		}
		keyRow, ok := keyRows.Get(rowIndex).(*values.List)
		if !ok {
			return nil, fmt.Errorf("group key row %d must be a list", rowIndex)
		}
		if expectedKeyArity == -1 {
			expectedKeyArity = keyRow.Len()
		} else if keyRow.Len() != expectedKeyArity {
			return nil, fmt.Errorf("group key row %d length mismatch: got %d, expected %d", rowIndex, keyRow.Len(), expectedKeyArity)
		}
		groupIndex, keySignature, found := findQueryGroup(groupIndices, keyRow)
		if !found {
			groupRow := createQueryGroupRow(sourceRow, scalarSlots)
			groups = append(groups, queryGroupState{row: groupRow})
			groupIndices[keySignature] = len(groups) - 1
			result.Append(groupRow)
			continue
		}
		appendQueryGroupRow(groups[groupIndex].row, sourceRow, scalarSlots)
	}
	return result, nil
}

func findQueryGroup(groupIndices queryGroupIndex, keyRow *values.List) (int, string, bool) {
	keySignature := queryGroupKeySignature(keyRow)
	groupIndex, ok := groupIndices[keySignature]
	if !ok {
		return 0, keySignature, false
	}
	return groupIndex, keySignature, true
}

func queryGroupKeySignature(keyRow *values.List) string {
	var b strings.Builder
	writeQueryGroupValueSignature(&b, keyRow, make(map[uintptr]bool))
	return b.String()
}

func writeQueryGroupValueSignature(b *strings.Builder, v values.BalValue, visited map[uintptr]bool) {
	switch typedValue := v.(type) {
	case nil:
		b.WriteString("nil;")
	case bool:
		b.WriteString("bool:")
		b.WriteString(strconv.FormatBool(typedValue))
		b.WriteByte(';')
	case int64:
		b.WriteString("int:")
		b.WriteString(strconv.FormatInt(typedValue, 10))
		b.WriteByte(';')
	case float64:
		writeQueryGroupFloatSignature(b, typedValue)
	case string:
		writeQueryGroupStringSignature(b, "string", typedValue)
	case *decimal.Decimal:
		b.WriteString("decimal:")
		writeQueryGroupFloatSignature(b, typedValue.Float64())
	case *values.List:
		writeQueryGroupListSignature(b, typedValue, visited)
	case *values.Map:
		writeQueryGroupMapSignature(b, typedValue, visited)
	default:
		b.WriteString("unknown:")
		b.WriteString(reflect.TypeOf(typedValue).String())
		b.WriteByte(';')
	}
}

func writeQueryGroupFloatSignature(b *strings.Builder, f float64) {
	b.WriteString("float:")
	switch {
	case math.IsNaN(f):
		b.WriteString("NaN")
	case f == 0:
		b.WriteByte('0')
	default:
		b.WriteString(strconv.FormatFloat(f, 'g', -1, 64))
	}
	b.WriteByte(';')
}

func writeQueryGroupStringSignature(b *strings.Builder, tag string, s string) {
	b.WriteString(tag)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(len(s)))
	b.WriteByte(':')
	b.WriteString(s)
	b.WriteByte(';')
}

func writeQueryGroupListSignature(b *strings.Builder, list *values.List, visited map[uintptr]bool) {
	ptr := uintptr(unsafe.Pointer(list))
	if visited[ptr] {
		b.WriteString("list:cycle;")
		return
	}
	visited[ptr] = true
	defer delete(visited, ptr)

	b.WriteString("list:")
	b.WriteString(strconv.Itoa(list.Len()))
	b.WriteByte('[')
	for i := 0; i < list.Len(); i++ {
		writeQueryGroupValueSignature(b, list.Get(i), visited)
	}
	b.WriteByte(']')
}

func writeQueryGroupMapSignature(b *strings.Builder, mapping *values.Map, visited map[uintptr]bool) {
	ptr := uintptr(unsafe.Pointer(mapping))
	if visited[ptr] {
		b.WriteString("map:cycle;")
		return
	}
	visited[ptr] = true
	defer delete(visited, ptr)

	keys := mapping.Keys()
	sort.Strings(keys)
	b.WriteString("map:")
	b.WriteString(strconv.Itoa(len(keys)))
	b.WriteByte('{')
	for _, key := range keys {
		writeQueryGroupStringSignature(b, "key", key)
		value, _ := mapping.Get(key)
		writeQueryGroupValueSignature(b, value, visited)
	}
	b.WriteByte('}')
}

func createQueryGroupRow(sourceRow *values.List, scalarSlots []bool) *values.List {
	groupRow := values.NewList(0, semtypes.LIST, nil)
	for slot, isScalar := range scalarSlots {
		value := sourceRow.Get(slot)
		if isScalar {
			groupRow.Append(value)
			continue
		}
		valuesForGroup := values.NewList(0, semtypes.LIST, nil)
		valuesForGroup.Append(value)
		groupRow.Append(valuesForGroup)
	}
	return groupRow
}

func appendQueryGroupRow(groupRow *values.List, sourceRow *values.List, scalarSlots []bool) {
	for slot, isScalar := range scalarSlots {
		if isScalar {
			continue
		}
		valuesForGroup := groupRow.Get(slot).(*values.List)
		valuesForGroup.Append(sourceRow.Get(slot))
	}
}

func queryCollect(rows *values.List, slotCount int, flattenFlags *values.List) (*values.List, error) {
	if slotCount < 0 {
		return nil, fmt.Errorf("slot count cannot be negative")
	}
	if flattenFlags.Len() != slotCount {
		return nil, fmt.Errorf("collect flatten flags length mismatch: flags=%d slots=%d", flattenFlags.Len(), slotCount)
	}
	flattenSlots := make([]bool, slotCount)
	for slot := 0; slot < slotCount; slot++ {
		flatten, ok := flattenFlags.Get(slot).(bool)
		if !ok {
			return nil, fmt.Errorf("collect flatten flag %d must be a bool", slot)
		}
		flattenSlots[slot] = flatten
	}
	resultRow := values.NewList(0, semtypes.LIST, nil)
	for slot := 0; slot < slotCount; slot++ {
		resultRow.Append(values.NewList(0, semtypes.LIST, nil))
	}
	for rowIndex := 0; rowIndex < rows.Len(); rowIndex++ {
		row, ok := rows.Get(rowIndex).(*values.List)
		if !ok {
			return nil, fmt.Errorf("query row %d must be a list", rowIndex)
		}
		if row.Len() != slotCount {
			return nil, fmt.Errorf("query row %d length mismatch: got %d, expected %d", rowIndex, row.Len(), slotCount)
		}
		for slot := 0; slot < slotCount; slot++ {
			valuesForSlot := resultRow.Get(slot).(*values.List)
			if !flattenSlots[slot] {
				valuesForSlot.Append(row.Get(slot))
				continue
			}
			nestedValues, ok := row.Get(slot).(*values.List)
			if !ok {
				return nil, fmt.Errorf("query row %d slot %d must be a list for flattening", rowIndex, slot)
			}
			for valueIndex := 0; valueIndex < nestedValues.Len(); valueIndex++ {
				valuesForSlot.Append(nestedValues.Get(valueIndex))
			}
		}
	}
	return resultRow, nil
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

func compareQuerySortValues(left values.BalValue, right values.BalValue, isAscending bool) values.CompareResult {
	return values.CompareK(left, right, isAscending)
}

func init() {
	runtime.RegisterModuleInitializer(initInternalModule)
}
