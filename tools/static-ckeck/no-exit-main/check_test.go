package noexitmain

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func isDetected(src string) bool {
	file, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		panic(err)
	}

	return DetectExitMain(file) != nil
}

func TestFindExitMain_Detected(t *testing.T) {
	{
		code := `
	package main
	import "fmt"

	func main() {
		os.Exit(0)
	}
	`
		require.True(t, isDetected(code))
	}
	{
		code := `
		package main
		import "fmt"
	
		func tmp() {
			os.Exit(0)
		}
		
		func main() {
			os.Exit(0)
		}
		`
		require.True(t, isDetected(code))
	}

}

func TestFindExitMain_NotDetected(t *testing.T) {
	{
		code := `
	package main
	import "fmt"

	func tmp() {
		os.Exit(0)
	}
	`
		require.False(t, isDetected(code))
	}
	{
		code := `
	package main
	import "fmt"

	func tmp() {
		// os.Exit
	}
	`
		require.False(t, isDetected(code))
	}
	{
		code := `
	package main
	import "fmt"

	var val = "os.Exit()"

	func tmp() {
		fmt.Println(val)
	}
	`
		require.False(t, isDetected(code))
	}
	{
		code := `
	package main
	import "fmt"

	func main() {
		fmt.Println("")
	}
	`
		require.False(t, isDetected(code))
	}
	{
		code := `
	package notmain
	import "fmt"

	func main() {
		os.Exit(-1)
	}
	`
		require.False(t, isDetected(code))
	}
}
