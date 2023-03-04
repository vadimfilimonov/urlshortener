CREATE TABLE urls
(
  id serial not null unique,
  user_id varchar(255) not null unique,
  shorten_url varchar(255) not null unique,
  original_url varchar(255) not null
);
