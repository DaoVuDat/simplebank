CREATE TABLE "user" (
                        username text primary key ,
                        hashed_password text not null,
                        full_name text not null ,
                        email text unique not null ,
                        password_changed_at timestamptz not null default '0001-01-01 00:00:00Z',
                        created_at timestamptz not null default now()
);

create unique index on account(owner, currency);

alter table account add foreign key ("owner") references "user"("username");