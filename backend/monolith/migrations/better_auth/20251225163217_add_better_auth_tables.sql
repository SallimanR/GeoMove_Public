-- migrate:up
create table "user" (
    id text not null primary key,
    name text not null,
    email text not null unique,
    "emailVerified" boolean not null,
    image text,
    "createdAt" timestamptz default CURRENT_TIMESTAMP not null,
    "updatedAt" timestamptz default CURRENT_TIMESTAMP not null
);

create table "account" (
    id text not null primary key,
    "accountId" text not null,
    "providerId" text not null,
    "userId" text not null references "user" (id) on delete cascade,
    "accessToken" text,
    "refreshToken" text,
    "idToken" text,
    "accessTokenExpiresAt" timestamptz,
    "refreshTokenExpiresAt" timestamptz,
    scope text,
    password text,
    "createdAt" timestamptz default CURRENT_TIMESTAMP not null,
    "updatedAt" timestamptz not null
);

create table "verification" (
    id text not null primary key,
    identifier text not null,
    value text not null,
    "expiresAt" timestamptz not null,
    "createdAt" timestamptz default CURRENT_TIMESTAMP not null,
    "updatedAt" timestamptz default CURRENT_TIMESTAMP not null
);

create index account_userId_idx on "account" ("userId");

create index verification_identifier_idx on "verification" (identifier);

-- migrate:down
drop table "verification";
drop table "account";
drop table "user";

-- drop index account_userId_idx;
-- drop index verification_identifier_idx;
