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

package corpus

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"ballerina-lang-go/ast"
	"ballerina-lang-go/bir"
	bircodec "ballerina-lang-go/bir/codec"
	"ballerina-lang-go/context"
	"ballerina-lang-go/desugar"
	"ballerina-lang-go/model"
	"ballerina-lang-go/model/symbolpool"
	"ballerina-lang-go/parser"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/runtime"
	"ballerina-lang-go/semantics"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util"
	"ballerina-lang-go/tools/diagnostics"
	"ballerina-lang-go/tools/text"

	_ "ballerina-lang-go/lib/rt"
)

const (
	corpusProjectBaseDir            = "../corpus/project"
	corpusProjectIntegrationBaseDir = "../corpus/integration/project"

	panicPrefix = "panic: "
)

var (
	update = flag.Bool("update", false, "update corpus integration test outputs")

	skipIntegrationTests = []string{
		"subset8/08-mutate/record1-v.bal",
		"subset8/08-colon/field1-v.bal",
		"subset8/08-mapping/8-v.bal",
		"subset8/08-ifelse/ifelse1-e.bal",
		"subset8/08-ifelse/ifelse2-e.bal",
		"subset8/08-ifelse/ifelse3-e.bal",
		"subset8/08-mutate/mappingassign1-e.bal",
		"subset8/08-mutate/mappingassign2-e.bal",
		"subset8/08-narrowing/3-e.bal",
		"subset8/08-rest/syntax1-e.bal",
		"subset8/08-rest/syntax2-e.bal",
		"subset8/08-typetest/not2-e.bal",
		// Tests that cause unrecoverable Go runtime errors.
		// https://github.com/ballerina-platform/ballerina-lang-go/issues/364
		"subset8/08-comparable/order5-v.bal",
		// Missing init exp for constants.
		"subset8/08-const/2-v.bal",
		"subset8/08-const/3-v.bal",
		"subset8/08-const/const3-v.bal",
		// Migrated from nballerina testSuite/12-nested/order4-e.bal: stack overflow in semantics.ResolveLocalNodes.
		"subset8/08-nested/order4-e.bal",

		// Expected error: migrated -e tests for which the frontend currently produces no
		// diagnostic. Skipped so we don't bake an empty stderr into the expected fixture.
		"subset8/08-bitwise/complement3-e.bal",
		"subset8/08-bug/assignforeach-e.bal",
		"subset8/08-bug/init2-e.bal",
		"subset8/08-bug/intersect1-e.bal",
		"subset8/08-bug/intersect2-e.bal",
		"subset8/08-bug/main2-e.bal",
		"subset8/08-bug/matchwild1-e.bal",
		"subset8/08-colon/ident1-e.bal",
		"subset8/08-colon/ident2-e.bal",
		"subset8/08-colon/ident3-e.bal",
		"subset8/08-const/1-e.bal",
		"subset8/08-const/10-e.bal",
		"subset8/08-const/11-e.bal",
		"subset8/08-const/12-e.bal",
		"subset8/08-const/13-e.bal",
		"subset8/08-const/14-e.bal",
		"subset8/08-const/15-e.bal",
		"subset8/08-const/16-e.bal",
		"subset8/08-const/17-e.bal",
		"subset8/08-const/18-e.bal",
		"subset8/08-const/23-e.bal",
		"subset8/08-const/7-e.bal",
		"subset8/08-const/8-e.bal",
		"subset8/08-const/9-e.bal",
		"subset8/08-decimal/const5-e.bal",
		"subset8/08-decimal/const6-e.bal",
		"subset8/08-equal/3-e.bal",
		"subset8/08-equal/4-e.bal",
		"subset8/08-equal/5-e.bal",
		"subset8/08-fill/10-e.bal",
		"subset8/08-fill/15-e.bal",
		"subset8/08-fill/18-e.bal",
		"subset8/08-fill/21-e.bal",
		"subset8/08-fill/22-e.bal",
		"subset8/08-float/5-e.bal",
		"subset8/08-float/7-e.bal",
		"subset8/08-hex/decimal1-e.bal",
		"subset8/08-inclusive/compoundassign3-e.bal",
		"subset8/08-inclusive/construct5-e.bal",
		"subset8/08-inclusive/duplicate2-e.bal",
		"subset8/08-infinite/infiniteRecord4-e.bal",
		"subset8/08-list/17-e.bal",
		"subset8/08-list/6-e.bal",
		"subset8/08-list/fixedlength1-e.bal",
		"subset8/08-list/fixedlength2-e.bal",
		"subset8/08-list/fixedlength3-e.bal",
		"subset8/08-list/fixedlength4-e.bal",
		"subset8/08-list/fixedlength5-e.bal",
		"subset8/08-list/fixedlength6-e.bal",
		"subset8/08-list/fixedlength7-e.bal",
		"subset8/08-list/fixedlength8-e.bal",
		"subset8/08-list/fixedlength9-e.bal",
		"subset8/08-map/compoundassign-e.bal",
		"subset8/08-mapping/14-e.bal",
		"subset8/08-mapping/4-e.bal",
		"subset8/08-mapping/6-e.bal",
		"subset8/08-mapping/7-e.bal",
		"subset8/08-match/19-e.bal",
		"subset8/08-match/3-e.bal",
		"subset8/08-match/7-e.bal",
		"subset8/08-narrowing/10-e.bal",
		"subset8/08-narrowing/11-e.bal",
		"subset8/08-narrowing/12-e.bal",
		"subset8/08-narrowing/15-e.bal",
		"subset8/08-narrowing/2-e.bal",
		"subset8/08-narrowing/4-e.bal",
		"subset8/08-narrowing/5-e.bal",
		"subset8/08-narrowing/6-e.bal",
		"subset8/08-narrowing/8-e.bal",
		"subset8/08-narrowing/if18-e.bal",
		"subset8/08-narrowing/unreach3-e.bal",
		"subset8/08-narrowing/unreach4-e.bal",
		"subset8/08-nillifting/compound1-e.bal",
		"subset8/08-nillifting/compound11-e.bal",
		"subset8/08-nillifting/compound2-e.bal",
		"subset8/08-nillifting/compound3-e.bal",
		"subset8/08-nillifting/compound5-e.bal",
		"subset8/08-nillifting/compound7-e.bal",
		"subset8/08-record/assign1-e.bal",
		"subset8/08-record/compoundassign4-e.bal",
		"subset8/08-semtype/xml-e.bal",
		"subset8/08-singleton/decimal10-e.bal",
		"subset8/08-singleton/decimal11-e.bal",
		"subset8/08-singleton/decimal12-e.bal",
		"subset8/08-singleton/decimal13-e.bal",
		"subset8/08-singleton/decimal2-e.bal",
		"subset8/08-singleton/decimal4-e.bal",
		"subset8/08-singleton/decimal5-e.bal",
		"subset8/08-singleton/decimal6-e.bal",
		"subset8/08-singleton/decimal7-e.bal",
		"subset8/08-singleton/decimal8-e.bal",
		"subset8/08-singleton/decimal9-e.bal",
		"subset8/08-singleton/nil1-e.bal",
		"subset8/08-singleton/not1-e.bal",
		"subset8/08-singleton/string1-e.bal",
		"subset8/08-singleton/stringconcat1-e.bal",
		"subset8/08-string/1-e.bal",
		"subset8/08-string/5-e.bal",
		"subset8/08-typecast/8-e.bal",
		"subset8/08-unused/unused1-e.bal",
		"subset8/08-unused/unused2-e.bal",
		"subset8/08-unused/unused3-e.bal",
		"subset8/08-unused/unused4-e.bal",
		"subset8/08-unused/unused5-e.bal",
		"subset8/08-unused/unused6-e.bal",

		// Expected clean run: migrated -v tests that produce diagnostics or runtime errors
		"subset8/08-bench/ackermann-v.bal",
		"subset8/08-bench/map-v.bal",
		"subset8/08-bitwise/shift1-v.bal",
		"subset8/08-bitwise/shift2-v.bal",
		"subset8/08-bitwise/shift3-v.bal",
		"subset8/08-bug/charcast1-v.bal",
		"subset8/08-bug/charcast2-v.bal",
		"subset8/08-bug/fill1-v.bal",
		"subset8/08-bug/shiftresulttype1-v.bal",
		"subset8/08-bug/shiftresulttype2-v.bal",
		"subset8/08-const/10-v.bal",
		"subset8/08-const/7-v.bal",
		"subset8/08-const/8-v.bal",
		"subset8/08-decimal/add1-v.bal",
		"subset8/08-decimal/add7-v.bal",
		"subset8/08-decimal/const1-v.bal",
		"subset8/08-decimal/const7-v.bal",
		"subset8/08-decimal/div1-v.bal",
		"subset8/08-decimal/div5-v.bal",
		// decimal equality
		"subset8/08-decimal/eq1-v.bal",
		"subset8/08-decimal/eq2-v.bal",
		"subset8/08-singleton/decimal1-v.bal",
		"subset8/08-singleton/decimal3-v.bal",
		"subset8/08-decimal/exacteq1-v.bal",
		"subset8/08-decimal/exacteq2-v.bal",
		// invalid number conversion
		"subset8/08-decimal/fromint1-v.bal",
		// decimal rounding error
		"subset8/08-decimal/map1-v.bal",
		"subset8/08-decimal/vardecl1-v.bal",
		"subset8/08-decimal/mul1-v.bal",
		"subset8/08-decimal/mul6-v.bal",
		"subset8/08-decimal/mul7-v.bal",
		"subset8/08-decimal/neg1-v.bal",
		"subset8/08-decimal/rem1-v.bal",
		"subset8/08-decimal/rem2-v.bal",
		"subset8/08-decimal/rem5-v.bal",
		"subset8/08-decimal/sub1-v.bal",
		"subset8/08-decimal/sub4-v.bal",
		"subset8/08-decimal/tofloat1-v.bal",
		"subset8/08-decimal/tofloat2-v.bal",
		"subset8/08-decimal/tofloat3-v.bal",
		"subset8/08-decimal/toint1-v.bal",
		"subset8/08-decimal/toint7-v.bal",
		"subset8/08-error/10-v.bal",
		"subset8/08-error/check1-v.bal",
		"subset8/08-error/check10-v.bal",
		"subset8/08-error/check3-v.bal",
		// invalid filling value
		"subset8/08-fill/1-v.bal",
		"subset8/08-fill/17-v.bal",
		"subset8/08-fill/2-v.bal",
		"subset8/08-fill/3-v.bal",
		"subset8/08-fill/4-v.bal",
		"subset8/08-fill/5-v.bal",
		"subset8/08-fill/8-v.bal",
		"subset8/08-fill/chain2-v.bal",
		"subset8/08-fill/fill1-v.bal",
		"subset8/08-fill/fill2-v.bal",
		"subset8/08-fill/fill3-v.bal",
		"subset8/08-fill/fill7-v.bal",
		"subset8/08-fill/methodcall1-v.bal",
		"subset8/08-float/10-v.bal",
		"subset8/08-float/12-v.bal",
		"subset8/08-float/14-v.bal",
		"subset8/08-float/16-v.bal",
		"subset8/08-float/18-v.bal",
		"subset8/08-float/19-v.bal",
		"subset8/08-float/2-v.bal",
		"subset8/08-float/20-v.bal",
		"subset8/08-float/21-v.bal",
		"subset8/08-float/22-v.bal",
		"subset8/08-float/23-v.bal",
		"subset8/08-float/24-v.bal",
		"subset8/08-float/9-v.bal",
		"subset8/08-float/const3-v.bal",
		"subset8/08-function/intersection11-v.bal",
		"subset8/08-function/intersection13-v.bal",
		"subset8/08-future/fieldexpr1-v.bal",
		"subset8/08-future/lib1-v.bal",
		"subset8/08-future/main-v.bal",
		"subset8/08-future/main2-v.bal",
		"subset8/08-future/never-v.bal",
		"subset8/08-future/xmlsubtype-v.bal",
		"subset8/08-ifelse/ifelse4-v.bal",
		// invalid list eq
		"subset8/08-list/1-v.bal",
		"subset8/08-list/14-v.bal",
		"subset8/08-list/equal-v.bal",
		"subset8/08-nested/eqcycle2-v.bal",
		// invalid compound assign result
		"subset8/08-list/compoundassign1-v.bal",
		"subset8/08-list/compoundassign2-v.bal",
		// invalid map eq
		"subset8/08-map/equal-v.bal",
		"subset8/08-mapping/1-v.bal",
		"subset8/08-mapping/5-v.bal",
		"subset8/08-nested/bdd1-v.bal",
		"subset8/08-nested/eqcycle1-v.bal",
		"subset8/08-match/18-v.bal",
		"subset8/08-match/2-v.bal",
		"subset8/08-match/4-v.bal",
		"subset8/08-match/float3-v.bal",
		// float +/- zero
		"subset8/08-narrowing/3-v.bal",
		"subset8/08-narrowing/7-v.bal",
		"subset8/08-nested/fill1-v.bal",
		"subset8/08-nested/push1-v.bal",
		"subset8/08-rest/construct7-v.bal",
		"subset8/08-semtype/array-v.bal",
		"subset8/08-semtype/not1-v.bal",
		"subset8/08-semtype/objectCompliment-v.bal",
		"subset8/08-semtype/optional-field-record1-v.bal",
		"subset8/08-semtype/optional-field-record3-v.bal",
		"subset8/08-semtype/proj10-v.bal",
		"subset8/08-semtype/proj2-v.bal",
		"subset8/08-semtype/proj3-v.bal",
		"subset8/08-semtype/proj7-v.bal",
		"subset8/08-semtype/proj8-v.bal",
		"subset8/08-semtype/readonly-record-field-v.bal",
		"subset8/08-semtype/readonly-record-field2-v.bal",
		"subset8/08-semtype/record-proj-v.bal",
		"subset8/08-singleton/float1-v.bal",
		"subset8/08-singleton/floattest1-v.bal",
		"subset8/08-singleton/floattest2-v.bal",
		"subset8/08-singleton/proj4-v.bal",
		"subset8/08-singleton/typecast1-v.bal",
		"subset8/08-string/10-v.bal",
		"subset8/08-string/11-v.bal",
		"subset8/08-string/12-v.bal",
		"subset8/08-string/13-v.bal",
		"subset8/08-string/15-v.bal",
		"subset8/08-string/16-v.bal",
		"subset8/08-string/17-v.bal",
		"subset8/08-tuple/context1-v.bal",
		"subset8/08-tuple/push2-v.bal",
		"subset8/08-tuple/tupleunion1-v.bal",
		// invalid int conversion
		"subset8/08-typecast/10-v.bal",
		"subset8/08-typecast/13-v.bal",
		"subset8/08-typecast/14-v.bal",
		"subset8/08-typecast/11-v.bal",
		"subset8/08-typecast/12-v.bal",
		"subset8/08-typecast/16-v.bal",
		"subset8/08-union/construct4-v.bal",

		// Expected runtime panic, but got nothing/wrong panic
		"subset8/08-bug/fill2-p.bal",
		"subset8/08-bug/fill4-p.bal",
		"subset8/08-bytearr/2-p.bal",
		"subset8/08-bytearr/3-p.bal",
		"subset8/08-bytearr/4-p.bal",
		"subset8/08-exact/array1-p.bal",
		"subset8/08-exact/map1-p.bal",
		"subset8/08-exact/push1-p.bal",
		"subset8/08-exact/record1-p.bal",
		"subset8/08-fill/14-p.bal",
		"subset8/08-fill/23-p.bal",
		"subset8/08-fill/9-p.bal",
		"subset8/08-inclusive/inherent1-p.bal",
		"subset8/08-inttest/typecast1-p.bal",
		"subset8/08-list/push6-p.bal",
		"subset8/08-map/int5-p.bal",
		"subset8/08-nested/exact2-p.bal",
		"subset8/08-nested/exact4-p.bal",
		"subset8/08-nested/exact5-p.bal",
		"subset8/08-nested/exact6-p.bal",
		"subset8/08-record/inherent1-p.bal",
		"subset8/08-record/inherent2-p.bal",
		"subset8/08-tuple/exact1-p.bal",
		"subset8/08-tuple/exact2-p.bal",
		"subset8/08-tuple/push3-p.bal",

		// invalid float overflow
		"subset8/08-decimal/fromfloat5-p.bal",
		"subset8/08-decimal/fromfloat6-p.bal",
		// Expected runtime panic, but got frontend error.
		"subset8/08-decimal/tofloat4-p.bal",
		"subset8/08-decimal/toint13-p.bal",
		"subset8/08-fill/fill4-p.bal",
		"subset8/08-nested/exact1-p.bal",
		"subset8/08-nested/proj1-p.bal",
		"subset8/08-rest/exact1-p.bal",
		"subset8/08-typecast/2-p.bal",
		"subset8/08-typecast/6-p.bal",

		// Expected clean run: migrated -v tests whose expected stdout contains a runtime
		// panic (). A -v test must complete without panicking.
		"subset8/08-bug/listfill1-v.bal",
		"subset8/08-fill/11-v.bal",
		"subset8/08-fill/12-v.bal",
		"subset8/08-fill/19-v.bal",
		"subset8/08-fill/20-v.bal",
		"subset8/08-fill/order-v.bal",
		"subset8/08-list/fixedlength1-v.bal",
		"subset8/08-list/fixedlength2-v.bal",
		"subset8/08-semtype/anydata-v.bal",
		"subset8/08-semtype/fixed-length-array-large-v.bal",
		"subset8/08-semtype/fixed-length-array-readonly-v.bal",
		"subset8/08-semtype/fixed-length-array-tuple-readonly-v.bal",
		"subset8/08-semtype/fixed-length-array-tuple-v.bal",
		"subset8/08-semtype/fixed-length-array-tuple2-v.bal",
		"subset8/08-semtype/fixed-length-array-v.bal",
		"subset8/08-semtype/fixed-length-array2-v.bal",
		"subset8/08-semtype/proj6-v.bal",
		"subset8/08-semtype/proj9-v.bal",
		"subset8/08-semtype/recurse-v.bal",
		"subset8/08-semtype/table-readonly-v.bal",
		"subset8/08-semtype/table-v.bal",
		"subset8/08-semtype/table2-v.bal",
		"subset8/08-semtype/table3-v.bal",
		"subset8/08-semtype/xml-complex-ro-v.bal",
		"subset8/08-semtype/xml-complex-rw-v.bal",
		"subset8/08-semtype/xml-never-v.bal",
		"subset8/08-semtype/xml-readonly-v.bal",
		"subset8/08-semtype/xml-sequence-v.bal",
		"subset8/08-tuple/comp9-v.bal",

		// Expected frontend error: migrated -e tests where pi did not catch the error in
		// the front-end. The expected stderr is either a runtime error () or a
		// compiler internal/unimplemented bailout (). The front-end should
		// detect these statically before reaching this stage.
		"subset8/08-bug/stringop1-e.bal",
		"subset8/08-bug/unusedimport-e.bal",
		"subset8/08-compoundassign/9-e.bal",
		"subset8/08-const/4-e.bal",
		"subset8/08-const/5-e.bal",
		"subset8/08-const/6-e.bal",
		"subset8/08-decimal/add2-e.bal",
		"subset8/08-decimal/add3-e.bal",
		"subset8/08-decimal/add4-e.bal",
		"subset8/08-decimal/add5-e.bal",
		"subset8/08-decimal/add6-e.bal",
		"subset8/08-decimal/div2-e.bal",
		"subset8/08-decimal/div3-e.bal",
		"subset8/08-decimal/div4-e.bal",
		"subset8/08-decimal/fromfloat2-e.bal",
		"subset8/08-decimal/fromfloat3-e.bal",
		"subset8/08-decimal/mul2-e.bal",
		"subset8/08-decimal/mul3-e.bal",
		"subset8/08-decimal/mul4-e.bal",
		"subset8/08-decimal/mul5-e.bal",
		"subset8/08-decimal/rem3-e.bal",
		"subset8/08-decimal/rem4-e.bal",
		"subset8/08-decimal/sub2-e.bal",
		"subset8/08-decimal/sub3-e.bal",
		"subset8/08-decimal/toint2-e.bal",
		"subset8/08-decimal/toint3-e.bal",
		"subset8/08-decimal/toint4-e.bal",
		"subset8/08-decimal/toint5-e.bal",
		"subset8/08-decimal/toint6-e.bal",
		"subset8/08-error/check8-e.bal",
		"subset8/08-float/15-e.bal",
		// rest param not supported in dependently typed functions
		"subset8/08-function/dependent-fn-5-e.bal",
		"subset8/08-future/langlib2-e.bal",
		"subset8/08-future/langlib3-e.bal",
		"subset8/08-inclusive/compoundassign2-e.bal",
		"subset8/08-inclusive/fieldlvalue5-e.bal",
		"subset8/08-list/10-e.bal",
		"subset8/08-list/compoundassign5-e.bal",
		"subset8/08-mapping/9-e.bal",
		"subset8/08-narrowing/7-e.bal",
		"subset8/08-narrowing/9-e.bal",
		"subset8/08-nillifting/compound10-e.bal",
		"subset8/08-nillifting/compound12-e.bal",
		"subset8/08-nillifting/compound4-e.bal",
		"subset8/08-nillifting/compound6-e.bal",
		"subset8/08-nillifting/compound8-e.bal",
		"subset8/08-nillifting/compound9-e.bal",
		"subset8/08-record/fieldlvalue6-e.bal",
		"subset8/08-tuple/construct4-e.bal",

		// Missing error location:
		"subset8/08-const/def-e.bal",

		// Wrong runtime panic: migrated -p tests where pi raises "unsupported type" instead
		// of the intended runtime panic (overflow/conversion/etc). The test exercises a feature
		// pi has not implemented at runtime (decimal arithmetic on *big.Rat, typed list/map
		// element checks), so the right panic is never produced.
		"subset8/08-decimal/add10-p.bal",
		"subset8/08-decimal/add11-p.bal",
		"subset8/08-decimal/add12-p.bal",
		"subset8/08-decimal/add8-p.bal",
		"subset8/08-decimal/add9-p.bal",
		"subset8/08-decimal/div6-p.bal",
		"subset8/08-decimal/div7-p.bal",
		"subset8/08-decimal/div8-p.bal",
		"subset8/08-decimal/div9-p.bal",
		"subset8/08-decimal/mul10-p.bal",
		"subset8/08-decimal/mul11-p.bal",
		"subset8/08-decimal/mul8-p.bal",
		"subset8/08-decimal/mul9-p.bal",
		"subset8/08-decimal/rem6-p.bal",
		"subset8/08-decimal/rem7-p.bal",
		"subset8/08-decimal/sub5-p.bal",
		"subset8/08-decimal/sub6-p.bal",
		"subset8/08-decimal/toint9-p.bal",
		"subset8/08-list/int2-p.bal",
		"subset8/08-list/int5-p.bal",
		"subset8/08-map/int2-p.bal",
		"subset8/08-fill/fill5-p.bal",
		"subset8/08-fill/fill6-p.bal",

		// Expected panic (-fp / "future panic" suffix): the source documents a runtime panic
		// (e.g. `// @panic bad mapping store`) but pi neither emits any diagnostic nor panics,
		// so the run completes silently.
		"subset8/08-future/fieldlvalue1-fp.bal",

		// Expected output: migrated -v tests where the source documents stdout via
		//  comments, but pi produces empty stdout. The runtime did not
		// emit the expected values, so the test would silently pass with a misleading fixture.
		"subset8/08-inclusive/fieldexpr6-v.bal",
		"subset8/08-inclusive/fieldexpr7-v.bal",
		"subset8/08-inclusive/fieldexpr8-v.bal",
		"subset8/08-nillifting/additive-1-v.bal",
		"subset8/08-nillifting/additive-4-v.bal",
		"subset8/08-nillifting/binary-bitwise-1-v.bal",
		"subset8/08-nillifting/binary-bitwise-4-v.bal",
		"subset8/08-nillifting/multiplicative-1-v.bal",
		"subset8/08-nillifting/multiplicative-4-v.bal",
		"subset8/08-nillifting/shift-1-v.bal",
		"subset8/08-nillifting/shift-4-v.bal",
		"subset8/08-nillifting/unary-5-v.bal",
		"subset8/08-record/fieldexpr9-v.bal",
		"subset8/08-vararg/lib2-v.bal",

		// Invalid output: migrated -v tests where pi emits stdout that differs from the
		// values documented by the source's `// @output ...` markers (i.e. the runtime took
		// the wrong branch).
		"subset8/08-intersect/mapping2-v.bal",
	}

	// Skip project-level integration tests with non-deterministic output.
	skipProjectIntegrationTests = []string{
		"multi-module-same-file-e",
		"inclusive-import1-v",
		"record-import1-v",
		// Migrated from nballerina testSuite/08-import/const4-e: cycle-detection picks a different
		// break point than the upstream compiler, so the reported error path is not stable.
		"import-const4-e",

		// Missing init exp for constants.
		"import-const1-v",

		// Expected error:
		"import-const5-e",
		"import-type3-e",

		// Expected clean run:
		"import-main-v",
		"import-type6-v",
	}
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

type testResult struct {
	success        bool
	expectedStdout string
	actualStdout   string
	expectedStderr string
	actualStderr   string
}

// caseRun is the full result of executing one corpus case (single-file or
// project): captured streams plus the resolved error diagnostics needed for
// `-e` annotation checks.
type caseRun struct {
	stdout string
	stderr string
	diags  []resolvedDiag
}

func TestIntegration(t *testing.T) {
	testPairs := test_util.GetTests(t, test_util.Integration, func(path string) bool {
		return true
	})

	for _, testPair := range testPairs {
		t.Run(testPair.Name, func(t *testing.T) {
			t.Parallel()
			testIntegration(t, testPair)
		})
	}
}

func TestProjectIntegration(t *testing.T) {
	if _, err := os.Stat(corpusProjectBaseDir); os.IsNotExist(err) {
		return
	}

	projectDirs := findProjectDirs(corpusProjectBaseDir)

	for _, projDir := range projectDirs {
		dirName := filepath.Base(projDir)
		txtarPath := filepath.Join(corpusProjectIntegrationBaseDir, dirName+".txtar")

		t.Run(dirName, func(t *testing.T) {
			t.Parallel()
			if isProjectTestSkipped(dirName) {
				t.Skipf("Skipping project integration test for %s", dirName)
			}
			testProjectIntegration(t, dirName, projDir, txtarPath)
		})
	}
}

func testIntegration(t *testing.T, testPair test_util.TestCase) {
	if isTestSkipped(testPair) {
		t.Skipf("Skipping integration test for %s", testPair.InputPath)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while running %s: %v", testPair.InputPath, r)
		}
	}()

	run := runIntegrationCase(testPair.InputPath)
	if *update {
		normalizedStderr := normalizeIntegrationStderr(run.stderr)
		if test_util.UpdateTxtarArchiveIfNeeded(t, testPair.ExpectedPath, test_util.TxtarFilesStdoutStderr(run.stdout, normalizedStderr)) {
			t.Fatalf("Updated expected file: %s", testPair.ExpectedPath)
		}
		return
	}

	expectedStdout, expectedStderr, err := test_util.LoadTxtarStdoutStderr(testPair.ExpectedPath)
	if err != nil {
		t.Fatalf("failed to load expected from %s: %v", testPair.ExpectedPath, err)
	}

	result := evaluateTestResult(expectedStdout, expectedStderr, run.stdout, run.stderr)
	assertAnnotations(t, collectSingleFileSources(testPair.InputPath), testPair.Name, run.stdout, run.stderr, run.diags)
	if result.success {
		return
	}

	stdoutMismatch := result.expectedStdout != result.actualStdout
	stderrMismatch := result.expectedStderr != normalizeIntegrationStderr(result.actualStderr)

	var msg strings.Builder
	if stdoutMismatch {
		fmt.Fprintf(&msg, "stdout mismatch\n%s", test_util.FormatExpectedGot(result.expectedStdout, result.actualStdout))
	}
	if stderrMismatch {
		if msg.Len() > 0 {
			msg.WriteString("\n\n")
		}
		fmt.Fprintf(&msg, "stderr mismatch\n%s", test_util.FormatExpectedGot(
			normalizeIntegrationStderr(result.expectedStderr),
			normalizeIntegrationStderr(result.actualStderr),
		))
	}
	t.Errorf("%s", msg.String())
}

