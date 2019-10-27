/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 16-12-2018
 * |
 * | File Name:     instance.go
 * +===============================================
 */

package actions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/toskatok/lg/models"
)

// instance creation request
type instanceReq struct {
	Name string
	models.Config
}

// InstancesHandler manages instances of load generators.
type InstancesHandler struct {
	// list of the running instances
	instances map[string]*models.Instance
}

// NewInstancesHandler creates new instance of load generator instances handler
func NewInstancesHandler() *InstancesHandler {
	return &InstancesHandler{
		instances: make(map[string]*models.Instance),
	}
}

// List returns all running instances This function is mapped
// to the path GET /instances
func (v *InstancesHandler) List(c echo.Context) error {
	return c.JSON(http.StatusOK, v.instances)
}

// Create runs new generator instance. This function is mapped
// to the path POST /instances
func (v *InstancesHandler) Create(c echo.Context) error {
	var req instanceReq

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var config models.Config = req.Config

	// check for duplicate name
	if _, ok := v.instances[req.Name]; ok {
		return echo.NewHTTPError(http.StatusBadRequest, "duplicate name")
	}

	rate, err := time.ParseDuration(c.QueryParam("rate"))
	if err != nil {
		rate = 1 * time.Millisecond
	}

	destination := c.QueryParam("destination")
	if destination == "" {
		destination = "mqtt://127.0.0.1:1883"
	}

	// create and run newly created instance
	i, err := models.NewInstance(config, rate, destination)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	i.Run()

	v.instances[req.Name] = i

	return c.JSON(http.StatusOK, true)
}

// Show shows the detail of given instance. This function is mapped
// to the path GET /instances/{instance_id}
func (v *InstancesHandler) Show(c echo.Context) error {
	id := c.Param("instance_id")

	i, ok := v.instances[id]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("there is no instnace with name %s", id))
	}

	return c.JSON(http.StatusOK, i.R.Count())
}

// Destroy stops given instance and removes it from the instances list.
// This function is mapped to the path DELETE /instances/{instance_id}
func (v *InstancesHandler) Destroy(c echo.Context) error {
	id := c.Param("instance_id")

	i, ok := v.instances[id]
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("there is no instnace with name %s", id))
	}

	i.Stop()
	delete(v.instances, id)

	return c.JSON(http.StatusOK, true)
}
