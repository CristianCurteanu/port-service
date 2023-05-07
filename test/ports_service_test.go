package test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/CristianCurteanu/koken-api/internal/domains/ports"
	httpApi "github.com/CristianCurteanu/koken-api/internal/infra/http"
	"github.com/CristianCurteanu/koken-api/internal/infra/storage/inmemory"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPortsFileUpload(t *testing.T) {
	t.Run("test success upload", func(t *testing.T) {
		router := httpApi.NewRouter(
			httpApi.PortHandlers(ports.NewPortService(ports.NewPortRepositories(inmemory.NewInMemoryStorage()))),
		)
		resp := httptest.NewRecorder()
		req, err := formFileUpload("/ports", "ports", "./fixtures/success.json")
		require.NoError(t, err)
		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusCreated, resp.Code)

		resp = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodGet, "/ports/AEJEA", nil)
		require.NoError(t, err)
		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusOK, resp.Code)

		body := resp.Body.String()
		require.Contains(t, body, "port_code\":\"AEJEA")
		require.Contains(t, body, "code\":\"52051")
	})

	t.Run("fail if data is passed as an array", func(t *testing.T) {
		router := httpApi.NewRouter(
			httpApi.PortHandlers(ports.NewPortService(ports.NewPortRepositories(inmemory.NewInMemoryStorage()))),
		)
		resp := httptest.NewRecorder()
		req, err := formFileUpload("/ports", "ports", "./fixtures/fail_as_array.json")
		require.NoError(t, err)
		router.ServeHTTP(resp, req)

		body := resp.Body.String()
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, body, "bad_json_file")
	})

	t.Run("fail if data in json is not json", func(t *testing.T) {
		router := httpApi.NewRouter(
			httpApi.PortHandlers(ports.NewPortService(ports.NewPortRepositories(inmemory.NewInMemoryStorage()))),
		)
		resp := httptest.NewRecorder()
		req, err := formFileUpload("/ports", "ports", "./fixtures/fail_as_array.json")
		require.NoError(t, err)
		router.ServeHTTP(resp, req)

		body := resp.Body.String()
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, body, "bad_json_file")
	})

	t.Run("fail if struct of value is not as expected", func(t *testing.T) {
		router := httpApi.NewRouter(
			httpApi.PortHandlers(ports.NewPortService(ports.NewPortRepositories(inmemory.NewInMemoryStorage()))),
		)
		resp := httptest.NewRecorder()
		req, err := formFileUpload("/ports", "ports", "./fixtures/fail_as_array.json")
		require.NoError(t, err)
		router.ServeHTTP(resp, req)

		body := resp.Body.String()
		require.Equal(t, http.StatusBadRequest, resp.Code)
		require.Contains(t, body, "bad_json_file")
	})

	t.Run("fail if service is not storing data correctly", func(t *testing.T) {
		serviceMock := new(MockPortsService)
		serviceMock.On("CreateOrUpdateMany", mock.Anything, mock.Anything).Return(errors.New("store failed error"))
		router := httpApi.NewRouter(
			httpApi.PortHandlers(serviceMock),
		)
		resp := httptest.NewRecorder()
		req, err := formFileUpload("/ports", "ports", "./fixtures/success.json")
		require.NoError(t, err)
		router.ServeHTTP(resp, req)

		require.Equal(t, http.StatusInternalServerError, resp.Code)

		body := resp.Body.String()
		require.Contains(t, body, "err_data_store")
	})
}

type MockPortsService struct {
	mock.Mock
}

func (m *MockPortsService) GetByPortCode(ctx context.Context, code string) (ports.Port, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(ports.Port), args.Error(1)
}

func (m *MockPortsService) CreateOrUpdate(ctx context.Context, port ports.Port) error {
	args := m.Called(ctx, port)
	return args.Error(0)
}

func (m *MockPortsService) CreateOrUpdateMany(ctx context.Context, ports []ports.Port) error {
	args := m.Called(ctx, ports)
	return args.Error(0)
}

func formFileUpload(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