func splitStderrDiagnostics(stderr string) []string {
	var diagnostics []string
	for part := range strings.SplitSeq(stderr, "\n\n") {
		diagnostic := strings.TrimSpace(part)
		if diagnostic != "" {
			diagnostics = append(diagnostics, diagnostic)
		}
	}
	return diagnostics
}

func normalizeIntegrationStderr(stderr string) string {
	stderr = strings.TrimSpace(stderr)
	if stderr == "" {
		return ""
	}

	diagnostics := splitStderrDiagnostics(stderr)

	slices.Sort(diagnostics)
	return strings.Join(diagnostics, "\n\n") + "\n"
}

func isTestSkipped(tc test_util.TestCase) bool {
	return isSkipKey(filepath.ToSlash(tc.Name))
}

func isSkipKey(key string) bool {
	return slices.Contains(skipIntegrationTests, key)
}

func isProjectTestSkipped(dirName string) bool {
	return slices.Contains(skipProjectIntegrationTests, dirName)
}

func resolveErrorDiagnostics(result projects.DiagnosticResult, de *diagnostics.DiagnosticEnv) []resolvedDiag {
	errs := result.Errors()
	if len(errs) == 0 {
		return nil
	}
	out := make([]resolvedDiag, 0, len(errs))
	for _, d := range errs {
		loc := d.Location()
		if !diagnostics.LocationHasSource(loc) {
			continue
		}
		out = append(out, resolvedDiag{
			file:      de.FileName(loc),
			startLine: de.StartLine(loc) + 1,
			endLine:   de.EndLine(loc) + 1,
		})
	}
	return out
}

