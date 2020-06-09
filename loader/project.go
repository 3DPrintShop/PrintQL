package loader

import (
	"context"
	"github.com/3DPrintShop/PrintQL/errors"
	"github.com/3DPrintShop/PrintQL/printdb"
	"github.com/graph-gophers/dataloader"
	"sync"
)

func LoadProject(ctx context.Context, projectId string) (printdb.Project, error) {
	var project printdb.Project

	ldr, err := extract(ctx, projectLoaderKey)
	if err != nil {
		return project, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(projectId))()
	if err != nil {
		return project, err
	}

	project, ok := data.(printdb.Project)
	if !ok {
		return project, errors.WrongType(project, data)
	}

	return project, nil
}

func LoadProjects(ctx context.Context, projectPageID string) ([]printdb.Project, error) {
	var projects []printdb.Project

	ldr, err := extract(ctx, projectPageLoaderKey)
	if err != nil {
		return projects, err
	}

	data, err := ldr.Load(ctx, dataloader.StringKey(projectPageID))()
	if err != nil {
		return projects, err
	}

	projectPage, ok := data.(printdb.ProjectPage)
	if !ok {
		return projects, errors.WrongType(projectPage, data)
	}

	for _, v := range projectPage.ProjectIds {
		project, err := LoadProject(ctx, v)
		if err != nil {
			return projects, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

type projectPageGetter interface {
	Projects(pageId *string) (printdb.ProjectPage, error)
}

type projectGetter interface {
	Project(projectId string) (printdb.Project, error)
}

type projectPageLoader struct {
	get projectPageGetter
}

type projectLoader struct {
	get projectGetter
}

func newProjectLoader(client projectGetter) dataloader.BatchFunc {
	return projectLoader{get: client}.loadBatch
}

func newProjectPageLoader(client projectPageGetter) dataloader.BatchFunc {
	return projectPageLoader{get: client}.loadBatch
}

func (ldr projectPageLoader) loadBatch(ctx context.Context, pageIDs dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(pageIDs)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, pageID := range pageIDs {
		go func(i int, url dataloader.Key) {
			defer wg.Done()

			idString := pageID.String()
			data, err := ldr.get.Projects(&idString)
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, pageID)
	}

	wg.Wait()

	return results
}

func (ldr projectLoader) loadBatch(ctx context.Context, urls dataloader.Keys) []*dataloader.Result {
	var (
		n       = len(urls)
		results = make([]*dataloader.Result, n)
		wg      sync.WaitGroup
	)

	wg.Add(n)

	for i, url := range urls {
		go func(i int, url dataloader.Key) {
			defer wg.Done()

			data, err := ldr.get.Project(url.String())
			results[i] = &dataloader.Result{Data: data, Error: err}
		}(i, url)
	}

	wg.Wait()

	return results
}
