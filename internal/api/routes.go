package api

import (
	"context"
	"net/http"

	"example.com/api/internal/setup"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func InitializeRoutes() {
	cfg := setup.GetConfig()
	cfg.GinEngine = gin.New()

	cfg.GinEngine.POST("api/v1/continent", createContinent(cfg.PgPool.QueryRow))
	cfg.GinEngine.GET("api/v1/continent/:id", getContinent(cfg.PgPool.QueryRow))
	cfg.GinEngine.GET("api/v1/continents", getAllContinents(cfg.PgPool.Query))
	cfg.GinEngine.PUT("api/v1/continent/:id", updateContinent(cfg.PgPool.Exec))
	cfg.GinEngine.DELETE("api/v1/continent/:id", deleteContinent(cfg.PgPool.Exec))

	cfg.GinEngine.POST("api/v1/country", createCountry(cfg.PgPool.QueryRow))
	cfg.GinEngine.GET("api/v1/country/:id", getCountry(cfg.PgPool.QueryRow))
	cfg.GinEngine.GET("api/v1/countries", getAllCountries(cfg.PgPool.Query))
	cfg.GinEngine.PUT("api/v1/country/:id", updateCountry(cfg.PgPool.Exec))
	cfg.GinEngine.DELETE("api/v1/country/:id", deleteCountry(cfg.PgPool.Exec))

	cfg.GinEngine.POST("api/v1/city", createCity(cfg.PgPool.QueryRow))
	cfg.GinEngine.GET("api/v1/city/:id", getCity(cfg.PgPool.QueryRow))
	cfg.GinEngine.GET("api/v1/cities", getAllCities(cfg.PgPool.Query))
	cfg.GinEngine.PUT("api/v1/city/:id", updateCity(cfg.PgPool.Exec))
	cfg.GinEngine.DELETE("api/v1/city/:id", deleteCity(cfg.PgPool.Exec))
}

func createContinent(queryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var id string
		err := queryRowFunc(context.Background(), "INSERT INTO continents (name) VALUES ($1) RETURNING id", input.Name).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func createCountry(queryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name        string `json:"name" binding:"required"`
			ContinentID int    `json:"continent_id" binding:"required"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var id string
		err := queryRowFunc(context.Background(), "INSERT INTO countries (name, continent_id) VALUES ($1, $2) RETURNING id", input.Name, input.ContinentID).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func createCity(queryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name      string `json:"name" binding:"required"`
			CountryID int    `json:"country_id" binding:"required"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var id string
		err := queryRowFunc(context.Background(), "INSERT INTO cities (name, country_id) VALUES ($1, $2) RETURNING id", input.Name, input.CountryID).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func getContinent(queryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var name string
		err := queryRowFunc(context.Background(), "SELECT name FROM continents WHERE id=$1", id).Scan(&name)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Continent not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"name": name})
	}
}

func getCountry(queryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var name string
		var continentID int
		err := queryRowFunc(context.Background(), "SELECT name, continent_id FROM countries WHERE id=$1", id).Scan(&name, &continentID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"name": name, "continent_id": continentID})
	}
}

func getCity(queryRowFunc func(ctx context.Context, sql string, args ...any) pgx.Row) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var name string
		var countryID int
		err := queryRowFunc(context.Background(), "SELECT name, country_id FROM cities WHERE id=$1", id).Scan(&name, &countryID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "City not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"name": name, "country_id": countryID})
	}
}

func getAllContinents(queryFunc func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := queryFunc(context.Background(), "SELECT id, name FROM continents")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var continents = make([]struct {
			ID   int    `json:"id" binding:"required"`
			Name string `json:"name" binding:"required"`
		}, 0)

		for rows.Next() {
			var continent struct {
				ID   int    `json:"id" binding:"required"`
				Name string `json:"name" binding:"required"`
			}
			if err := rows.Scan(&continent.ID, &continent.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			continents = append(continents, continent)
		}

		c.JSON(http.StatusOK, continents)
	}
}

func getAllCountries(queryFunc func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := queryFunc(context.Background(), "SELECT id, name, continent_id FROM countries")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var countries = make([]struct {
			ID          int    `json:"id" binding:"required"`
			Name        string `json:"name" binding:"required"`
			ContinentID int    `json:"continent_id" binding:"required"`
		}, 0)

		for rows.Next() {
			var country struct {
				ID          int    `json:"id" binding:"required"`
				Name        string `json:"name" binding:"required"`
				ContinentID int    `json:"continent_id" binding:"required"`
			}
			if err := rows.Scan(&country.ID, &country.Name, &country.ContinentID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			countries = append(countries, country)
		}

		c.JSON(http.StatusOK, countries)
	}
}

func getAllCities(queryFunc func(ctx context.Context, sql string, args ...any) (pgx.Rows, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := queryFunc(context.Background(), "SELECT id, name, country_id FROM cities")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var cities = make([]struct {
			ID        int    `json:"id" binding:"required"`
			Name      string `json:"name" binding:"required"`
			CountryID int    `json:"country_id" binding:"required"`
		}, 0)

		for rows.Next() {
			var city struct {
				ID        int    `json:"id" binding:"required"`
				Name      string `json:"name" binding:"required"`
				CountryID int    `json:"country_id" binding:"required"`
			}
			if err := rows.Scan(&city.ID, &city.Name, &city.CountryID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			cities = append(cities, city)
		}

		c.JSON(http.StatusOK, cities)
	}
}

func updateContinent(execFunc func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input struct {
			Name string `json:"name"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		_, err := execFunc(context.Background(), "UPDATE continents SET name=$1 WHERE id=$2", input.Name, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

func updateCountry(execFunc func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input struct {
			Name        string `json:"name" binding:"required"`
			ContinentID int    `json:"continent_id" binding:"required"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		_, err := execFunc(context.Background(), "UPDATE countries SET name=$1, continent_id=$2 WHERE id=$3", input.Name, input.ContinentID, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

func updateCity(execFunc func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input struct {
			Name      string `json:"name" binding:"required"`
			CountryID int    `json:"country_id" binding:"required"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := execFunc(context.Background(), "UPDATE cities SET name=$1, country_id=$2 WHERE id=$3", input.Name, input.CountryID, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

func deleteContinent(execFunc func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := execFunc(context.Background(), "DELETE FROM continents WHERE id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

func deleteCountry(execFunc func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := execFunc(context.Background(), "DELETE FROM countries WHERE id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

func deleteCity(execFunc func(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		_, err := execFunc(context.Background(), "DELETE FROM cities WHERE id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}
