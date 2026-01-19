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

package elements

type AttachPointPoint uint8

const (
	ATTACH_POINT_POINT_TYPE AttachPointPoint = iota
	ATTACH_POINT_POINT_OBJECT
	ATTACH_POINT_POINT_FUNCTION
	ATTACH_POINT_POINT_OBJECT_METHOD
	ATTACH_POINT_POINT_SERVICE_REMOTE
	ATTACH_POINT_POINT_PARAMETER
	ATTACH_POINT_POINT_RETURN
	ATTACH_POINT_POINT_SERVICE
	ATTACH_POINT_POINT_FIELD
	ATTACH_POINT_POINT_OBJECT_FIELD
	ATTACH_POINT_POINT_RECORD_FIELD
	ATTACH_POINT_POINT_LISTENER
	ATTACH_POINT_POINT_ANNOTATION
	ATTACH_POINT_POINT_EXTERNAL
	ATTACH_POINT_POINT_VAR
	ATTACH_POINT_POINT_CONST
	ATTACH_POINT_POINT_WORKER
	ATTACH_POINT_POINT_CLASS
)

var attachPointPointValues = map[AttachPointPoint]string{
	ATTACH_POINT_POINT_TYPE:           "type",
	ATTACH_POINT_POINT_OBJECT:         "objecttype",
	ATTACH_POINT_POINT_FUNCTION:       "function",
	ATTACH_POINT_POINT_OBJECT_METHOD:  "objectfunction",
	ATTACH_POINT_POINT_SERVICE_REMOTE: "serviceremotefunction",
	ATTACH_POINT_POINT_PARAMETER:      "parameter",
	ATTACH_POINT_POINT_RETURN:         "return",
	ATTACH_POINT_POINT_SERVICE:        "service",
	ATTACH_POINT_POINT_FIELD:          "field",
	ATTACH_POINT_POINT_OBJECT_FIELD:   "objectfield",
	ATTACH_POINT_POINT_RECORD_FIELD:   "recordfield",
	ATTACH_POINT_POINT_LISTENER:       "listener",
	ATTACH_POINT_POINT_ANNOTATION:     "annotation",
	ATTACH_POINT_POINT_EXTERNAL:       "external",
	ATTACH_POINT_POINT_VAR:            "var",
	ATTACH_POINT_POINT_CONST:          "const",
	ATTACH_POINT_POINT_WORKER:         "worker",
	ATTACH_POINT_POINT_CLASS:          "class",
}

func (p AttachPointPoint) GetValue() string {
	return attachPointPointValues[p]
}

func (p AttachPointPoint) String() string {
	return p.GetValue()
}

type AttachPoint struct {
	Point  AttachPointPoint
	Source bool
}

func NewAttachPoint(point AttachPointPoint, source bool) AttachPoint {
	return AttachPoint{
		Point:  point,
		Source: source,
	}
}

func GetAttachmentPoint(value string, source bool) *AttachPoint {
	for point, val := range attachPointPointValues {
		if val == value {
			ap := NewAttachPoint(point, source)
			return &ap
		}
	}
	return nil
}
