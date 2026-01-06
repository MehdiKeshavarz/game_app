-- +migrate Up
CREATE TABLE `users` (
                       `iD` INT PRIMARY KEY AUTO_INCREMENT,
                       `name` VARCHAR(191) NOT NULL ,
                       `phone_number` VARCHAR(191) NOT NULL UNIQUE
);

-- +migrate Down
DROP TABLE  `users`;