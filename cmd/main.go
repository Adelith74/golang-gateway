package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type routes struct {
	Auth      string `json:"auth"`
	Routes    string `json:"routes"`
	Hotels    string `json:"hotels"`
	Events    string `json:"events"`
	Booking   string `json:"booking"`
	Transport string `json:"transport"`
}

func load_routes() ([]byte, error) {
	file_data, err := os.ReadFile("../routes.json")
	if err != nil {
		return nil, err
	} else {
		return file_data, nil
	}
}

func main() {
	json_routes, err := load_routes()
	if err != nil {
		panic(err)
	}
	//записываем в переменную маршруты до API из файла
	var data routes
	json.Unmarshal(json_routes, &data)
	fmt.Println(data)
	r := gin.Default()
	add_route("/auth", data.Auth, r)
	add_route("/transport", data.Transport, r)
	add_route("/routes", data.Routes, r)
	add_route("/hotels", data.Hotels, r)
	add_route("/events", data.Events, r)
	add_route("/booking", data.Booking, r)
	r.Run(":8080")
}

func add_route(path string, dest string, r *gin.Engine) {
	if path == "" || dest == "" {
		log.Fatalf("unable to define route:'%s'", path)
	}

	r.Any(path+"/:params", func(c *gin.Context) {
		params := c.Param("params")

		// Создаем новый HTTP-клиент
		client := &http.Client{}

		// Создаем новый запрос на основе исходного запроса
		req, err := http.NewRequest(c.Request.Method, dest+"/"+params, c.Request.Body)
		if err != nil {
			log.Fatalf("error creating request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Копируем заголовки из исходного запроса в новый запрос
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Отправляем новый запрос
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("error sending request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer resp.Body.Close()

		// Читаем тело ответа
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("error reading response body: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Отправляем ответ обратно клиенту
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	})
}