func runIntegrationCase(balFile string) caseRun {
	var stdoutBuf, stderrBuf bytes.Buffer

	birPkg, diags, compileErr := runCompilePhase(balFile, &stdoutBuf, &stderrBuf)
	if birPkg == nil || compileErr != nil {
		return caseRun{stdout: stdoutBuf.String(), stderr: stderrBuf.String(), diags: diags}
	}

	runInterpretPhase(birPkg, &stdoutBuf, &stderrBuf)
	return caseRun{stdout: stdoutBuf.String(), stderr: stderrBuf.String(), diags: diags}
}

func evaluateTestResult(expectedStdout, expectedStderr, actualStdout, actualStderr string) testResult {
	stderrMatch := expectedStderr == normalizeIntegrationStderr(actualStderr)
	return testResult{
		success:        actualStdout == expectedStdout && stderrMatch,
		expectedStdout: expectedStdout,
		actualStdout:   actualStdout,
		expectedStderr: expectedStderr,
		actualStderr:   actualStderr,
	}
}

func runCompilePhase(balFile string, stdoutBuf, stderrBuf *bytes.Buffer) (pkg *bir.BIRPackage, diags []resolvedDiag, err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%v", r)
			msg = strings.TrimPrefix(msg, panicPrefix)
			fmt.Fprintf(stdoutBuf, "%s%s\n", panicPrefix, msg)
			err = fmt.Errorf("compile panic")
		}
	}()

	fsys := os.DirFS(filepath.Dir(balFile))

	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		fmt.Fprintf(stdoutBuf, "%s\n", err.Error())
		return nil, nil, err
	}
	ballerinaEnvFs := os.DirFS(ballerinaEnvPath)

	result, err := projects.Load(fsys, filepath.Base(balFile), projects.ProjectLoadConfig{
		BallerinaEnvFs: ballerinaEnvFs,
	})
	if err != nil {
		fmt.Fprintf(stdoutBuf, "%s\n", err.Error())
		return nil, nil, err
	}
	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()

	printDiagnostics(fsys, stderrBuf, compilation.DiagnosticResult(), compilation.DiagnosticEnv())
	diags = resolveErrorDiagnostics(compilation.DiagnosticResult(), compilation.DiagnosticEnv())
	if compilation.DiagnosticResult().HasErrors() {
		return nil, diags, nil
	}

	backend := projects.NewBallerinaBackend(compilation)
	return backend.BIR(), diags, nil
}

