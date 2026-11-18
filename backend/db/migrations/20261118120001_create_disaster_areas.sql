-- +migrate Up
CREATE TABLE disaster_areas (
    id BIGSERIAL PRIMARY KEY,
    disaster_id BIGINT NOT NULL,
    city_name VARCHAR(100) NOT NULL,
    damage_summary TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_disaster_areas_disaster
        FOREIGN KEY (disaster_id)
        REFERENCES disasters(id)
        ON DELETE CASCADE
);

-- インデックス
CREATE INDEX idx_disaster_areas_disaster_id ON disaster_areas(disaster_id);
CREATE INDEX idx_disaster_areas_city_name ON disaster_areas(city_name);

-- コメント
COMMENT ON TABLE disaster_areas IS '災害の被災地域';
COMMENT ON COLUMN disaster_areas.id IS '主キー（連番）';
COMMENT ON COLUMN disaster_areas.disaster_id IS '災害ID（外部キー）';
COMMENT ON COLUMN disaster_areas.city_name IS '市区町村名';
COMMENT ON COLUMN disaster_areas.damage_summary IS '被害概要';

-- +migrate Down
DROP TABLE IF EXISTS disaster_areas;
