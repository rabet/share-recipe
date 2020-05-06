-- Write your migrate up statements here
create table recipe(
  id serial primary key,
  people_id INTEGER REFERENCES people(id),
  category_id INTEGER REFERENCES category(id),
  title varchar(50) not null,
  descr text,
  link text,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_title 
ON recipe(title);

---- create above / drop below ----

drop table recipe;
drop index idx_title
