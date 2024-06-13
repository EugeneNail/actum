CREATE TABLE photos
(
    id int AUTO_INCREMENT,
    name      varchar(50) UNIQUE NOT NULL,
    record_id INT,
    user_id   INT      NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (record_id) REFERENCES records (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
)