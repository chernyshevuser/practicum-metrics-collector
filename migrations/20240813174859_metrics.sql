-- +goose Up
-- +goose StatementBegin
CREATE TABLE  public.Metrics (
    "key" BIGINT PRIMARY KEY,
    "id" TEXT,
    "type" TEXT,
    "val" DOUBLE PRECISION DEFAULT 0,
    "delta" BIGINT DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.Metrics;
-- +goose StatementEnd
