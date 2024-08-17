-- +goose Up
-- +goose StatementBegin
CREATE TABLE  public.Metrics (
    "id" TEXT PRIMARY KEY,
    "type" TEXT,
    "val" TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.Metrics;
-- +goose StatementEnd
