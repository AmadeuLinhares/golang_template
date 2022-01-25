-- +goose Up
-- +goose StatementBegin
CREATE TABLE item(
    Id VARCHAR(36) NOT NULL,
    Token VARCHAR(200) NOT NULL,
    UserId VARCHAR(200),
    DeviceId VARCHAR(200),
    Roles VARCHAR(200),
    IsOAuth BOOLEAN NOT NULL DEFAULT true,
    CreatedAt datetime NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    ExpiresAt datetime,
    LastActivityAt datetime,
    PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE item;
-- +goose StatementEnd
