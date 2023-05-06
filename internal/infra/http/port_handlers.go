package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/CristianCurteanu/koken-api/errors"
	"github.com/CristianCurteanu/koken-api/internal/domains/ports"
	"github.com/gin-gonic/gin"
)

func PortHandlers(service ports.PortService) DomainHandler {
	return DomainHandler{
		Path: "/",
		Middlewares: []gin.HandlerFunc{
			requestLogMiddleware,
		},
		Routes: []Route{
			{
				Path:   "/ports",
				Method: http.MethodPost,
				Middlewares: []gin.HandlerFunc{
					requestLogMiddleware,
				},
				Handler: createPortsHandler(service),
			},
			{
				Path:    "/ports/:port_code",
				Method:  http.MethodGet,
				Handler: getPortByPortCodeHandler(service),
			},
		},
	}
}

func requestLogMiddleware(c *gin.Context) {
	log.Printf("PORTS[%s][received]: %q", c.Request.Method, c.Request.URL.Path)
}

func getPortByPortCodeHandler(service ports.PortService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		portCode, found := ctx.Params.Get("port_code")
		if !found || portCode == "" {
			ctx.SecureJSON(http.StatusBadRequest, errors.ApiError{
				Code:    "no_port_code",
				Message: "No port_code provided, please check the URL, to make sure you added port_code path param",
			})
			return
		}

		port, err := service.GetByPortCode(ctx, portCode)
		if err != nil {
			ctx.SecureJSON(http.StatusNotFound, errors.ApiError{
				Code:    "not_found",
				Message: "No port found with the specified port code",
			})
			return
		}

		ctx.SecureJSON(http.StatusOK, portResponse(port))
	}
}

func createPortsHandler(service ports.PortService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, header, _ := ctx.Request.FormFile("ports")

		// file_buf.copy
		buf := bytes.NewBuffer(nil)
		copied, err := io.Copy(buf, file)
		if err != nil {
			log.Printf("PORTS[CREATE][file_buf.copy], error=%q\n", err)
			ctx.SecureJSON(http.StatusInternalServerError, errors.ApiError{
				Code:    "internal_error",
				Message: "Please check with the administrator",
			})
			return
		}

		if copied != header.Size {
			log.Printf("PORTS[CREATE][file_buf.copy], error=%q\n", fmt.Errorf("copied `%d`, while header-size is: `%d`", copied, header.Size))
			ctx.SecureJSON(http.StatusInternalServerError, errors.ApiError{
				Code:    "internal_error",
				Message: "Please check with the administrator",
			})
			return
		}

		// file_buf.decode
		ports, err := decodePortsBody(buf.Bytes())
		if err != nil {
			log.Printf("PORTS[CREATE][file_buf.unmarshal], error=%q\n", err)

			ctx.SecureJSON(http.StatusBadRequest, errors.ApiError{
				Code:    "bad_json_file",
				Message: "Please check your json file, there might be syntax issues",
			})
			return
		}

		// service.create_or_update_many
		err = service.CreateOrUpdateMany(ctx, ports)
		if err != nil {
			log.Printf("PORTS[CREATE][service.create_or_update_many], error=%q\n", err)
			ctx.SecureJSON(http.StatusInternalServerError, errors.ApiError{
				Code:    "invalid_data",
				Message: "Error while storing the data; please contact administrator to check the reason of failure",
			})
			return
		}

		ctx.SecureJSON(http.StatusCreated, "{}")
	}
}

type portRequest struct {
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Code        string    `json:"code"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
}

func decodePortsBody(data []byte) (res []ports.Port, err error) {
	var portsReq map[string]portRequest
	err = json.Unmarshal(data, &portsReq)
	if err != nil {
		return res, err
	}

	for code, portBody := range portsReq {
		res = append(res, ports.Port{
			PortCode:    code,
			Name:        portBody.Name,
			City:        portBody.City,
			Country:     portBody.Country,
			Code:        portBody.Code,
			Alias:       portBody.Alias,
			Regions:     portBody.Regions,
			Coordinates: portBody.Coordinates,
			Province:    portBody.Province,
			Timezone:    portBody.Timezone,
			Unlocs:      portBody.Unlocs,
		})
	}

	return res, nil
}

type portResponse struct {
	PortCode    string    `json:"port_code,omitempty"`
	Name        string    `json:"name,omitempty"`
	City        string    `json:"city,omitempty"`
	Country     string    `json:"country,omitempty"`
	Code        string    `json:"code,omitempty"`
	Alias       []string  `json:"alias,omitempty"`
	Regions     []string  `json:"regions,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty"`
	Province    string    `json:"province,omitempty"`
	Timezone    string    `json:"timezone,omitempty"`
	Unlocs      []string  `json:"unlocs,omitempty"`
}
