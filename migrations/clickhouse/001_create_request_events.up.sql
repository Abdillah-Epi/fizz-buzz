CREATE TABLE request_events
(
    timestamp DateTime64(3, 'UTC'),

    int1 UInt32,
    int2 UInt32,
    limit_value UInt32,

    str1 String,
    str2 String
)
ENGINE = MergeTree
ORDER BY (timestamp, int1, int2);
