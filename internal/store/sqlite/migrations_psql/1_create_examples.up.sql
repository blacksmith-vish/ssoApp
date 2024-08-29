CREATE TABLE examples (
    id SERIAL PRIMARY KEY,
    text VARCHAR(200) NOT NULL
);

INSERT INTO
    examples (text)
VALUES
    ('val1'), ('val2'), ('val3');