-- +goose Up
-- +goose StatementBegin
INSERT INTO item (ID, Token, UserId, DeviceId, Roles, IsOAuth, CreatedAt, ExpiresAt, LastActivityAt) 
    VALUES (0, 123, 123, 123, 123, FALSE, CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM item WHERE id = 0;
-- +goose StatementEnd