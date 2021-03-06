-- Database Schema

create table versions (
    applied timestamp without time zone default (now()  at time zone 'utc')
    , file varchar(150)
);

create table maps (
    map_id serial primary key
    , region_name varchar(50)
    , region_type varchar(50)
    , created_at timestamp without time zone default (now() at time zone 'utc')
    , updated_at timestamp without time zone default (now() at time zone 'utc')
    , data json
);

create table users (
    user_id serial primary key 
    , login varchar(50) unique
    , email varchar(50) unique
    , password varchar(256)
    , enabled boolean 
    , key bytea
    , admin boolean
    , last_login timestamp without time zone default (now() at time zone 'utc')
    , data json -- dont know maybe will have user preferences and stuff like that ;)
);

create table corporations (
    corporation_id serial primary key
    , map_id integer references maps on delete cascade
    , data json
    , name varchar(50)
    , user_id integer references users(user_id) on delete set NULL default NULL
);

create table cities (
    city_id serial primary key 
    , map_id integer references maps on delete cascade default NULL 
    , city_name varchar(50) 
    , updated_at timestamp  without time zone default (now() at time zone 'utc')
    , data json
    , corporation_id integer references corporations on delete set NULL default NULL
);

create table neighbouring_cities (
    neighbouring_cities serial primary key
    , from_city_id integer references cities(city_id) on delete cascade
    , to_city_id integer references cities(city_id) on delete cascade
);

create table caravans (
    caravan_id serial primary key
    , origin_corporation_id integer references corporations on delete set NULL default NULL
    , target_corporation_id integer references corporations on delete set NULL default NULL
    , origin_city_id integer references cities(city_id) on delete cascade
    , target_city_id integer references cities(city_id) on delete cascade
    , state integer default 0
    , map_id integer references maps on delete cascade 
    , updated_at timestamp  without time zone default (now() at time zone 'utc')
    , data json
);

create table user_logs (
    user_log_id  serial primary key
    , user_id integer references users(user_id) on delete cascade
    , message varchar(200) 
    , gravity integer 
    , inserted timestamp  without time zone default (now() at time zone 'utc')
    , acknowledged timestamp  without time zone default NULL
);


