package devicectrl

import (
	"github.com/arfan21/mertani/internal/device"
	"github.com/arfan21/mertani/internal/model"
	"github.com/arfan21/mertani/pkg/exception"
	"github.com/arfan21/mertani/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type ControllerHTTP struct {
	svc device.Service
}

func New(svc device.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Create Device
// @Description Create a new device
// @Tags Device
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer token"
// @Param device body model.DeviceCreateRequest true "Device data"
// @Success 201 {object} pkgutil.HTTPResponse
// @Router /api/v1/devices [post]
func (c ControllerHTTP) Create(ctx *fiber.Ctx) error {
	var req model.DeviceCreateRequest
	err := ctx.BodyParser(&req)
	exception.PanicIfNeeded(err)

	err = c.svc.Create(ctx.UserContext(), req)
	exception.PanicIfNeeded(err)

	return ctx.Status(fiber.StatusCreated).JSON(pkgutil.HTTPResponse{
		Code:    fiber.StatusCreated,
		Message: "Device created successfully",
	})
}

// @Summary Get Device By ID
// @Description Get device by ID
// @Tags Device
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer token"
// @Param id path string true "Device ID"
// @Success 200 {object} pkgutil.HTTPResponse{data=model.DeviceResponse}
// @Router /api/v1/devices/{id} [get]
func (c ControllerHTTP) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	res, err := c.svc.GetByID(ctx.UserContext(), id)
	exception.PanicIfNeeded(err)

	return ctx.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Code:   fiber.StatusOK,
		Status: "success",
		Data:   res,
	})
}

// @Summary Get All Devices
// @Description Get all devices
// @Tags Device
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer token"
// @Success 200 {object} pkgutil.HTTPResponse{data=[]model.DeviceResponse}
// @Router /api/v1/devices [get]
func (c ControllerHTTP) GetAll(ctx *fiber.Ctx) error {
	res, err := c.svc.GetAll(ctx.UserContext())
	exception.PanicIfNeeded(err)

	return ctx.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Code:   fiber.StatusOK,
		Status: "success",
		Data:   res,
	})
}

// @Summary Update Device
// @Description Update device by ID
// @Tags Device
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer token"
// @Param id path string true "Device ID"
// @Param device body model.DeviceUpdateRequest true "Device data"
// @Success 200 {object} pkgutil.HTTPResponse
// @Router /api/v1/devices/{id} [put]
func (c ControllerHTTP) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var req model.DeviceUpdateRequest
	err := ctx.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.ID = id

	err = c.svc.Update(ctx.UserContext(), req)
	exception.PanicIfNeeded(err)

	return ctx.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Code:    fiber.StatusOK,
		Message: "Device updated successfully",
	})
}

// @Summary Delete Device
// @Description Delete device by ID
// @Tags Device
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer token"
// @Param id path string true "Device ID"
// @Success 200 {object} pkgutil.HTTPResponse
// @Router /api/v1/devices/{id} [delete]
func (c ControllerHTTP) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.svc.Delete(ctx.UserContext(), id)
	exception.PanicIfNeeded(err)

	return ctx.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Code:    fiber.StatusOK,
		Message: "Device deleted successfully",
	})
}
