CREATE TABLE IF NOT EXISTS memes.file_upload_status
(
    id smallint NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT file_processing_statuses_pkey PRIMARY KEY (id),
    CONSTRAINT file_processing_statuses_name_key UNIQUE (name)
)

INSERT INTO memes.file_upload_status (id, name) VALUES
(1, 'pending'),
(2, 'processed'),
(3, 'failed')