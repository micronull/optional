# Optional Package

The `optional` package in Go provides a generic type that can represent values which may or may not be set, 
including the concept of null. This allows handling scenarios where a value might be missing or explicitly set to null in JSON. 
The package is designed to be flexible and efficient, with support for custom marshalling and unmarshalling functions.

## Features

- **Generic Type**: Supports any type `T`.
- **Null Handling**: Distinguishes between unset values, null values, and non-null values.
- **Custom Marshalling/Unmarshalling**: Allows changing the JSON marshalling/unmarshalling implementation, such as using a faster library like [json-iterator](https://pkg.go.dev/github.com/json-iterator/go).

## Installation

To install the `optional` package, use the following command:

```bash
go get github.com/micronull/optional
```

## Usage

### Basic Usage

Here's a simple example demonstrating how to use the `Type` struct and its methods:

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/micronull/optional"
)

func main() {
	// Create a new optional value that is not null
	optVal := optional.New[string]("hello", false)
	fmt.Println(optVal.IsSet())    // Output: true
	fmt.Println(optVal.IsSetNull()) // Output: false

	// Marshal the optional value to JSON
	jsonBytes, _ := json.Marshal(optVal)
	fmt.Println(string(jsonBytes)) // Output: "hello"

	// Create a new optional value that is explicitly null
	optNull := optional.New[string]("", true)
	fmt.Println(optNull.IsSet())    // Output: true
	fmt.Println(optNull.IsSetNull()) // Output: true

	// Marshal the null optional value to JSON
	jsonBytes, _ = json.Marshal(optNull)
	fmt.Println(string(jsonBytes)) // Output: null

	// Unmarshal JSON into an optional value
	var opt Type[string]
	json.Unmarshal([]byte(`"world"`), &opt)
	fmt.Println(opt.V) // Output: world

	json.Unmarshal([]byte(`null`), &opt)
	fmt.Println(opt.IsSetNull()) // Output: true
}
```

### Custom Marshalling/Unmarshalling

You can replace the default JSON marshalling and unmarshalling functions with your own implementations, such as using [json-iterator](https://pkg.go.dev/github.com/json-iterator/go):

```go
package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"

	"github.com/micronull/optional"
)

func init() {
	// Replace the default marshaller with json-iterator's marshaller
	optional.ChangeMarshal(jsoniter.Marshal)
	// Replace the default unmarshaller with json-iterator's unmarshaller
	optional.ChangeUnmarshal(jsoniter.Unmarshal)
}

func main() {
	optVal := optional.New[string]("hello", false)

	jsonBytes, _ := jsoniter.Marshal(optVal)
	fmt.Println(string(jsonBytes)) // Output: "hello"
}
```

## Contributing

Contributions are welcome! If you have any suggestions or find a bug, please open an issue on the [GitHub repository](https://github.com/micronull/optional).

## License

This package is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
