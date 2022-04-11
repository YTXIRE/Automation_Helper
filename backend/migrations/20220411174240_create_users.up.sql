CREATE table users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    email varchar not null unique,
    password varchar not null,
    login varchar not null unique,
    created_at integer DEFAULT null,
    updated_at integer DEFAULT null,
    last_login integer DEFAULT null
);