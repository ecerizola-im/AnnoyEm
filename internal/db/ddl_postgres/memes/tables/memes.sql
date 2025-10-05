CREATE TABLE IF NOT EXISTS memes.meme
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    user_id bigint NOT NULL,
    original_file_name text NOT NULL,
    mime_type text NOT NULL,
    size_bytes bigint NOT NULL,
    uuid uuid,
    upload_status_id smallint NOT NULL DEFAULT 0,
    category text,
    created_at timestamp NOT NULL DEFAULT now(),
    processed_at timestamp,
    updated_at timestamp NOT NULL DEFAULT now(),
    CONSTRAINT memes_pkey PRIMARY KEY (id),
    CONSTRAINT memes_uuid_key UNIQUE (uuid),
    -- CONSTRAINT memes_user_id_fkey FOREIGN KEY (user_id)
    --     REFERENCES users.user (id) MATCH SIMPLE
    --     ON UPDATE RESTRICT
    --     ON DELETE RESTRICT,
    CONSTRAINT memes_upload_status_id_fkey FOREIGN KEY (upload_status_id)
        REFERENCES memes.file_upload_status (id) MATCH SIMPLE
        ON UPDATE RESTRICT
        ON DELETE RESTRICT,
    CONSTRAINT memes_size_bytes_check CHECK (size_bytes >= 0)
)

CREATE INDEX IF NOT EXISTS idx_memes_user_id ON memes.meme (user_id);
CREATE INDEX IF NOT EXISTS idx_memes_status ON memes.meme (upload_status_id);
CREATE INDEX IF NOT EXISTS idx_memes_created_at ON memes.meme (created_at DESC);