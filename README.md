# Filler
## Description
Fills golang objects with predefined values
## Usage
All fields must be exported or you can define default value for type with unexported fields
```go
    var obj = struct {
        M map[string]struct {
            Name  string
            Value int
        }
        S []struct {
            Name  string
            Value int
        }
    }{}
    f := filler.New()
    f.RegisterType("exampleString")
    f.RegisterType(int(1))
    if err := f.Fill(&obj); err != nil {
        panic(err)
    }
    fmt.Printf("result: %+v", obj) 
    // result: {M:map[exampleString:{Name:exampleString Value:1}] S:[{Name:exampleString Value:1}]}
```