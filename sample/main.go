package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	cache "github.com/takuoki/gocodecache"
)

const (
	dbPingRetryLimit    = 30
	dbPingRetryInterval = 1 * time.Second
	reloadInterval      = 1 * time.Minute
)

func main() {
	ctx := context.Background()

	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			"localhost",
			"5432",
			"root",
			"root",
			"postgres",
			"disable",
		),
	)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	for i := 0; ; i++ {
		if i >= dbPingRetryLimit {
			log.Fatalf("failed to ping database: %v", err)
		}
		if err = db.Ping(); err == nil {
			break
		}
		time.Sleep(dbPingRetryInterval)
	}

	c1, err := cache.New(ctx,
		cache.RdbSource(db, "codes", []string{"key1", "key2"}, "value"), 2)
	if err != nil {
		log.Fatalf("failed to create codes cache: %v", err)
	}
	go reload(ctx, c1)

	c2, err := cache.New(ctx,
		cache.RdbSource(db, "codes_lang", []string{"key1", "key2", "lang"}, "value"), 3)
	if err != nil {
		log.Fatalf("failed to create codes_lang cache: %v", err)
	}
	go reload(ctx, c2)

	h := handler{
		codesCache:     c1,
		codesLangCache: c2,
	}

	e := echo.New()
	e.GET("/codes", h.codes)
	e.GET("/codes-lang", h.codesLang)
	e.Logger.Fatal(e.Start(":1323"))
}

func reload(ctx context.Context, c *cache.Cache) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(reloadInterval)
		if err := c.Reload(ctx); err != nil {
			log.Printf("failed to reload cache: %v", err)
		}
	}
}

type handler struct {
	codesCache     *cache.Cache
	codesLangCache *cache.Cache
}

func (h *handler) codes(c echo.Context) error {
	str, err := h.codesCache.GetValue(
		c.Request().Context(),
		c.QueryParam("key1"),
		c.QueryParam("key2"))
	if err != nil {
		if err == cache.ErrCodeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Code not found")
		}
		return err
	}

	return c.String(http.StatusOK, str)
}

func (h *handler) codesLang(c echo.Context) error {
	str, err := h.codesLangCache.GetValue(
		c.Request().Context(),
		c.QueryParam("key1"),
		c.QueryParam("key2"),
		c.QueryParam("lang"))
	if err != nil {
		if err == cache.ErrCodeNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Code not found")
		}
		return err
	}

	return c.String(http.StatusOK, str)
}
