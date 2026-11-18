-- +migrate Up
CREATE TABLE earthquake_risks (
    mesh_code VARCHAR(20) PRIMARY KEY,
    risk_level VARCHAR(10) NOT NULL,
    probability_30year DECIMAL(5,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- コメント
COMMENT ON TABLE earthquake_risks IS '地震危険度データ（メッシュ単位）';
COMMENT ON COLUMN earthquake_risks.mesh_code IS 'メッシュコード（主キー）';
COMMENT ON COLUMN earthquake_risks.risk_level IS 'リスクレベル（high/medium/low）';
COMMENT ON COLUMN earthquake_risks.probability_30year IS '30年以内の発生確率（%）';

-- +migrate Down
DROP TABLE IF EXISTS earthquake_risks;
