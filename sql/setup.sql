create table if not exists ads (
    id        int unsigned auto_increment not null,
    title     varchar(128) not null,
    start_at  timestamp not null,
    end_at    timestamp not null,
    age_start int unsigned not null,
    age_end   int unsigned not null,
    primary key (id)
);

create table if not exists genders (
    id int unsigned auto_increment not null,
    gender varchar(2) not null,
    primary key (id),
    unique key gender (gender)
);

create table if not exists platforms (
    id int unsigned auto_increment not null,
    platform varchar(8) not null,
    primary key (id),
    unique key (platform)
);

CREATE TABLE IF NOT EXISTS `countries` (
    `id` int(3) NOT NULL,
    `country` varchar(2) NOT NULL,
    `alpha3` varchar(3) NOT NULL,
    `langCS` varchar(45) NOT NULL,
    `langDE` varchar(45) NOT NULL,
    `langEN` varchar(45) NOT NULL,
    `langES` varchar(45) NOT NULL,
    `langFR` varchar(45) NOT NULL,
    `langIT` varchar(45) NOT NULL,
    `langNL` varchar(45) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `alpha2` (`country`),
    UNIQUE KEY `alpha3` (`alpha3`)
);


create table if not exists ad_gender (
    ad_id int unsigned not null,
    gender_id int unsigned not null,
    constraint ad_gender_ad foreign key (ad_id) references ads(id),
    constraint ad_gender_gender foreign key (gender_id) references genders(id)
);

create table if not exists ad_platform (
    ad_id int unsigned not null,
    platform_id int unsigned not null,
    constraint ad_platform_ad foreign key (ad_id) references ads(id),
    constraint ad_platform_platform foreign key (platform_id) references platforms(id)
);

create table if not exists ad_country (
    ad_id int unsigned not null,
    country_id int(3) not null,
    constraint ad_contry_ad foreign key (ad_id) references ads(id),
    constraint ad_contry_contry foreign key (country_id) references countries(id)
);