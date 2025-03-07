package userctrl

import (
	"github.com/arfan21/mertani/internal/model"
	"github.com/arfan21/mertani/internal/user"
	"github.com/arfan21/mertani/pkg/exception"
	"github.com/arfan21/mertani/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type ControllerHTTP struct {
	svc user.Service
}

func New(svc user.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Register user
// @Description Register user
// @Tags user
// @Accept json
// @Produce json
// @Param body body model.UserRegisterRequest true "Payload user Register Request"
// @Success 201 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/users/register [post]
func (ctrl ControllerHTTP) Register(c *fiber.Ctx) error {
	var req model.UserRegisterRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	err = ctrl.svc.Register(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusCreated).JSON(pkgutil.HTTPResponse{
		Code: fiber.StatusCreated,
	})
}

// @Summary Login user
// @Description Login user
// @Tags user
// @Accept json
// @Produce json
// @Param body body model.UserLoginRequest true "Payload user Login Request"
// @Success 200 {object} pkgutil.HTTPResponse{data=model.UserLoginResponse}
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/users/login [post]
func (ctrl ControllerHTTP) Login(c *fiber.Ctx) error {
	var req model.UserLoginRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	res, err := ctrl.svc.Login(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Code: fiber.StatusOK,
		Data: res,
	})
}
