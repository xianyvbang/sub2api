package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultOpenAICodexImageGenerationBridgeDisabled(t *testing.T) {
	t.Run("sets disabled by default for openai", func(t *testing.T) {
		extra := withDefaultOpenAICodexImageGenerationBridgeDisabled(PlatformOpenAI, nil)

		require.Equal(t, false, extra[featureKeyCodexImageGenerationBridge])
	})

	t.Run("preserves explicit top level override", func(t *testing.T) {
		extra := map[string]any{featureKeyCodexImageGenerationBridge: true}

		got := withDefaultOpenAICodexImageGenerationBridgeDisabled(PlatformOpenAI, extra)

		require.Equal(t, true, got[featureKeyCodexImageGenerationBridge])
	})

	t.Run("preserves explicit nested override", func(t *testing.T) {
		extra := map[string]any{
			PlatformOpenAI: map[string]any{"codex_image_generation_bridge_enabled": true},
		}

		got := withDefaultOpenAICodexImageGenerationBridgeDisabled(PlatformOpenAI, extra)

		require.NotContains(t, got, featureKeyCodexImageGenerationBridge)
	})

	t.Run("leaves non openai accounts unchanged", func(t *testing.T) {
		extra := map[string]any{}

		got := withDefaultOpenAICodexImageGenerationBridgeDisabled(PlatformAnthropic, extra)

		require.Empty(t, got)
	})
}
