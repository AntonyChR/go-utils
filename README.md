# go-utils

A collection of simple and generic utility functions for Go, organized into logical packages.

## Description

This project provides a set of reusable Go packages, each focusing on a specific domain of functionality. The goal is to offer clean, tested, and idiomatic Go code to speed up development by handling common tasks.

## Installation

To use these utilities in your project, you can fetch them with `go get`:

```bash
go get github.com/AntonyChR/go-utils
```

## Available Packages

Below is a list of the available packages and a brief description of their purpose.

-   **`array`**: Contains functions for common slice operations like `Map`, `Filter`, `Flatten`, and `Contains`.
    ```go
    import "github.com/AntonyChR/go-utils/array"
    ```
-   **`assert`**: Provides helper functions for writing tests, such as `Assert`, `AssertEq`, and `AssertNe`.
    ```go
    import "github.com/AntonyChR/go-utils/assert"
    ```
-   **`bytes`**: Includes utilities for byte manipulation.
    ```go
    import "github.com/AntonyChR/go-utils/bytes"
    ```
-   **`math`**: Offers mathematical helper functions like `Round`.
    ```go
    import "github.com/AntonyChR/go-utils/math"
    ```
-   **`network`**: A package for network-related tasks, like retrieving the local IPv4 address.
    ```go
    import "github.com/AntonyChR/go-utils/network"
    ```
-   **`terminal`**: Provides functions to interact with the terminal, such as getting the terminal size.
    ```go
    import "github.com/AntonyChR/go-utils/terminal"
    ```
-   **`vec`**: Implements a `Vec` type for 2D and 3D vector mathematics, including operations like dot product, cross product, and rotation.
    ```go
    import "github.com/AntonyChR/go-utils/vec"
    ```

## Usage Example

Here is a simple example of how to use the `array` package to double the numbers in a slice:

```go
package main

import (
	"fmt"
	"github.com/AntonyChR/go-utils/array"
)

func main() {
	nums := []int{1, 2, 3, 4, 5}

	// Use the Map function from the array package
	doubled := array.Map(nums, func(n int) int {
		return n * 2
	})

	fmt.Println(doubled) // Output: [2 4 6 8 10]
}
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the MIT License.
