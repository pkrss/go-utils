package simple

import (
	"github.com/pkrss/go-utils/pqsql"
)

type SimpleAuthRestController struct {
	SimpleAuthController
	Model  pqsql.BaseModelInterface
	Helper SimpleAuthRestHelper
}

func (this *SimpleAuthRestController) OnPrepare() {
	this.Helper = CreateSimpleAuthRestHelper(this, this.Model)
}

func (this *SimpleAuthRestController) OnLeave() {
	this.Helper = nil
	this.Model = nil
}
