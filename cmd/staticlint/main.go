/*
Package main provides a static analyzer for the project, using the multichecker package from `golang.org/x/tools/go/analysis`.
This analyzer gathers various code analyzers and runs them concurrently to check Go projects.

Analyzers used:
1. `assign.Analyzer` - checks for correct assignments.
2. `findcall.Analyzer` - finds function calls based on specified patterns.
3. `inspect.Analyzer` - provides a basic mechanism to walk through the code tree (AST).
4. `printf.Analyzer` - checks the correctness of formatting function calls (e.g., `fmt.Printf`).
5. `shadow.Analyzer` - detects cases of variable shadowing.
6. `shift.Analyzer` - checks for correct usage of shift operations.
7. `structtag.Analyzer` - verifies the correctness of struct tags.
8. `staticcheck` - a library of static code analyzers, including checks like `SA` and `ST` (finding potential issues).
9. `quickfix.Analyzers` - a set of quick-fix analyzers, including `QF1001` (suggests possible improvements).

Mechanism of work:
The program collects a set of analyzers that will be executed on the project code. The `multichecker.Main()` method is used to run these analyzers in sequence, analyzing the code and outputting the results.
*/
package main

import (
	"strings"

	noexitmain "github.com/chernyshevuser/practicum-metrics-collector/tools/static-ckeck/no-exit-main"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/staticcheck"
)

// main function gathers all analyzers and runs them using multichecker.
func main() {
	// Default list of analyzers used to check the code.
	analyzers := []*analysis.Analyzer{
		assign.Analyzer,    // Checks for correct assignments
		findcall.Analyzer,  // Finds function calls by pattern
		inspect.Analyzer,   // Provides a mechanism for AST traversal (code tree)
		printf.Analyzer,    // Verifies the correct use of printf-style functions
		shadow.Analyzer,    // Detects shadowed variables
		shift.Analyzer,     // Checks the correctness of shift operations
		structtag.Analyzer, // Verifies the correctness of struct tags
	}

	// Adding static checks from the staticcheck package that start with "SA"
	checks := staticcheck.Analyzers
	for _, v := range checks {
		if strings.HasPrefix(strings.ToUpper(v.Analyzer.Name), "SA") {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	// Adding static checks from staticcheck with prefix "ST"
	for _, v := range checks {
		if strings.HasPrefix(strings.ToUpper(v.Analyzer.Name[:2]), "ST") {
			analyzers = append(analyzers, v.Analyzer)
			break
		}
	}

	// Adding quickfix analyzers, including "QF1001"
	for _, v := range quickfix.Analyzers {
		if strings.ToUpper(v.Analyzer.Name) == "QF1001" {
			analyzers = append(analyzers, v.Analyzer)
			break
		}
	}

	// Adding a custom analyzer that checks for the presence of os.Exit in the main function.
	analyzers = append(analyzers, noexitmain.Analyzer)

	// Running the multichecker with the collected list of analyzers.
	multichecker.Main(analyzers...)
}
