create table musicians (
        id uuid primary key, 
        name varchar(250) 
);

create table musicians_tracks (
        id uuid primary key ,
        musician_id uuid references musicians(id), 
        title varchar(250), 
        played_times bigint default 0
);

create table musicians_plays_by_month (
        musician_id uuid references musicians(id) primary key, 
        month smallint default 1, 
        year integer default 1970,
        plays_total bigint default 0
);