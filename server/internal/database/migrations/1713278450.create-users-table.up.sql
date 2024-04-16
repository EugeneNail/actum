CREATE TABLE users
(
    id       INT AUTO_INCREMENT,
    name     VARCHAR(255)  NOT NULL,
    email    VARCHAR(255)  NOT NULL,
    password VARCHAR(1024) NOT NULL,
    PRIMARY KEY (id)
)