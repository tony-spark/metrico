CREATE TABLE gauges (
    name VARCHAR NOT NULL PRIMARY KEY,
    value DOUBLE PRECISION NOT NULL
);

CREATE TABLE counters (
    name VARCHAR NOT NULL PRIMARY KEY,
    value BIGINT NOT NULL
);