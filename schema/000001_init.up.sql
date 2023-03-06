CREATE TABLE urls
(
  id serial not null unique,
  user_id varchar(255) not null,
  shorten_url varchar(255) not null,
  original_url varchar(255) not null
);
