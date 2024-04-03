package sensorctrl

import (
	"github.com/arfan21/mertani/internal/model"
	"github.com/arfan21/mertani/internal/sensor"
	"github.com/arfan21/mertani/pkg/exception"
	"github.com/arfan21/mertani/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type ControllerHTTP struct {
	svc sensor.Service
}

func New(svc sensor.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Create Sensor
// @Description Create Sensor
// @Tags Sensor
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer token"
// @Param body body model.SensorCreateRequest true "Payload Sensor Create Request"
// @Success 201 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/sensors [post]
func (c ControllerHTTP) Create(ctx *fiber.Ctx) error {
	var req model.SensorCreateRequest
	err := ctx.BodyParser(&req)
	exception.PanicIfNeeded(err)

	err = c.svc.Create(ctx.UserContext(), req)
	exception.PanicIfNeeded(err)

	return ctx.Status(fiber.StatusCreated).JSON(pkgutil.HTTPResponse{
		Code:    fiber.StatusCreated,
		Message: "Sensor created successfully",
	})
}

// @Summary Get Sensor By ID
// @Description Get Sensor By ID
// @Tags Sensor
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer token"
// @Param id path string true "Sensor ID"
// @Success 200 {object} pkgutil.HTTPResponse{data=model.SensorResponse} "Success"
// @Failure 404 {object} pkgutil.HTTPResponse
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/sensors/{id} [get]
func (c ControllerHTTP) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	res, err := c.svc.GetByID(ctx.UserContext(), id)
	exception.PanicIfNeeded(err)

	return ctx.JSON(pkgutil.HTTPResponse{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    res,
	})
}

// @Summary Update Sensor
// @Description Update Sensor
// @Tags Sensor
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer token"
// @Param id path string true "Sensor ID"
// @Param body body model.SensorUpdateRequest true "Payload Sensor Update Request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 404 {object} pkgutil.HTTPResponse
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/sensors/{id} [put]
func (c ControllerHTTP) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	var req model.SensorUpdateRequest
	err := ctx.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.ID = id
	err = c.svc.Update(ctx.UserContext(), req)
	exception.PanicIfNeeded(err)

	return ctx.JSON(pkgutil.HTTPResponse{
		Code:    fiber.StatusOK,
		Message: "Sensor updated successfully",
	})
}

// @Summary Delete Sensor
// @Description Delete Sensor
// @Tags Sensor
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer token"
// @Param id path string true "Sensor ID"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 404 {object} pkgutil.HTTPResponse
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /api/v1/sensors/{id} [delete]
func (c ControllerHTTP) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.svc.Delete(ctx.UserContext(), id)
	exception.PanicIfNeeded(err)

	return ctx.JSON(pkgutil.HTTPResponse{
		Code:    fiber.StatusOK,
		Message: "Sensor deleted successfully",
	})
}
