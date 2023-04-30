package install

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

var settings *cli.EnvSettings = cli.New()

// RepoAdd adds repo with given name and url
func RepoAdd(name, url string, log *logrus.Entry) string {
	repoFile := settings.RepositoryConfig

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error(err.Error())
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer func() {
			if err := fileLock.Unlock(); err != nil {
				log.Error(err.Error())
			}
		}()
	}
	if err != nil {
		log.Error(err.Error())
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		log.Error(err.Error())
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		log.Error(err.Error())
	}

	if f.Has(name) {
		log.Info("repository name (%s) already exists\n", name)
		return "repository name already exists"
	}

	c := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		log.Error(err.Error())
	}

	if _, err := r.DownloadIndexFile(); err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
		log.Error(err.Error())
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		log.Error(err.Error())
	}
	log.Info("%q has been added to your repositories\n", name)
	return "Repo added"
}

// RepoUpdate updates charts for all helm repos
func RepoUpdate(log *logrus.Entry) string {
	repoFile := settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		defer func() {
			if r := recover(); r != nil {
				log.Error("no repositories found. You must add one before updating)", r)
			}
		}()
		log.Panic("no repositories found. You must add one before updating. Error: " + err.Error())
		return "no repositories found. You must add one before updating"
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			defer func() {
				if r := recover(); r != nil {
					log.Error("Recovered in RepoUpdate()", r)
				}
			}()
			log.Panic("Error: " + err.Error())
			return err.Error()
		}
		repos = append(repos, r)
	}

	log.Info("Hang tight while we grab the latest from your chart repositories...\n")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				log.Error("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
				log.Error(err.Error())
			} else {
				log.Info("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	log.Info("Update Complete. ⎈ Happy Helming!⎈\n")
	return "Update Complete. ⎈ Happy Helming!⎈"
}

// InstallChart
func InstallChart(name, repo, chart, namespace string, log *logrus.Entry) string {
	os.Setenv("HELM_NAMESPACE", namespace)
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	client := action.NewInstall(actionConfig)

	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}
	//name, chart, err := client.NameAndChart(args)
	client.ReleaseName = name
	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", repo, chart), settings)
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}

	debug("CHART PATH: %s\n", cp)

	p := getter.All(settings)
	valueOpts := &values.Options{}
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}

	// Add args
	// if err := strvals.ParseInto(args["set"], vals); err != nil {
	// 	log.Fatal(errors.Wrap(err, "failed parsing --set data"))
	// }

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		log.Error(err.Error())
		return err.Error()
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stdout,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					log.Error(err.Error())
				}
			} else {
				log.Error(err.Error())
			}
		}
	}

	client.Namespace = settings.Namespace()
	release, err := client.Run(chartRequested, vals)
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	log.Info(release.Manifest)
	return "Chart installed"
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	err := log.Output(2, fmt.Sprintf(format, v...))
	if err != nil {
		log.Print(err.Error())
	}
}

func DeleteChart(name, namespace string, log *logrus.Entry) string {
	os.Setenv("HELM_NAMESPACE", namespace)
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug); err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	client := action.NewUninstall(actionConfig)
	res, err := client.Run(name)
	if err != nil {
		log.Error(err.Error())
		return err.Error()
	}
	return res.Info
}
