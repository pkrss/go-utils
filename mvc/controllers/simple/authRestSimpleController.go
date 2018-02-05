package simple

import (
	"github.com/pkrss/go-utils/mvc/controllers"
	"github.com/pkrss/go-utils/orm"
)

type SimpleAuthRestController struct {
	controllers.AuthController
	Model  orm.BaseModelInterface
	Helper *SimpleAuthRestHelper
}

func (this *SimpleAuthRestController) OnPrepare() {
	this.Helper = CreateSimpleAuthRestHelper(this, this.Model)
}

func (this *SimpleAuthRestController) OnLeave() {
	this.Helper = nil
	this.Model = nil
}