func runInterpretPhase(birPkg *bir.BIRPackage, stdoutBuf, stderrBuf *bytes.Buffer) {
	if birPkg == nil {
		return
	}

	rt := runtime.NewRuntime(test_util.TestPal(stdoutBuf, stderrBuf))
	if err := rt.Interpret(*birPkg); err != nil {
		// For now just write the error string to stderr to match corpus expectations
		fmt.Fprintln(stderrBuf, err.Error())
	}
}

func findProjectDirs(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var dirs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, "-v") || strings.HasSuffix(name, "-e") || strings.HasSuffix(name, "-p") {
			dirs = append(dirs, filepath.Join(dir, name))
		}
	}
	return dirs
}

func testProjectIntegration(t *testing.T, dirName, projDir, txtarPath string) {
	if isSkipKey("project/" + dirName) {
		t.Skipf("Skipping project integration test for %s", dirName)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while running %s: %v", dirName, r)
		}
	}()

	run := runProjectIntegrationCase(projDir)
	if *update {
		normalizedStderr := normalizeIntegrationStderr(run.stderr)
		if test_util.UpdateTxtarArchiveIfNeeded(t, txtarPath, test_util.TxtarFilesStdoutStderr(run.stdout, normalizedStderr)) {
			t.Fatalf("Updated expected file: %s", txtarPath)
		}
		return
	}

	expectedStdout, expectedStderr, err := test_util.LoadTxtarStdoutStderr(txtarPath)
	if err != nil {
		t.Fatalf("failed to load expected from %s: %v", txtarPath, err)
	}

	result := evaluateTestResult(expectedStdout, expectedStderr, run.stdout, run.stderr)

	projectSources, srcErr := collectProjectSources(projDir)
	if srcErr != nil {
		t.Errorf("failed to collect project sources: %v", srcErr)
	} else {
		assertAnnotations(t, projectSources, dirName, run.stdout, run.stderr, run.diags)
	}
	if result.success {
		return
	}

	stdoutMismatch := result.expectedStdout != result.actualStdout
	stderrMismatch := result.expectedStderr != result.actualStderr

	var msg strings.Builder
	if stdoutMismatch {
		fmt.Fprintf(&msg, "stdout mismatch\n%s", test_util.FormatExpectedGot(
			result.expectedStdout,
			result.actualStdout,
		))
	}
	if stderrMismatch {
		if msg.Len() > 0 {
			msg.WriteString("\n\n")
		}
		fmt.Fprintf(&msg, "stderr mismatch\n%s", test_util.FormatExpectedGot(
			result.expectedStderr,
			normalizeIntegrationStderr(result.actualStderr),
		))
	}
	t.Errorf("%s", msg.String())
}

