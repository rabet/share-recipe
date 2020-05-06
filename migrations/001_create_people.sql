create table people(
  id serial primary key,
  username varchar(50) UNIQUE not null,
  code integer not null,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_code 
ON people(code);

---- create above / drop below ----

drop table people;
drop index idx_code
