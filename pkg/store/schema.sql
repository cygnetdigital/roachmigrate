CREATE TABLE roach_migrations (
  id           STRING PRIMARY KEY,
  key          STRING NOT NULL,
  filename     STRING NOT NULL,
  completed    BOOL NOT NULL DEFAULT false,
  failed       BOOL NOT NULL DEFAULT false,
  fail_reason  STRING NULL,
  created_at    TIMESTAMP NOT NULL DEFAULT now(),
  updated_at    TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX key_listing_idx ON roach_migrations (key);