func runProjectIntegrationCase(projectDir string) caseRun {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	birPkgs, diags, compileErr := runProjectCompilePhase(projectDir, &stdoutBuf, &stderrBuf)
	if birPkgs == nil || compileErr != nil {
		return caseRun{stdout: stdoutBuf.String(), stderr: stderrBuf.String(), diags: diags}
	}

	runProjectInterpretPhase(birPkgs, &stdoutBuf, &stderrBuf)
	return caseRun{stdout: stdoutBuf.String(), stderr: stderrBuf.String(), diags: diags}
}

func runProjectCompilePhase(projectDir string, stdoutBuf, stderrBuf *bytes.Buffer) (pkgs []*bir.BIRPackage, diags []resolvedDiag, err error) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%v", r)
			msg = strings.TrimPrefix(msg, panicPrefix)
			fmt.Fprintf(stdoutBuf, "%s%s\n", panicPrefix, msg)
			err = fmt.Errorf("compile panic")
		}
	}()

	fsys := os.DirFS(projectDir)

	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		fmt.Fprintf(stdoutBuf, "%s\n", err.Error())
		return nil, nil, err
	}
	ballerinaEnvFs := os.DirFS(ballerinaEnvPath)

	result, err := projects.Load(fsys, ".", projects.ProjectLoadConfig{
		BallerinaEnvFs: ballerinaEnvFs,
	})
	if err != nil {
		fmt.Fprintf(stdoutBuf, "%s\n", err.Error())
		return nil, nil, err
	}
	currentPkg := result.Project().CurrentPackage()
	compilation := currentPkg.Compilation()

	printDiagnostics(fsys, stderrBuf, compilation.DiagnosticResult(), compilation.DiagnosticEnv())
	diags = resolveErrorDiagnostics(compilation.DiagnosticResult(), compilation.DiagnosticEnv())
	if compilation.DiagnosticResult().HasErrors() {
		return nil, diags, nil
	}

	backend := projects.NewBallerinaBackend(compilation)
	return backend.BIRPackages(), diags, nil
}

