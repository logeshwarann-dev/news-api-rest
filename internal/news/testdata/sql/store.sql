CREATE EXTENSION IF NOT EXIST "uuid-ossp";


CREATE TABLE NEWS IF NOT EXISTS (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    content TEXT NOT NULL,
    source TEXT[] NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
);

INSERT INTO news(id, author, title, summary, content, source, tags, created_at, updated_at)
VALUES(
    'c2f92052-348f-4372-b4bc-43dbbc88445a',
    'Batman',
    'Breaking NEWS',
    'A brief summary of news',
    'Batman is a hero with super powers who saves people from devil',
    'https://www.google.com',
    ARRAY ['marvel', 'sci-fi'],
    NOW(),
    NOW()
);

INSERT INTO news(id, author, title, summary, content, source, tags created_at, updated_at)
VALUES(
    'd2f92090-348f-4372-b4bc-43dbbc88445z',
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

INSERT INTO news(id, author, title, summary, content, source, tags created_at, updated_at, deleted_at)
VALUES(
    'd2f92090-348f-4372-b4bc-43dbbc88445z',
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