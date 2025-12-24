-- +migrate Up
CREATE TABLE users (
                       iD int primary key AUTO_INCREMENT,
                       name varchar(255) not null,
                       phone_number varchar(255) not null unique
);

-- +migrate Down
DROP TABLE  users;