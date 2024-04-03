package server

import (
	devicectrl "github.com/arfan21/mertani/internal/device/controller"
	devicerepo "github.com/arfan21/mertani/internal/device/repository"
	devicesvc "github.com/arfan21/mertani/internal/device/service"
	sensorctrl "github.com/arfan21/mertani/internal/sensor/controller"
	sensorrepo "github.com/arfan21/mertani/internal/sensor/repository"
	sensorsvc "github.com/arfan21/mertani/internal/sensor/service"
	userctrl "github.com/arfan21/mertani/internal/user/controller"
	userrepo "github.com/arfan21/mertani/internal/user/repository"
	usersvc "github.com/arfan21/mertani/internal/user/service"
	"github.com/arfan21/mertani/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Routes() {

	api := s.app.Group("/api")
	api.Get("/health-check", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	userRepo := userrepo.New(s.db)
	userSvc := usersvc.New(userRepo)
	userCtrl := userctrl.New(userSvc)

	deviceRepo := devicerepo.New(s.db)
	deviceSvc := devicesvc.New(deviceRepo)
	deviceCtrl := devicectrl.New(deviceSvc)

	sensorRepo := sensorrepo.New(s.db)
	sensorSvc := sensorsvc.New(sensorRepo)
	sensorCtrl := sensorctrl.New(sensorSvc)

	s.RoutesCustomer(api, userCtrl)
	s.RoutesDevice(api, deviceCtrl)
	s.RoutesSensor(api, sensorCtrl)
}

func (s Server) RoutesCustomer(route fiber.Router, ctrl *userctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	usersV1 := v1.Group("/users")
	usersV1.Post("/register", ctrl.Register)
	usersV1.Post("/login", ctrl.Login)
}

func (s Server) RoutesDevice(route fiber.Router, ctrl *devicectrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	deviceV1 := v1.Group("/devices", middleware.JWTAuth)
	deviceV1.Post("", ctrl.Create)
	deviceV1.Get("/:id", ctrl.GetByID)
	deviceV1.Get("", ctrl.GetAll)
	deviceV1.Put("/:id", ctrl.Update)
	deviceV1.Delete("/:id", ctrl.Delete)
}

func (s Server) RoutesSensor(route fiber.Router, ctrl *sensorctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	sensorV1 := v1.Group("/sensors", middleware.JWTAuth)
	sensorV1.Post("", ctrl.Create)
	sensorV1.Get("/:id", ctrl.GetByID)
	sensorV1.Put("/:id", ctrl.Update)
	sensorV1.Delete("/:id", ctrl.Delete)
}
