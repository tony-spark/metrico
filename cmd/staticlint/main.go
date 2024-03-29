// staticlint is a linter used in this project
//
// staticlint performs following checks:
//
// 1. From standard library
//   - assign
//   - bools
//   - composite
//   - copylock
//   - deepequalerrors
//   - errorsas
//   - httpresponse
//   - ifaceassert
//   - loopclosure
//   - lostcancel
//   - nilfunc
//   - nilness
//   - printf
//   - shadow
//   - sigchanyzer
//   - stdmethods
//   - stringintconv
//   - structtag
//   - tests
//   - unmarshal
//   - unreachable
//   - unusedresult
//   - unusedwrite
//
// 2. From staticcheck - SA (common) and S1 (code simplification) checks
//
// 3. errcheck - to check for unchecked errors
//
// 4. wrapcheck - to check that errors from external packages are wrapped during return to help identify the error source during debugging
//
// 5. os.Exit check - to check there's no os.Exit() call in the main function
package main

import (
	"strings"

	"github.com/kisielk/errcheck/errcheck"
	"github.com/tomarrell/wrapcheck/v2/wrapcheck"
	"github.com/tony-spark/metrico/internal/staticlint/osexit"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var checks = []*analysis.Analyzer{
		assign.Analyzer,
		bools.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		printf.Analyzer,
		shadow.Analyzer,
		sigchanyzer.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
	}

	for _, a := range staticcheck.Analyzers {
		if strings.HasPrefix(a.Analyzer.Name, "SA") || strings.HasPrefix(a.Analyzer.Name, "S1") {
			checks = append(checks, a.Analyzer)
		}
	}

	errcheck.DefaultExcludedSymbols = append(errcheck.DefaultExcludedSymbols,
		"(io.ReadCloser).Close",
		"(*database/sql.Tx).Rollback",
		"(*database/sql.Stmt).Close",
		"(*database/sql.Rows).Close",
	)

	checks = append(checks, errcheck.Analyzer)

	wConfig := wrapcheck.NewDefaultConfig()
	wConfig.IgnoreSigRegexps = append(wConfig.IgnoreSigRegexps,
		`.*github.com/tony-spark/metrico/internal/.*`, // ignore error wrapping in internal packages
		`.*github.com/tony-spark/metrico/gen/.*`,      // ignore error wrapping in generated code
		`.*google.golang.org/grpc.*`,
	)
	wConfig.IgnoreSigs = append(wConfig.IgnoreSigs,
		"func github.com/hashicorp/go-multierror.Append(err error, errs ...error) *github.com/hashicorp/go-multierror.Error",
		"func (*golang.org/x/sync/errgroup.Group).Wait() error")
	checks = append(checks, wrapcheck.NewAnalyzer(wConfig))

	checks = append(checks, osexit.Analyzer)

	multichecker.Main(
		checks...,
	)
}
