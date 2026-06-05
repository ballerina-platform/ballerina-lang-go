type Info record {|
    int code;
|};

annotation Info info on type;

@info {code: "bad"} // @error
type Person record {|
    string name;
|};
