CREATE TABLE [examples] (
    [id] [int] NOT NULL DEFAULT 0,
    [text] [nvarchar](260) NOT NULL
);

INSERT INTO [examples] 
    ([text])
VALUES
    ('val1'), ('val2'), ('val3');