package cmd

import (
	"fmt"
	"strconv"
)

type ArrayFlags []float64

func (i *ArrayFlags) String() string {
	return fmt.Sprintf("%g", *i)
}

func (i *ArrayFlags) Set(value string) error {
	tmp, err := strconv.ParseFloat(value, 64)
	if err != nil {
		*i = append(*i, -1)
	} else {
		*i = append(*i, tmp)
	}
	return nil
}