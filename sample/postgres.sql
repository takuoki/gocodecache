-- codes
CREATE TABLE codes (
    key1 VARCHAR(32) NOT NULL,
    key2 VARCHAR(32) NOT NULL,
    value VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    PRIMARY KEY (key1, key2)
);

INSERT INTO codes (key1, key2, value, created_at, updated_at) 
VALUES 
  ('account_type', '1', 'Anonymous account', NOW(), NOW()),
  ('account_type', '2', 'General account', NOW(), NOW()),
  ('account_type', '3', 'Administrator account', NOW(), NOW()),
  ('visibility_level', '1', 'Private', NOW(), NOW()),
  ('visibility_level', '2', 'Public', NOW(), NOW());

-- codes_lang
CREATE TABLE codes_lang (
    key1 VARCHAR(32) NOT NULL,
    key2 VARCHAR(32) NOT NULL,
    lang VARCHAR(32) NOT NULL,
    value VARCHAR(64) NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    PRIMARY KEY (key1, key2, lang)
);

INSERT INTO codes_lang (key1, key2, lang, value, created_at, updated_at) 
VALUES 
  ('account_type', '1', 'en-US', 'Anonymous account', NOW(), NOW()),
  ('account_type', '1', 'ja-JP', '匿名アカウント', NOW(), NOW()),
  ('account_type', '2', 'en-US', 'General account', NOW(), NOW()),
  ('account_type', '2', 'ja-JP', '一般アカウント', NOW(), NOW()),
  ('account_type', '3', 'en-US', 'Administrator account', NOW(), NOW()),
  ('account_type', '3', 'ja-JP', '管理者アカウント', NOW(), NOW()),
  ('visibility_level', '1', 'en-US', 'Private', NOW(), NOW()),
  ('visibility_level', '1', 'ja-JP', '非公開', NOW(), NOW()),
  ('visibility_level', '2', 'en-US', 'Public', NOW(), NOW()),
  ('visibility_level', '2', 'ja-JP', '公開', NOW(), NOW());
