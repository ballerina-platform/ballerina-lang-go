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

# Represents Foo object
# 
# + i - Expected iter
public class Foo {
    public int i = 1;
}

# Represents Bar object
# 
# + i - Expected iter
# + s - String str
# + foos - Foo object inside 
public class Bar {
    public int i = 1;
    public string s = "str";
    public Foo foos = new;
}

# Represents student object
# 
# + f - Bar type object cast inside student
# + name - Name of the student
public class Student {
    public Foo f = new Bar();
    public string name = "John";
}
