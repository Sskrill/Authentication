CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL ,
                       password VARCHAR(255) NOT NULL,
                       registered_at DATE NOT NULL ,
                       email VARCHAR(255) UNIQUE NOT NULL
);
