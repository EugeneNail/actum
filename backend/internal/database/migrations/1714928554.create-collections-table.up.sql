CREATE TABLE collections
(
    id      INT AUTO_INCREMENT,
    name    VARCHAR(255) NOT NULL,
    user_id INT          NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
)