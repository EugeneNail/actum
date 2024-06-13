CREATE TABLE activities
(
    id            INT AUTO_INCREMENT,
    name          VARCHAR(100) NOT NULL,
    icon          VARCHAR(50)  NOT NULL,
    user_id       INT          NOT NULL,
    collection_id INT          NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (collection_id) REFERENCES collections (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
)