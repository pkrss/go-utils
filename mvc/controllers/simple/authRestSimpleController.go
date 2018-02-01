package simple

import (
	"github.com/pkrss/go-utils/orm"
)

type SimpleAuthRestController struct {
	SimpleAuthController
	Model  orm.BaseModelInterface
	Helper SimpleAuthRestHelper
}

func (this *SimpleAuthRestController) OnPrepare() {
	this.Helper = CreateSimpleAuthRestHelper(this, this.Model)
}

func (this *SimpleAuthRestController) OnLeave() {
	this.Helper = nil
	this.Model = nil
}

func (this *SimpleAuthRestController) CloneAttribute(src ControllerInterface) {
	this.SimpleAuthController.CloneAttribute(src)
	if src == nil {
		return
	}
	s := src.(SimpleAuthRestController)
	if s != nil {
		this.Model = s.Model
		this.Helper = s.Helper
	}
}
