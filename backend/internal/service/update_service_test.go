//go:build unit

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type updateServiceCacheStub struct {
	data string
}

func (s *updateServiceCacheStub) GetUpdateInfo(context.Context) (string, error) {
	if s.data == "" {
		return "", errors.New("cache miss")
	}
	return s.data, nil
}

func (s *updateServiceCacheStub) SetUpdateInfo(_ context.Context, data string, _ time.Duration) error {
	s.data = data
	return nil
}

type updateServiceGitHubClientStub struct {
	release *GitHubRelease
	repo    string
}

func (s *updateServiceGitHubClientStub) FetchLatestRelease(_ context.Context, repo string) (*GitHubRelease, error) {
	s.repo = repo
	return s.release, nil
}

func (s *updateServiceGitHubClientStub) DownloadFile(context.Context, string, string, int64) error {
	panic("DownloadFile should not be called when no update is available")
}

func (s *updateServiceGitHubClientStub) FetchChecksumFile(context.Context, string) ([]byte, error) {
	panic("FetchChecksumFile should not be called when no update is available")
}

func TestUpdateServicePerformUpdateNoUpdateReturnsSentinel(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{
			release: &GitHubRelease{
				TagName: "v0.1.132",
				Name:    "v0.1.132",
			},
		},
		"0.1.132",
		"release",
	)

	err := svc.PerformUpdate(context.Background())

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrNoUpdateAvailable))
	require.ErrorIs(t, err, ErrNoUpdateAvailable)
}

func TestUpdateServiceUsesConfiguredReleaseRepository(t *testing.T) {
	github := &updateServiceGitHubClientStub{
		release: &GitHubRelease{
			TagName: "v0.1.132",
			Name:    "v0.1.132",
		},
	}
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		github,
		"0.1.132",
		"release",
	)

	_, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.Equal(t, "xianyvbang/sub2api", github.repo)
}

func TestCompareVersionsSupportsFourSegmentVersions(t *testing.T) {
	testCases := []struct {
		name     string
		current  string
		latest   string
		expected int
	}{
		{
			name:     "same four segment version",
			current:  "1.2.3.4",
			latest:   "1.2.3.4",
			expected: 0,
		},
		{
			name:     "detects fourth segment upgrade",
			current:  "1.2.3.4",
			latest:   "1.2.3.5",
			expected: -1,
		},
		{
			name:     "detects fourth segment downgrade",
			current:  "1.2.3.5",
			latest:   "1.2.3.4",
			expected: 1,
		},
		{
			name:     "treats missing segment as zero",
			current:  "1.2.3",
			latest:   "1.2.3.0",
			expected: 0,
		},
		{
			name:     "treats prefixed version equivalently",
			current:  "v1.2.3.4",
			latest:   "1.2.3.4",
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, compareVersions(tc.current, tc.latest))
		})
	}
}

func TestUpdateServiceDetectsFourSegmentUpdate(t *testing.T) {
	svc := NewUpdateService(
		&updateServiceCacheStub{},
		&updateServiceGitHubClientStub{
			release: &GitHubRelease{
				TagName: "v1.2.3.5",
				Name:    "v1.2.3.5",
			},
		},
		"1.2.3.4",
		"release",
	)

	info, err := svc.CheckUpdate(context.Background(), true)

	require.NoError(t, err)
	require.True(t, info.HasUpdate)
	require.Equal(t, "1.2.3.5", info.LatestVersion)
}
