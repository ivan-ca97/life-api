package server

import (
	"log/slog"

	"github.com/ivan-ca97/life/internal/features/user/ports"
)

func seed(userService ports.UserService, email, password string) error {
	if email == "" || password == "" {
		return nil
	}

	_, err := userService.GetByEmail(email)
	if err == nil {
		return nil
	}

	_, err = userService.Create(email, password)
	if err != nil {
		return err
	}

	slog.Info("admin user seeded", "email", email)
	return nil
}
