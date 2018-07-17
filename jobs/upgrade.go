package jobs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/aenthill/aenthill/errors"

	"github.com/thegomachine/github-store"
	update "github.com/tj/go-update"
	"github.com/tj/go-update/progress"
)

type upgradeJob struct {
	target  string
	version string
	project *update.Manager
}

func NewUpgradeJob(target, version string) Job {
	p := &update.Manager{
		Command: "aenthill",
		Store: &github.Store{
			Owner:   "aenthill",
			Repo:    "aenthill",
			Version: version,
		},
	}
	return &upgradeJob{target, version, p}
}

func (j *upgradeJob) Execute() error {
	fmt.Println("Fetching release...")
	r, err := j.fetchRelease()
	if err != nil {
		return err
	}
	if r == nil {
		fmt.Println("No updates available")
		return nil
	}
	fmt.Println("Downloading release...")
	a, err := j.downloadArchive(r)
	if err != nil {
		return err
	}
	fmt.Println("Installing release...")
	if err := j.installArchive(a); err != nil {
		return err
	}
	fmt.Printf(`Aenthill version "%s" installed!`, r.Version)
	return nil
}

func (j *upgradeJob) fetchRelease() (*update.Release, error) {
	if j.target == "" {
		return j.fetchLatestRelease()
	}
	return j.project.Store.GetRelease(j.target)
}

func (j *upgradeJob) fetchLatestRelease() (*update.Release, error) {
	releases, err := j.project.Store.LatestReleases()
	if err != nil {
		return nil, errors.Wrap("upgrade job", err)
	}
	if len(releases) == 0 {
		return nil, nil
	}
	return releases[0], nil
}

func (j *upgradeJob) downloadArchive(release *update.Release) (string, error) {
	asset := release.FindTarball(runtime.GOOS, runtime.GOARCH)
	if asset == nil {
		return "", errors.Errorf("upgrade job", `failed to find a binary for "%s" "%s"`, runtime.GOOS, runtime.GOARCH)
	}
	archive, err := asset.DownloadProxy(progress.Reader)
	if err != nil {
		return "", errors.Wrap("upgrade job", err)
	}
	return archive, nil
}

func (j *upgradeJob) installArchive(archive string) error {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return errors.Wrap("upgrade job", err)
	}
	dst := filepath.Dir(path)
	return errors.Wrap("upgrade job", j.project.InstallTo(archive, dst))
}
