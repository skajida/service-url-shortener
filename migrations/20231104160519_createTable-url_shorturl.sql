CREATE TABLE url_shorturl(
    url_short_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    origin_url   VARCHAR(255) UNIQUE,
    short_url    VARCHAR(10) UNIQUE
);
