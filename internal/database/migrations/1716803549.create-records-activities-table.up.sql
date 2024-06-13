CREATE TABLE records_activities
(
    record_id   INT,
    activity_id INT,

    PRIMARY KEY (record_id, activity_id),
    FOREIGN KEY (record_id) REFERENCES records (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (activity_id) REFERENCES activities (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
)