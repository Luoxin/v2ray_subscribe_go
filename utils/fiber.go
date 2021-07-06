package utils

import (
	"github.com/Ferluci/fast-realip"
	"github.com/gofiber/fiber/v2"
)

func GetRealIpFromCtx(ctx *fiber.Ctx) string {
	return realip.FromRequest(ctx.Context())
}
