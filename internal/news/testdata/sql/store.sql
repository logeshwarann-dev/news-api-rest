CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS news (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    content TEXT NOT NULL,
    source TEXT NOT NULL,
    tags TEXT[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


INSERT INTO news (id, author, title, summary, content, source, tags, created_at, updated_at)
VALUES (
    'c2f92252-348f-7374-b4bc-43dbbc88445a',
    'Batman',
    'Breaking NEWS',
    'A brief summary of news',
    'Batman is a hero with super powers who saves people from devil',
    'https://www.google.com',
    ARRAY ['marvel', 'sci-fi'],
    NOW(),
    NOW()
);

INSERT INTO news (id, author, title, summary, content, source, tags created_at, updated_at)
VALUES (
    'd2f92090-348f-4372-b4bg-43dbbc88445z',
    'Superman',
    'Breaking NEWS',
    'A brief summary of news',
    'Superman is a hero with super powers who saves people from devil',
    'https://www.google.com',
    ARRAY ['marvel', 'sci-fi'],
    NOW(),
    NOW()
);

---Deleted News---

INSERT INTO news (id, author, title, summary, content, source, tags created_at, updated_at, deleted_at)
VALUES (
    'n2f92090-348f-4372-b4bv-43dbbc88445i',
    'Spiderman',
    'Breaking NEWS',
    'A brief summary of news',
    'Spiderman is a hero with super powers who saves people from devil',
    'https://www.google.com',
    ARRAY ['marvel', 'sci-fi'],
    NOW(),
    NOW(),
    NOW()
);