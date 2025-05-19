// @@docs Main Package
// This is the main package of the Go application.
// It contains the entry point.
package main

import "fmt"

func greet(name string) string {
	// @@docs Greeting
	// Functions related to greeting users.
	return fmt.Sprintf("Hello, %s!", name)
}

// @@docs Configuration
// App configuration helpers.
var config = struct {
	Env string
}{
	Env: "development",
}

// @@/code

func main() {
	/* @@docs Main Logic
	The actual application logic.
	This section has reduced indentation. */
	fmt.Println(greet("World"))
}
