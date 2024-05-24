package authz

import (
	"context"
	"errors"
	"github.com/google/uuid"
)

func GetAuthzUser(ctx context.Context) (uuid.UUID, error) {
	userId, ok := ctx.Value("AuthorizedUserId").(uuid.UUID)

	if !ok {
		return userId, errors.New("failed to obtain AuthorizedUserId from context")
	}

	return userId, nil
}
