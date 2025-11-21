# go-fault
`go-fault` provides excellent error handling features to your application.<br>
<br>
You can...<br>
- Add stack trace
- Add ErrorType and Compare errors by ErrorType
- Bind requestId
- Bind context data with error
- Convert error into JSON string



## Compatibility with standard library errors and other libraries

| Method          | Library                      | Compatibility |
|-----------------|------------------------------|---------|
| errors.Is()     | Go standard package "errors" | ✅       |
| errors.As()     | Go standard package "errors" | ✅       |
| errors.Unwrap() | Go standard package "errors" | ✅       |   
| errors.Join()   | Go standard package "errors" | ✅       |
| errors.Wrap()   | pkg/errors                   | ✅       |
| errors.Cause()  | pkg/errors                   | ✅       |

## How to use

### `New()`

`New()` return an error with stack trace.

```go
import (
"fmt"
"github.com/hinoguma/go-fault/errors"
)

// error with stack trace
err := errors.New("example error")
fmt.Printf("%+v", err)
// Output:
// main_error:
//     message: example error
//     type: none
//     stack trace:
//         example.exampleFunction() /path/to/your/file.go:15
//         example.exampleFunction2() /path/to/your/file.go:20
```

<br>

### `Wrap()`

`Wrap()` adds `stack trace` to existing error if it doesn't have one.

```go
// standard library error
originalErr := fmt.Errorf("original error") 

// add stack trace
wrappedErr = errors.Wrap(originalErr, "wrapped error")

fmt.Printf("%+v", wrappedErr)
// Output:
// main_error:
//     message: wrapped error: original error
//     type: none
//     stack trace:
//         example.exampleFunction() /path/to/your/file.go:15
//         example.exampleFunction2() /path/to/your/file.go:20
```


if the error already has stack trace,` Wrap()` does **not add new stack trace**.

```go
import 	"github.com/hinoguma/go-fault/errors"

// go-fault/errors.New() return error with stack trace
errWithStack := errors.New("error with stack")

// no new stack trace added, just wrap error
errors.Wrap(errWithStack, "another wrap")
```

<br>

### `Lift()`

Wrap() needs message parameter.<br>
**if you bather to o it every time**, you can use `Lift()` instead of Wrap().

`Lift()` just **adds stack trace** to existing error if it doesn't have one.

```go
import 	"github.com/hinoguma/go-fault/errors"

// add stack trace to existing error
originalErr := fmt.Errorf("original error")

// adding stack trace
errors.Lift(originalErr)

// not adding stack trace if error already has one
errWithStack := errors.New("error with stack")

// no new stack trace added
errors.Lift(errWithStack)
```

<br>

### Error Type

You can **add type to error** and branch your error handling logic based on error type.

```go
const CustomType1 fault.ErrorType = "CustomType1"

// Set error type as CustomType1 
err := errors.New("error with type")
err = errors.With(err).Type(CustomType1).Err()
```
<br>

 You can check error type using `IsType()`.<br>
`errors.Is()` checks identity of errors but `IsType()` check **error type only**.
```go
if errors.IsType(err, CustomType1) {
    fmt.Println("This is a CustomType1 error")
}

// if wrapping error, IsType() can still check error type in the chain
wrappedErr := fmt.Errorf("wrapping error: %w", err)
if errors.IsType(wrappedErr, CustomType1) {
	    fmt.Println("This is a CustomType1 error")
}
```


<br>

### Additional Context

You can **add context** to error using `With()`.<br>
`With()` provides various methods to add context to error and **methods can be chained**.

```go
err = errors.With(err).
    Type(CustomType1).
    When(time.Now()).
    RequestID("request-1234").
    AddTagString("key1", "value1").
    Err()
```

#### When it happened? Which Request?

Use `When()` and `RequestID()`.

```go
err = errors.With(err).
	When(time.Now()).
	RequestID("request-1234").
    Err()

fmt.Printf("%+v", err)
// Output:
// main_error:
//     message: example error
//     type: none
//     when: 2024-01-01T12:00:00Z
//     request id: request-1234
//     stack trace:
//         example.exampleFunction() /path/to/your/file.go:15
//         example.exampleFunction2() /path/to/your/file.go:20
```

<br>

##### More Context
Use `AddTagXXX()` to add more context to error.
```go
err = errors.With(err).
	AddTagString("key1", "value1").
	AddTagInt("key2", 42).
	AddTagFloat("key3", 3.14).
	AddTagBool("key4", true).
	Err()

fmt.Printf("%+v", err)
// Output:
// main_error:
//     message: example error
//     type: none
//     tags:
//         key1: value1
//         key2: 42
//         key3: 3.14
//         key4: true
//     stack trace:
//         example.exampleFunction() /path/to/your/file.go:15
//         example.exampleFunction2() /path/to/your/file.go:20
```

<br>

### Logging

#### Log as JSON string
type, message, stacktrace are always included in the JSON output.<br>
when, request_id, tags are included only if they are set.

```go
js := errors.ToJsonString(err)
fmt.Println(js)
// Output:
// {"type":"none","error":"go standard error","when":"2024-06-01T12:00:00Z","request_id":"12345","tags":{"tag1":"value1","key2":42,"key3":3.14,"key4":true},"stacktrace":[{"file":"main.go","line":75,"function":"github.com/hinoguma/go-fault.main"}],"sub_errors":[{"type":"none","message":"go standard error","stacktrace":[]},{"type":"none","message":"go standard error2","stacktrace":[]}]}

// Output is one-liner JSON string.
// Below is pretty printed version.
{
  "type": "none",
  "error": "go standard error",
  "when": "2024-06-01T12:00:00Z",
  "request_id": "12345",
  "tags": {
    "tag1": "value1",
	"key2":42,
	"key3":3.14,
	"key4":true
  },
  "stacktrace": [
    {
      "file": "main.go",
      "line": 75,
      "function": "github.com/hinoguma/go-fault.main"
    }
  ],
  "sub_errors": [
    {
      "type": "none",
      "message": "go standard error",
      "stacktrace": []
    },
    {
      "type": "none",
      "message": "go standard error2",
      "stacktrace": []
    }
  ]
}
```

#### Log as plain string
use fmt with `%+v` verb to print error with all details including stack trace.
```go
txt := fmt.SprintF("%+v\n", err)
// txt is:
// main_error:
//     message: example error
//     type: none
//     tags:
//         key1: value1
//         key2: 42
//         key3: 3.14
//         key4: true
//     stack trace:
//         example.exampleFunction() /path/to/your/file.go:15
//         example.exampleFunction2() /path/to/your/file.go:20
```


