-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.CounterMetrics (
    type TEXT,
    id TEXT PRIMARY KEY,
    delta BIGINT DEFAULT 0,
    value DOUBLE PRECISION DEFAULT 0
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
