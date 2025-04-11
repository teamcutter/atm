package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/teamcutter/atm/client"
)

func main() {
	cli, err := client.New("localhost:8001", "user", "12345")
	if err != nil {
		panic(err)
	}

	app := echo.New()
	app.GET("/set/:key/:value", func(ctx echo.Context) error {
		key := ctx.Param("key")
		value := ctx.Param("value")

		cli.Set(key, value)
		return ctx.JSON(200, map[string]string{
			"key": key,
			"value": value,
		})
	})

	app.GET("/get/:key", func(ctx echo.Context) error {
		key := ctx.Param("key")

		value, err := cli.Get(key)
		if err != nil {
			log.Println(err)
		}
		return ctx.JSON(200, map[string]string{
			"key": key,
			"value": value,
		})
	})

	log.Fatal(app.Start(":8080"))
}