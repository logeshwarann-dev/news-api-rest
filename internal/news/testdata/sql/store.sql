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
    '67a7f4b7-2261-4121-a578-bf9da06aa0f3',
    'Batman',
    'Breaking NEWS',
    'A brief summary of news',
    'Batman is a hero with super powers who saves people from devil',
    'https://www.google.com',
    ARRAY ['marvel', 'sci-fi'],
    NOW(),
    NOW()
);

INSERT INTO news (id, author, title, summary, content, source, tags, created_at, updated_at)
VALUES (
    '79d8b51d-97dd-48f4-b619-1f637f15b395',
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

INSERT INTO news (id, author, title, summary, content, source, tags, created_at, updated_at, deleted_at)
VALUES (
    'd635e974-be17-49a9-b962-d508381939ed',
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