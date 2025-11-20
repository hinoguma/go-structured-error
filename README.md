# go-fault

Go Error library with Flexibility and Functionality

## Features
| Feature                                             | go-fault |
|-----------------------------------------------------|----------|
| stack trace                                         | ✅        |
| add error type                                      | ✅        |
| errors.IsType()                                     | ✅        |
| add time when error occurred                        | ✅        |
| add request id to trace which request               | ✅        |
| additional context data you want to bind with error | ✅        |
| convert error into JSON string for logging          | ✅        |

## Compatibility with standard library errors and other libraries

| Method          | library                      | compatibility |
|-----------------|------------------------------|---------------|
| errors.Is()     | Go standard package "errors" | ✅             |
| errors.As()     | Go standard package "errors" | ✅             |
| errors.Unwrap() | Go standard package "errors" | ✅             |   
| errors.Join()   | Go standard package "errors" | ✅             |
| errors.Wrap()   | pkg/errors                   | ✅             |
| errors.Cause()  | pkg/errors                   | No support    |

## How to use

#### New()

New() return an error with stack trace.

```go
import (
"fmt"
"github.com/hinoguma/go-fault/errors"
)

// error with stack trace
err := errors.New("example error")
fmt.Printf("%+v", err)
// Output:
// [Type:none] [Message:example error]
// Stack trace:
// main.main
//     /path/to/your/file.go:10
```

<br>

#### Wrap()

Wrap() adds stack trace to existing error if it doesn't have one.

```go
// add stack trace to existing error
originalErr := fmt.Errorf("original error") // standard library error
originalErr = errors.Wrap(originalErr, "wrapped error") // adding stack trace

fmt.Printf("%+v", originalErr)
// Output:
// [Type:none] [Message:wrapped error: original error]
// Stack trace:
// main.main
//     /path/to/your/file.go:15
```

if the error already has stack trace, Wrap() does not add new stack trace.

```go
// not adding stack trace if error already has one
errWithStack := errors.New("error with stack")
errors.Wrap(errWithStack, "another wrap") // no new stack trace added
```

<br>

#### Lift()

Wrap() needs message parameter.<br>
if you bather that and stack trace upto error occur is enough, you can use Lift() instead of Wrap().

Lift() just adds stack trace to existing error if it doesn't have one.

```go
// add stack trace to existing error
originalErr := fmt.Errorf("original error")
errors.Lift(originalErr) // adding stack trace

// not adding stack trace if error already has one
errWithStack := errors.New("error with stack")
errors.Lift(errWithStack) // no new stack trace added
```

<br>

#### Error Type

You can add type to error and branch your error handling logic based on error type.

```go
const CustomType1 fault.ErrorType = "CustomType1"

// Set error type as CustomType1 
err := errors.New("error with type")
err = errors.With(err).Type(CustomType1).Err()
```

 You can check error type using IsType().
 errors.Is() check  identity ob errors but IsType() just check error type.
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

#### Additional Context

You can add context to error using With().<br>
With() provides various methods to add context to error and methods can be chained.

```go
err = errors.With(err).
    Type(CustomType1).
    When(time.Now()).
    RequestID("request-1234").
    AddTagString("key1", "value1").
    Err()
```

##### When it happened? Which Request?

Use When() and RequestID().

```go
err = errors.With(err).
	When(time.Now()).
	RequestID("request-1234").
    Err()

fmt.Printf("%+v", err)
// Output:
// [Type:none] [Message:error with type]
```

<br>

##### More Context
Use AddTagXXX() to add more context to error.
```go
err = errors.With(err).
	AddTagString("key1", "value1").
	AddTagInt("key2", 42).
	AddTagFloat("key3", 3.14).
	AddTagBool("key4", true).
	Err()

fmt.Printf("%+v", err)
// Output:
// [Type:none] [Message:error with type]
```

<br>

#### Logging

##### Log as JSON string
type, message, stacktrace are always included in the JSON output.<br>
when, request_id, tags are included only if they are set.

```go
js := errors.Converter(err).JsonString()
fmt.Println(js)
// Output:
// {"type":"CustomType1","message":"error with type","when":"2024-06-01T12:00:00Z","request_id":"request-1234","tags":{"key1":"value1","key2":42,"key3":3.14,"key4":true}}

// Output is one-liner JSON string.
// Below is pretty printed version.
// {
//   "type": "CustomType1",
//   "message": "error with type",
//   "when": "2024-06-01T12:00:00Z",
//   "request_id": "request-1234",
//   "tags": {
//     "key1": "value1",
//     "key2": 42,
//     "key3": 3.14,
//     "key4": true
//   }
// }
```

#### Log as plain string

```go
fmt.PrintF("%+v\n", err)
// Output:
// [Type:CustomType1] [Message:error with type] [When:2024-06-01 12:00:00 +0000 UTC] [RequestID:request-1234] [Tags:key1=value1,key2=42,key3=3.14,key4=true]
```

#### Is() and As()

