USE fitbois;

CREATE TABLE fit_boi_record(
    id int AUTO_INCREMENT primary key,
    user_id bigint,
    activity varchar(25),
    month varchar(2),
    day varchar(2),
    year varchar(4),
    FOREIGN KEY (user_id) REFERENCES fit_boi_user(id)
);

CREATE TABLE fit_boi_user (
    id bigint primary key,
    name varchar(50),
    group_id long,
    version varchar(50)
);

CREATE TABLE fit_boi_gg (
    id int AUTO_INCREMENT primary key,
    user_id bigint,
    group_id long,
    year varchar(4),
    fast_gg_count int DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES fit_boi_user(id)
);