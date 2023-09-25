CREATE TYPE CURRENCY AS ENUM ('USD', 'EUR');

CREATE TABLE account
(
    id         BIGSERIAL PRIMARY KEY,
    owner      TEXT                      NOT NULL,
    balance    BIGINT                    NOT NULL,
    currency   TEXT                      NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL
);

CREATE INDEX idx_owner ON account (owner);

CREATE TABLE entry
(
    id         BIGSERIAL PRIMARY KEY,
    account_id BIGINT REFERENCES account (id) ON DELETE CASCADE NOT NULL,
    amount     BIGINT                                           NOT NULL,
    created_at timestamptz DEFAULT now()                        NOT NULL
);

CREATE INDEX idx_account_id ON entry (account_id);


CREATE TABLE transfer
(
    id              BIGSERIAL PRIMARY KEY,
    from_account_id BIGINT REFERENCES account (id) NOT NULL,
    to_account_id   BIGINT REFERENCES account (id) NOT NULL,
    amount          BIGINT                         NOT NULL,
    created_at      timestamptz DEFAULT now()      NOT NULL
);

CREATE INDEX idx_from_account_id ON transfer (from_account_id);
CREATE INDEX idx_to_account_id ON transfer (to_account_id);
CREATE INDEX idx_from_to_account_id ON transfer (from_account_id, to_account_id);

COMMENT ON COLUMN transfer.amount IS 'must be positive';
COMMENT ON COLUMN entry.amount IS 'can be negative or positive';
