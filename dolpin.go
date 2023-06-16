package dolpin

type dolpin struct {
	config         config
	mockFunction   mockFunction
	funcDecoration string
}

type mockFunction struct {
	funcName string
	args     []any
	results  []any
}

type ()

func Args(args []any) []any {
	return args
}

func Results(results []any) []any {
	return results
}

func New(conf Config) *dolpin {
	return &dolpin{config: config{mockPath: conf.MockPath}}
}

func (_dolpin *dolpin) CallFunction(funcName string, args []any, results []any) {
	mockFunc := mockFunction{
		funcName: funcName,
		args:     args,
		results:  results,
	}

	_dolpin.mockFunction = mockFunc
}

func (_dolpin *dolpin) SetDecoration(decoration string) {
	_dolpin.funcDecoration = decoration
}
