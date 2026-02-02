## Coding style

- Don't make symbols public unless asked for or needed
- Constructor methods should data for all the fields unless their is default initialization
    - Map values should always be initialized to an empty map

- If multiple structs needs to holds same set of fields and implement methods on those fields add *Base struct and use type inclusion on other structs
    - Make this base struct private
    - Implement the relevant methods on the base struct
