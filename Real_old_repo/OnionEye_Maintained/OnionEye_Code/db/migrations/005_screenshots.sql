-- OnionEye — Per-URL Screenshots (render-sandbox pipeline)
--
-- One screenshot per unique URL. UNIQUE(url) is the upsert anchor: a re-capture
-- (UI "CAPTURE SCREENSHOT" click) overwrites the previous PNG instead of
-- accumulating rows. png_data is stored raw BYTEA — pixels only, no
-- executable payload — produced by the air-gapped render-sandbox from deeply
-- sanitized HTML (the raw bytes of a hostile page never leave the crawler).

CREATE TABLE IF NOT EXISTS screenshots (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url        TEXT        NOT NULL,
    png_data   BYTEA       NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_screenshots_url UNIQUE (url)
);

CREATE INDEX IF NOT EXISTS idx_screenshots_url ON screenshots (url);