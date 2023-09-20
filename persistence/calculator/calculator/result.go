package calculator

import "github.com/ytake/protoactor-go-example/persistence/calculator/protobuf"

// CalculationResult is a struct that holds the result of a calculation.
type CalculationResult struct {
	protobuf.CalculationResult
}

func (c *CalculationResult) Reset() *CalculationResult {
	return &CalculationResult{
		CalculationResult: protobuf.CalculationResult{},
	}
}

// Add adds a value to the result.
func (c *CalculationResult) Add(value float64) *CalculationResult {
	return &CalculationResult{
		CalculationResult: protobuf.CalculationResult{
			Result: c.Result + value,
		},
	}
}

// Subtract subtracts a value from the result.
func (c *CalculationResult) Subtract(value float64) *CalculationResult {
	return &CalculationResult{
		CalculationResult: protobuf.CalculationResult{
			Result: c.Result - value,
		},
	}
}

// Divide divides the result by a value.
func (c *CalculationResult) Divide(value float64) *CalculationResult {
	return &CalculationResult{
		CalculationResult: protobuf.CalculationResult{
			Result: c.Result / value,
		},
	}
}

// Multiply multiplies the result by a value.
func (c *CalculationResult) Multiply(value float64) *CalculationResult {
	return &CalculationResult{
		CalculationResult: protobuf.CalculationResult{
			Result: c.Result * value,
		},
	}
}
