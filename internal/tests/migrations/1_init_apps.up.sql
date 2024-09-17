INSERT INTO apps (id, name, secret)
VALUES ('a16fcc5e-d4de-4cf9-813f-e7ccf36f29d3', 'test', 'secret')
ON CONFLICT DO NOTHING;