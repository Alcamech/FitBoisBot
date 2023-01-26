USE fitbois;

CREATE TABLE fit_boi_record(
    id int AUTO_INCREMENT primary key,
    name varchar(25),
    activity varchar(25),
    month varchar(2),
    day varchar(2),
    year varchar(4)
);

CREATE TABLE fit_boi_user (
    id bigint primary key,
    name varchar(50),
    group_id long,
    fast_gg_count int DEFAULT 0,
    version varchar(50)
);