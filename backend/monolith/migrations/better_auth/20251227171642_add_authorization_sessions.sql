-- migrate:up
create table session (
    id text not null primary key,
    "expiresAt" timestamptz not null,
    token text not null unique,
    "createdAt" timestamptz default CURRENT_TIMESTAMP not null,
    "updatedAt" timestamptz not null,
    "ipAddress" text,
    "userAgent" text,
    "userId" text not null references "user" (id) on delete cascade
);

create index session_userId_idx on session ("userId");

-- migrate:down
drop table session;

-- drop index "session_userId_idx";
