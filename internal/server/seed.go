package server

import (
	"log/slog"

	authzPorts "github.com/ivan-ca97/life/internal/features/authorization/ports"
	"github.com/ivan-ca97/life/internal/features/user/ports"
)

func seed(userService ports.UserService, roleRepo authzPorts.RoleRepository, email, password string) error {
	if email == "" || password == "" {
		return nil
	}

	user, err := userService.GetByEmail(email)
	if err != nil {
		user, err = userService.Create(email, password)
		if err != nil {
			return err
		}
		slog.Info("admin user seeded", "email", email)
	}

	// Grant the admin role idempotently: AssignRoleByName is not a no-op on
	// conflict, so guard it with UserHasRole so re-seeding never fails and an
	// already-seeded user missing the role is healed.
	hasAdmin, err := roleRepo.UserHasRole(user.Id, "admin")
	if err != nil {
		return err
	}
	if !hasAdmin {
		err = roleRepo.AssignRoleByName(user.Id, "admin")
		if err != nil {
			return err
		}
		slog.Info("admin role granted to seed user", "email", email)
	}

	return nil
}
