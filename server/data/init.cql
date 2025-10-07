CREATE KEYSPACE IF NOT EXISTS main
    WITH REPLICATION = {
        'class' : 'SimpleStrategy', 'replication_factor' : 1
    };

CREATE TABLE IF NOT EXISTS main.account (
    session    text,
    resource   uuid,
    port       bigint,
    PRIMARY KEY (session)
);

CREATE TABLE IF NOT EXISTS main.event (
    id text,
    PRIMARY KEY (id)
);