func runProjectInterpretPhase(birPkgs []*bir.BIRPackage, stdoutBuf, stderrBuf *bytes.Buffer) {
	if len(birPkgs) == 0 {
		return
	}

	rt := runtime.NewRuntime(test_util.TestPal(stdoutBuf, stderrBuf))
	for _, birPkg := range birPkgs {
		if err := rt.Interpret(*birPkg); err != nil {
			fmt.Fprintln(stderrBuf, err.Error())
			return
		}
	}
}

func TestProjectSerializationRoundtrip(t *testing.T) {
	flag.Parse()

	if _, err := os.Stat(corpusProjectBaseDir); os.IsNotExist(err) {
		return
	}

	projectDirs := findProjectDirs(corpusProjectBaseDir)

	for _, projDir := range projectDirs {
		dirName := filepath.Base(projDir)
		if !strings.HasSuffix(dirName, "-v") {
			continue
		}
		txtarPath := filepath.Join(corpusProjectIntegrationBaseDir, dirName+".txtar")

		t.Run(dirName, func(t *testing.T) {
			t.Parallel()
			// Roundtrip test reuses the integration project skip list because any project
			// skipped at the integration level has no usable expected fixture.
			if isProjectTestSkipped(dirName) {
				t.Skipf("Skipping project serialization roundtrip for %s", dirName)
			}
			testProjectSerializationRoundtrip(t, dirName, projDir, txtarPath)
		})
	}
}

