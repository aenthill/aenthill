package jobs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/aenthill/aenthill/app/context"
	"github.com/aenthill/aenthill/app/jobs/stores"

	"github.com/aenthill/log"
	update "github.com/tj/go-update"
)

type selfUpdateJob struct {
	target  string
	version string
	project *update.Manager
	appCtx  *context.AppContext
}

// NewSelfUpdateJob creates a selfUpdateJob instance.
func NewSelfUpdateJob(target string, version string, appCtx *context.AppContext) Job {
	p := &update.Manager{
		Command: "aenthill",
		Store: &github.Store{
			Owner:   "aenthill",
			Repo:    "aenthill",
			Version: version,
		},
	}

	return &selfUpdateJob{target, version, p, appCtx}
}

// Run implements Job.
// The logic was strongly inspired
// by https://github.com/apex/up/blob/master/internal/cli/upgrade/upgrade.go
func (job *selfUpdateJob) Run() error {
	start := time.Now()

	r, err := job.run()
	if err != nil {
		log.Errorf(job.appCtx.EntryContext, err, "job has failed after %0.2fs", time.Since(start).Seconds())
	} else if r != nil {
		log.Infof(job.appCtx.EntryContext, "updated from %s to %s successfully after %0.2fs", job.version, r.Version, time.Since(start).Seconds())
	} else {
		log.Info(job.appCtx.EntryContext, "no updates available")
	}

	return err
}

func (job *selfUpdateJob) run() (*update.Release, error) {
	log.Info(job.appCtx.EntryContext, "fetching release")
	r, err := job.fetchRelease()
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, nil
	}

	log.Infof(job.appCtx.EntryContext, "downloading version %s", r.Version)
	a, err := job.downloadArchive(r)
	if err != nil {
		return nil, err
	}

	log.Infof(job.appCtx.EntryContext, "installing version %s", r.Version)
	return r, job.installArchive(a)
}

func (job *selfUpdateJob) fetchRelease() (*update.Release, error) {
	if job.target == "" {
		return job.fetchLatestRelease()
	}

	return job.project.Store.GetRelease(job.target)
}

func (job *selfUpdateJob) fetchLatestRelease() (*update.Release, error) {
	releases, err := job.project.Store.LatestReleases()
	if err != nil {
		return nil, err
	}

	if len(releases) == 0 {
		return nil, nil
	}

	return releases[0], nil
}

type noBinaryFoundError struct{}

const noBinaryFoundErrorMessage = "failed to find a binary for %s %s"

func (e *noBinaryFoundError) Error() string {
	return fmt.Sprintf(noBinaryFoundErrorMessage, runtime.GOOS, runtime.GOARCH)
}

func (job *selfUpdateJob) downloadArchive(release *update.Release) (string, error) {
	asset := release.FindTarball(runtime.GOOS, runtime.GOARCH)
	if asset == nil {
		return "", &noBinaryFoundError{}
	}

	archive, err := asset.Download()
	if err != nil {
		return "", err
	}

	return archive, nil
}

func (job *selfUpdateJob) installArchive(archive string) error {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}
	dst := filepath.Dir(path)

	return job.project.InstallTo(archive, dst)
}
