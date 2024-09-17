INSERT INTO apps (id, name, secret)
VALUES ('iotafull', 'test', 'secret')
ON CONFLICT DO NOTHING;