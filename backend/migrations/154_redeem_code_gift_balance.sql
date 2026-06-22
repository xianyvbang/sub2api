ALTER TABLE redeem_codes
    ADD COLUMN IF NOT EXISTS usage_limit INT NOT NULL DEFAULT 1;

ALTER TABLE redeem_codes
    ADD COLUMN IF NOT EXISTS per_user_limit INT NOT NULL DEFAULT 1;

ALTER TABLE redeem_codes
    ADD COLUMN IF NOT EXISTS gift_parent_id BIGINT REFERENCES redeem_codes(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_redeem_codes_gift_parent_id
    ON redeem_codes(gift_parent_id);
