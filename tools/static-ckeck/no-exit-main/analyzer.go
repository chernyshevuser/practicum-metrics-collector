package noexitmain

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "noexitmain",
	Doc:  "check if os.Exit() is exist in package 'main' in func 'main'",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		foundPos := DetectExitMain(file)
		if foundPos != nil {
			pos := *foundPos
			pass.Reportf(pos, "os.Exit() call in main func")
		}
	}
	return nil, nil
}
