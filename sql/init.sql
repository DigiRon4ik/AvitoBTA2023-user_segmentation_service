CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(150),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS segments
(
    id          SERIAL PRIMARY KEY,
    slug        VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_segments
(
    user_id         INT NOT NULL,
    segment_id      INT NOT NULL,
    expiration_time TIMESTAMP,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, segment_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (segment_id) REFERENCES segments (id) ON DELETE CASCADE
);

-- Опционально: история изменений сегментов пользователя
CREATE TABLE IF NOT EXISTS user_segments_history
(
    id         SERIAL PRIMARY KEY,
    user_id    INT                                             NOT NULL,
    segment_id INT                                             NOT NULL,
    action     VARCHAR(10) CHECK (action IN ('ADD', 'REMOVE')) NOT NULL, -- PostgreSQL совместимый ENUM
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (segment_id) REFERENCES segments (id) ON DELETE CASCADE
);
