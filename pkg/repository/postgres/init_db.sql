CREATE TABLE segment
(
    slug VARCHAR(255) PRIMARY KEY
);

CREATE TABLE users
(
    id INT PRIMARY KEY
);

CREATE TABLE user_segment
(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER NOT NULL,
    segment_slug VARCHAR(255) NOT NULL,
    add_time TIMESTAMP NOT NULL,
    expire_time TIMESTAMP,
    FOREIGN KEY (segment_slug) REFERENCES segment(slug),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

INSERT INTO users (id) VALUES (1);

INSERT INTO users (id) VALUES (2);

INSERT INTO users (id) VALUES (3);

INSERT INTO segment (slug) VALUES ('test_slug1');

INSERT INTO segment (slug) VALUES ('test_slug2');

INSERT INTO segment (slug) VALUES ('test_slug3');

INSERT INTO user_segment
    (user_id, segment_slug, add_time, expire_time)
VALUES
    (1, 'test_slug1', now(), now() + INTERVAL '20 hours');

INSERT INTO user_segment
    (user_id, segment_slug, add_time, expire_time)
VALUES
    (1, 'test_slug2', now(), now() + INTERVAL '1 hours');

INSERT INTO user_segment
    (user_id, segment_slug, add_time, expire_time)
VALUES
    (1, 'test_slug3', now(), now() - INTERVAL '1 hour');

INSERT INTO user_segment
    (user_id, segment_slug, add_time, expire_time)
VALUES
    (2, 'test_slug2', now(), now() + INTERVAL '1 hour');

INSERT INTO user_segment
    (user_id, segment_slug, add_time, expire_time)
VALUES
    (2, 'test_slug3', now(), now() + INTERVAL '3 hour');

INSERT INTO user_segment
    (user_id, segment_slug, add_time, expire_time)
VALUES
    (3, 'test_slug3', now(), now() + INTERVAL '1 hour');







