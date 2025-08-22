package handlers

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/nathannkweto/mailer/applog"
	"github.com/nathannkweto/mailer/concurrency"
	"github.com/nathannkweto/mailer/config"
	"github.com/nathannkweto/mailer/email"
	"github.com/nathannkweto/mailer/models"
	"github.com/nathannkweto/mailer/utils"
	"github.com/nathannkweto/mailer/validation"
)

// SendEmailHandler handles /send requests
func SendEmailHandler(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.EmailRequest
		if err := c.BodyParser(&req); err != nil {
			applog.Log.WithError(err).Warn("invalid json body")
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				OK:      false,
				Message: "invalid request body",
				Details: err.Error(),
			})
		}

		// validate
		if err := validation.Validator.Struct(req); err != nil {
			applog.Log.WithError(err).Warn("validation failed")
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				OK:      false,
				Message: "validation error",
				Details: err.Error(),
			})
		}

		// concurrency limit (non-blocking)
		if !concurrency.Acquire() {
			applog.Log.Warn("concurrency limit reached")
			return c.Status(fiber.StatusTooManyRequests).JSON(models.APIResponse{
				OK:      false,
				Message: "too many concurrent send requests",
			})
		}
		defer concurrency.Release()

		// save attachments (if any)
		var attachmentPaths []string
		if len(req.Attachments) > 0 {
			// convert attachments to lightweight struct expected by utils
			inline := make([]struct {
				Filename string `json:"filename" validate:"required"`
				Content  string `json:"content"  validate:"required"`
			}, len(req.Attachments))
			for i, a := range req.Attachments {
				inline[i].Filename = a.Filename
				inline[i].Content = a.Content
			}

			paths, err := utils.SaveBase64Attachments(inline)
			if err != nil {
				applog.Log.WithError(err).Warn("failed to save attachments")
				return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
					OK:      false,
					Message: "invalid attachment",
					Details: err.Error(),
				})
			}
			attachmentPaths = paths
			defer utils.CleanupFiles(attachmentPaths)
		}

		// create context with timeout from cfg
		ctx, cancel := context.WithTimeout(context.Background(), cfg.SendTimeout)
		defer cancel()

		// sending email
		if err := email.SendEmail(ctx, req, attachmentPaths); err != nil {
			applog.Log.WithFields(map[string]interface{}{
				"sender":    req.Sender,
				"smtp":      fmt.Sprintf("%s:%d", req.SMTPServer, req.SMTPPort),
				"recipient": req.Recipient,
				"subject":   req.Subject,
			}).WithError(err).Error("email send failed")
			return c.Status(fiber.StatusBadGateway).JSON(models.APIResponse{
				OK:      false,
				Message: "failed to send email",
				Details: err.Error(),
			})
		}

		applog.Log.WithFields(map[string]interface{}{
			"sender":    req.Sender,
			"recipient": req.Recipient,
			"subject":   req.Subject,
		}).Info("email sent")
		return c.Status(fiber.StatusOK).JSON(models.APIResponse{OK: true, Message: "email sent"})
	}
}
