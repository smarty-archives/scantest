package core

import (
	"path/filepath"

	"github.com/smartystreets/scantest/scantest/contract"
)

type Checksummer struct {
	previous int64
}

func NewChecksummer() *Checksummer {
	return &Checksummer{}
}

func (this *Checksummer) Handle(context *contract.Context) {
	var sum int64
	for _, file := range context.Files {
		if filepath.Base(file.Path) == "generated_by_gunit_test.go" {
			continue
		}
		if file.IsGoFile || filepath.Base(file.Path) == ".testargs" {
			sum += file.Size + file.Modified
		}
	}
	if sum == this.previous {
		context.Error = contract.ContextComplete
	}
	this.previous = sum
}
