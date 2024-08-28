CREATE TABLE refresh_tokens (
                               id SERIAL PRIMARY KEY,
                               user_id INT NOT NULL,
                               token varchar(250) NOT NULL,
                               expires_at TIMESTAMP NOT NULL,
                               FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
