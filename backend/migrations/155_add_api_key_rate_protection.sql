-- Add API key rate multiplier protection.
ALTER TABLE api_keys
    ADD COLUMN IF NOT EXISTS rate_protection_enabled BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE api_keys
    ADD COLUMN IF NOT EXISTS max_rate_multiplier DECIMAL(10,4) NOT NULL DEFAULT 0;

UPDATE api_keys ak
SET max_rate_multiplier = COALESCE(
    (
        SELECT ugrm.rate_multiplier
        FROM user_group_rate_multipliers ugrm
        WHERE ugrm.user_id = ak.user_id
          AND ugrm.group_id = ak.group_id
          AND ugrm.rate_multiplier IS NOT NULL
        LIMIT 1
    ),
    g.rate_multiplier,
    1.0
)
FROM groups g
WHERE ak.group_id = g.id
  AND ak.max_rate_multiplier = 0;

UPDATE api_keys
SET max_rate_multiplier = 1.0
WHERE group_id IS NULL
  AND max_rate_multiplier = 0;

COMMENT ON COLUMN api_keys.rate_protection_enabled IS 'Whether this API key blocks usage when effective group rate exceeds max_rate_multiplier';
COMMENT ON COLUMN api_keys.max_rate_multiplier IS 'Maximum effective group rate multiplier allowed for this API key (0 = unset)';
