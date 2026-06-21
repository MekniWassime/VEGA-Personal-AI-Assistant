CREATE TYPE conversation_state AS ENUM ('idle', 'processing', 'completed', 'errored');
CREATE TYPE conversation_type AS ENUM ('task', 'conversation');

CREATE TABLE conversations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    state       conversation_state NOT NULL DEFAULT 'idle',
    type        conversation_type  NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE workers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
);

CREATE TABLE contexts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    timestamp       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE messages (
    id          SERIAL PRIMARY KEY,
    context_id  UUID        NOT NULL REFERENCES contexts(id) ON DELETE CASCADE,
    role        TEXT        NOT NULL,
    content     TEXT        NOT NULL,
    timestamp   TIMESTAMPTZ NOT NULL DEFAULT now(),
    worker_id   UUID REFERENCES workers(id) ON DELETE SET NULL
);

---- create above / drop below ----

DROP TABLE messages;
DROP TABLE contexts;
DROP TABLE workers;
DROP TABLE conversations;
DROP TYPE conversation_type;
DROP TYPE conversation_state;
