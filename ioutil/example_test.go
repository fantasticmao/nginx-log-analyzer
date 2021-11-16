package ioutil

import "os"

func ExampleFatal() {
	fatal(os.Stdout, func(i int) {}, "log: %v\n", "Hello, World!")
	// Output:
	// log: Hello, World!
}
