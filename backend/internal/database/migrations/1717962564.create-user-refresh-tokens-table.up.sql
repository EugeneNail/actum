CREATE TABLE user_refresh_tokens
(
    uuid       VARCHAR(255) NOT NULL,
    expired_at DATETIME     NOT NULL,
    user_id    INT          NOT NULL,

    PRIMARY KEY (uuid),
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
)