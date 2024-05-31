UPDATE activities
SET icon = '';

ALTER TABLE `activities`
    MODIFY COLUMN `icon` SMALLINT NOT NULL
