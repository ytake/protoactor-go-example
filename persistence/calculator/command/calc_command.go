package command

type Clear struct{}
type Add struct {
	Value float64
}

type Subtract struct {
	Value float64
}

type Divide struct {
	Value float64
}

type Multiply struct {
	Value float64
}

type PrintResult struct{}
type GetResult struct{}
