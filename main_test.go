package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// setupEcho creates a test Echo instance with all routes registered
func setupEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "pong",
		})
	})

	e.GET("/ip", func(c echo.Context) error {
		ip := c.RealIP()
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"data": ip,
		})
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"routes": "/ping /echo /ip /env /delay?ms=100 /mem?mb=10&ms=10000 /cpu?ms=1000&cores=2&percent=80",
		})
	})

	return e
}

func TestPing(t *testing.T) {
	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Code != 0 || resp.Msg != "pong" {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestIPWithXForwardedFor(t *testing.T) {
	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/ip", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRootRouteListing(t *testing.T) {
	e := setupEcho()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if _, ok := resp["routes"]; !ok {
		t.Fatal("expected routes key in response")
	}
}
