-- New API keys should not enable rate multiplier protection unless requested.
ALTER TABLE api_keys
    ALTER COLUMN rate_protection_enabled SET DEFAULT FALSE;
