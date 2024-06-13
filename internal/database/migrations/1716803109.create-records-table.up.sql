CREATE TABLE records
(
    id      INT AUTO_INCREMENT,
    mood    TINYINT NOT NULL,
    date    DATE    NOT NULL,
    notes   VARCHAR(5000),
    user_id INT     NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
)