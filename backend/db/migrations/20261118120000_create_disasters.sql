-- +migrate Up
CREATE TABLE disasters (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL,
    damage TEXT,
    prefecture_code VARCHAR(2) NOT NULL,
    city_name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- インデックス
CREATE INDEX idx_disasters_type ON disasters(type);
CREATE INDEX idx_disasters_year ON disasters(year);
CREATE INDEX idx_disasters_prefecture_code ON disasters(prefecture_code);

-- コメント
COMMENT ON TABLE disasters IS '過去の災害データ';
COMMENT ON COLUMN disasters.id IS '主キー（連番）';
COMMENT ON COLUMN disasters.type IS '災害種別（earthquake/flood/tsunami/landslide/typhoon/other）';
COMMENT ON COLUMN disasters.name IS '災害名';
COMMENT ON COLUMN disasters.year IS '発生年';
COMMENT ON COLUMN disasters.damage IS '被害内容';
COMMENT ON COLUMN disasters.prefecture_code IS '都道府県コード（01-47）';
COMMENT ON COLUMN disasters.city_name IS '市区町村名';

-- +migrate Down
DROP TABLE IF EXISTS disasters;
