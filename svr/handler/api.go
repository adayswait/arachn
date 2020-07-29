package handler

import "time"
import "fmt"
import "github.com/gofiber/fiber"

func Hello(c *fiber.Ctx) {
	fmt.Println(c.Cookies("jesse"))
	c.JSON(fiber.Map{"status": "success", "message": "Hello i'm ok!", "data": nil})
	cookie := new(fiber.Cookie)
	cookie.Name = "jesse"
	cookie.Value = "sdajsdsahjdkhsakd"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.Cookie(cookie)
}

func Login(c *fiber.Ctx) {
	c.JSON(fiber.Map{"status": "success", "message": "Login ok!", "data": nil})
}