func testProjectSerializationRoundtrip(t *testing.T, dirName, projDir, txtarPath string) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic while running %s: %v", dirName, r)
		}
	}()

	expectedStdout, expectedStderr, err := test_util.LoadTxtarStdoutStderr(txtarPath)
	if err != nil {
		t.Fatalf("failed to load expected from %s: %v", txtarPath, err)
	}

	stdout, stderr := runProjectSerializationRoundtrip(projDir)
	result := evaluateTestResult(expectedStdout, expectedStderr, stdout, stderr)
	if result.success {
		return
	}

	stdoutMismatch := result.expectedStdout != result.actualStdout
	stderrMismatch := result.expectedStderr != result.actualStderr

	var msg strings.Builder
	if stdoutMismatch {
		fmt.Fprintf(&msg, "stdout mismatch\n%s", test_util.FormatExpectedGot(result.expectedStdout, result.actualStdout))
	}
	if stderrMismatch {
		if msg.Len() > 0 {
			msg.WriteString("\n\n")
		}
		fmt.Fprintf(&msg, "stderr mismatch\n%s", test_util.FormatExpectedGot(result.expectedStderr, result.actualStderr))
	}
	t.Errorf("%s", msg.String())
}

func runProjectSerializationRoundtrip(projectDir string) (stdout, stderr string) {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	absProjectDir, err := filepath.Abs(projectDir)
	if err != nil {
		fmt.Fprintf(&stdoutBuf, "%s\n", err.Error())
		return stdoutBuf.String(), stderrBuf.String()
	}

	fsys := os.DirFS(projectDir)
	ballerinaEnvPath, err := getBallerinaEnvPath()
	if err != nil {
		fmt.Fprintf(&stdoutBuf, "%s\n", err.Error())
		return stdoutBuf.String(), stderrBuf.String()
	}
	ballerinaEnvFs := os.DirFS(ballerinaEnvPath)
	result, err := projects.Load(fsys, ".", projects.ProjectLoadConfig{
		BallerinaEnvFs: ballerinaEnvFs,
	})
	if err != nil {
		fmt.Fprintf(&stdoutBuf, "%s\n", err.Error())
		return stdoutBuf.String(), stderrBuf.String()
	}
	project := result.Project()
	currentPkg := project.CurrentPackage()
	compilation := currentPkg.Compilation()

	printDiagnostics(fsys, &stderrBuf, compilation.DiagnosticResult(), compilation.DiagnosticEnv())
	if compilation.DiagnosticResult().HasErrors() {
		return stdoutBuf.String(), stderrBuf.String()
	}

	backend := projects.NewBallerinaBackend(compilation)
	birPkgs := backend.BIRPackages()
	exportedSymbols := backend.ExportedSymbols()

	if len(birPkgs) == 0 {
		return stdoutBuf.String(), stderrBuf.String()
	}

	deps := birPkgs[:len(birPkgs)-1]

	// Step 1: Serialize dep symbols and BIR to byte arrays
	type serializedModule struct {
		symBytes []byte
		birBytes []byte
	}
	serializedDeps := make([]serializedModule, 0, len(deps))

	for _, dep := range deps {
		pkgIdent := semantics.PackageIdentifier{
			OrgName:    dep.PackageID.OrgName.Value(),
			ModuleName: dep.PackageID.PkgName.Value(),
		}
		exported, ok := exportedSymbols[pkgIdent]
		if !ok {
			fmt.Fprintf(&stdoutBuf, "exported symbols not found for %s/%s\n", pkgIdent.OrgName, pkgIdent.ModuleName)
			return stdoutBuf.String(), stderrBuf.String()
		}

		symBytes, err := symbolpool.Marshal(exported, dep.TypeEnv)
		if err != nil {
			fmt.Fprintf(&stdoutBuf, "symbol serialization failed: %v\n", err)
			return stdoutBuf.String(), stderrBuf.String()
		}

		birBytes, err := bircodec.Marshal(dep)
		if err != nil {
			fmt.Fprintf(&stdoutBuf, "BIR serialization failed: %v\n", err)
			return stdoutBuf.String(), stderrBuf.String()
		}

		serializedDeps = append(serializedDeps, serializedModule{symBytes: symBytes, birBytes: birBytes})
	}

	// Step 2: Create fresh compiler and deserialize dep symbols + BIR
	freshEnv := context.NewCompilerEnvironment(semtypes.CreateTypeEnv(), false)
	publicSymbols := make(map[semantics.PackageIdentifier]model.ExportedSymbolSpace)
	deserialized := make([]*bir.BIRPackage, 0, len(birPkgs))

	for i, sd := range serializedDeps {
		exported, err := symbolpool.Unmarshal(freshEnv, sd.symBytes)
		if err != nil {
			fmt.Fprintf(&stdoutBuf, "symbol deserialization failed: %v\n", err)
			return stdoutBuf.String(), stderrBuf.String()
		}

		dep := deps[i]
		pkgIdent := semantics.PackageIdentifier{
			OrgName:    dep.PackageID.OrgName.Value(),
			ModuleName: dep.PackageID.PkgName.Value(),
		}
		publicSymbols[pkgIdent] = exported

		freshCtx := context.NewCompilerContext(freshEnv)
		deserializedPkg, err := bircodec.Unmarshal(freshCtx, sd.birBytes)
		if err != nil {
			fmt.Fprintf(&stdoutBuf, "BIR deserialization failed: %v\n", err)
			return stdoutBuf.String(), stderrBuf.String()
		}

		deserialized = append(deserialized, deserializedPkg)
	}

	// Step 3: Recompile the main (default) module from source using deserialized dep symbols
	defaultModule := currentPkg.DefaultModule()
	defaultDesc := defaultModule.Descriptor()
	defaultOrg := defaultDesc.Org().Value()

	mainBirPkg, err := compileModuleFromSource(freshEnv, project, defaultModule, absProjectDir, publicSymbols, defaultOrg)
	if err != nil {
		fmt.Fprintf(&stdoutBuf, "main module recompilation failed: %v\n", err)
		return stdoutBuf.String(), stderrBuf.String()
	}

	deserialized = append(deserialized, mainBirPkg)

	runProjectInterpretPhase(deserialized, &stdoutBuf, &stderrBuf)
	return stdoutBuf.String(), stderrBuf.String()
}

