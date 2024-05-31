UPDATE activities
SET icon = 0;

ALTER TABLE activities
    MODIFY COLUMN icon VARCHAR(50) NOT NULL