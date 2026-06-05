import ballerina/io;

type Info record {|
    string name;
    int code;
|};

annotation Info info on type;
annotation marker on type;

@info {name: "person", code: 7}
@marker
type Person record {|
    string name;
|};

public function main() {
    Info? infoValue = Person.@info;
    if infoValue is Info {
        io:println(infoValue.name); // @output person
        io:println(infoValue.code); // @output 7
    }

    boolean? markerValue = Person.@marker;
    if markerValue is boolean {
        io:println(markerValue); // @output true
    }
}
