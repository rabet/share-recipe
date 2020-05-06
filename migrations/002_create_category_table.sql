create table category(
  id serial primary key,
  title varchar(50) UNIQUE not null,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cat_title 
ON category(title);

---- create above / drop below ----

drop table category;
drop index idx_cat_title
