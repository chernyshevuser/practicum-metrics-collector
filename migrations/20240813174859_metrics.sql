-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.counterMetrics (
    id TEXT PRIMARY KEY,
    delta BIGINT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS public.gaugeMetrics (
    id TEXT PRIMARY KEY,
    value DOUBLE PRECISION DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.counterMetrics;

DROP TABLE public.gaugeMetrics;
-- +goose StatementEnd