func compileModuleFromSource(env *context.CompilerEnvironment, project projects.Project, module *projects.Module,
	absProjectDir string, publicSymbols map[semantics.PackageIdentifier]model.ExportedSymbolSpace, defaultOrg string,
) (*bir.BIRPackage, error) {
	cx := context.NewCompilerContext(env)

	// Register source files with DiagnosticEnv
	de := cx.DiagnosticEnv()
	for _, docID := range module.DocumentIDs() {
		relPath := project.DocumentPath(docID)
		absPath := filepath.Join(absProjectDir, relPath)
		content, err := os.ReadFile(absPath)
		if err == nil {
			de.RegisterFile(absPath, text.NewStringTextDocument(string(content)))
		}
	}

	// Parse all source files in the module
	docIDs := module.DocumentIDs()
	var syntaxTrees []*ast.BLangCompilationUnit
	for _, docID := range docIDs {
		relPath := project.DocumentPath(docID)
		absPath := filepath.Join(absProjectDir, relPath)
		st, err := parser.GetSyntaxTree(cx, absPath)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %v", relPath, err)
		}
		cu := ast.GetCompilationUnit(cx, st)
		syntaxTrees = append(syntaxTrees, cu)
	}

	// Build package from compilation units
	var pkg *ast.BLangPackage
	if len(syntaxTrees) == 1 {
		pkg = ast.ToPackage(syntaxTrees[0])
	} else {
		pkg = &ast.BLangPackage{}
		for _, cu := range syntaxTrees {
			if pkg.PackageID == nil {
				pkg.PackageID = cu.GetPackageID()
			}
			for _, node := range cu.GetTopLevelNodes() {
				switch n := node.(type) {
				case *ast.BLangImportPackage:
					pkg.Imports = append(pkg.Imports, *n)
				case *ast.BLangConstant:
					pkg.Constants = append(pkg.Constants, *n)
				case *ast.BLangService:
					pkg.Services = append(pkg.Services, *n)
				case *ast.BLangFunction:
					pkg.Functions = append(pkg.Functions, *n)
				case *ast.BLangTypeDefinition:
					pkg.TypeDefinitions = append(pkg.TypeDefinitions, *n)
				case *ast.BLangAnnotation:
					pkg.Annotations = append(pkg.Annotations, *n)
				case *ast.BLangClassDefinition:
					pkg.ClassDefinitions = append(pkg.ClassDefinitions, *n)
				default:
					pkg.TopLevelNodes = append(pkg.TopLevelNodes, node)
				}
			}
		}
	}

	// Set the package ID to match the module descriptor
	desc := module.Descriptor()
	orgName := model.Name(desc.Org().Value())
	moduleName := desc.Name().String()
	nameComps := make([]model.Name, 0)
	for _, part := range strings.Split(moduleName, ".") {
		nameComps = append(nameComps, model.Name(part))
	}
	version := model.Name(desc.Version().String())
	if version == "" {
		version = model.DEFAULT_VERSION
	}
	pkg.PackageID = cx.NewPackageID(orgName, nameComps, version)

	// Run compilation pipeline
	importedSymbols := semantics.ResolveImports(cx, pkg, semantics.GetImplicitImports(cx), publicSymbols, defaultOrg)
	semantics.ResolveSymbols(cx, pkg, importedSymbols)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("symbol resolution failed")
	}

	semantics.ResolveTopLevelNodes(cx, pkg, importedSymbols)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("top-level type resolution failed")
	}

	semantics.ResolveLocalNodes(cx, pkg, importedSymbols)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("local type resolution failed")
	}

	analyzer := semantics.NewSemanticAnalyzer(cx)
	analyzer.Analyze(pkg)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("semantic analysis failed")
	}

	cfg := semantics.CreateControlFlowGraph(cx, pkg)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("CFG creation failed")
	}

	semantics.AnalyzeCFG(cx, pkg, cfg)
	if cx.HasDiagnostics() {
		return nil, fmt.Errorf("CFG analysis failed")
	}

	pkg = desugar.DesugarPackage(cx, pkg, importedSymbols)

	return bir.GenBir(cx, pkg), nil
}

var skipBenchmarkIntegrationTests = []string{
	// error: interface conversion: values.BalValue is nil, not string
	"08-bench/map-v.bal",
}

func BenchmarkIntegration(b *testing.B) {
	testPairs := test_util.GetTests(b, test_util.Bench, func(path string) bool {
		return true
	})
	for _, testPair := range testPairs {
		if slices.Contains(skipBenchmarkIntegrationTests, filepath.ToSlash(testPair.Name)) {
			continue
		}
		b.Run(testPair.Name, func(b *testing.B) {
			expectedStdout, expectedStderr, err := test_util.LoadTxtarStdoutStderr(testPair.ExpectedPath)
			if err != nil {
				b.Fatalf("failed to load expected from %s: %v", testPair.ExpectedPath, err)
			}

			var run caseRun
			b.ResetTimer()
			for b.Loop() {
				run = runIntegrationCase(testPair.InputPath)
			}
			b.StopTimer()

			result := evaluateTestResult(expectedStdout, expectedStderr, run.stdout, run.stderr)
			if !result.success {
				b.Fatalf("output mismatch for %s:\nstdout:\n%s\nstderr:\n%s",
					testPair.InputPath,
					test_util.FormatExpectedGot(result.expectedStdout, result.actualStdout),
					test_util.FormatExpectedGot(
						normalizeIntegrationStderr(result.expectedStderr),
						normalizeIntegrationStderr(result.actualStderr),
					))
			}
		})
	}
}

func getBallerinaEnvPath() (string, error) {
	if balEnv := os.Getenv(projects.BallerinaEnvVar); balEnv != "" {
		return balEnv, nil
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userHome, projects.UserHomeDirName), nil
}
