package printdb

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

// Component is something you can build that is a piece of a project.
type Component struct {
	ID       string
	Name     string
	Type     string
	Projects ComponentProjectPage
}

// ComponentPage is a page of identifiers for components.
type ComponentPage struct {
	ComponentIds []string
	NextKey      *string
}

// ComponentProjectPage is a paginated list of project ids that a component belongs to.
type ComponentProjectPage struct {
	ComponentID string
	NextPage    *string
	ProjectIds  []string
}

// Component retrieves a component based on componentID.
func (c *Client) Component(componentID string) (Component, error) {
	var component Component

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ComponentBucket))
		if b == nil {
			return fmt.Errorf("database not setup")
		}
		pv := b.Bucket([]byte(componentID))

		if pv == nil {
			return fmt.Errorf("project by that id does not exist")
		}

		component = Component{
			ID:   componentID,
			Name: string(pv.Get([]byte(Name))),
			Type: string(pv.Get([]byte(Type))),
		}

		return nil
	})

	if err != nil {
		return component, err
	}

	projects, err := c.GetProjectsForComponent(componentID)
	component.Projects = projects

	return component, err
}

// Components gets a paginated list of component IDs starting at pageID.
func (c *Client) Components(pageID *string) (ComponentPage, error) {
	var componentIds []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ComponentBucket))
		c := b.Cursor()

		if pageID != nil {
			c.Seek([]byte(*pageID))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			componentIds = append(componentIds, string(k))
		}

		return nil
	})

	return ComponentPage{ComponentIds: componentIds, NextKey: nil}, err
}

// NewComponentRequest is the required data to create a new component.
type NewComponentRequest struct {
	Name string
	Type string
}

// CreateComponent creates a new component based on the data in request, and returns an ID for the created component.
func (c *Client) CreateComponent(request NewComponentRequest) (string, error) {
	id := uuid.New()
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ComponentBucket))
		cb, err := b.CreateBucket([]byte(id.String()))
		if err != nil {
			return err
		}

		_, err = cb.CreateBucket([]byte(ComponentProjectBucket))

		if err != nil {
			return err
		}

		cb.Put([]byte(Name), []byte(request.Name))
		cb.Put([]byte(Type), []byte(request.Type))
		return nil
	})

	return id.String(), err
}

// GetProjectsForComponent returns a paginated list of project IDs that a component is part of.
func (c *Client) GetProjectsForComponent(componentID string) (ComponentProjectPage, error) {
	return c.GetNextComponentProjectPage(ComponentProjectPage{
		ComponentID: componentID,
	})
}

// GetNextComponentProjectPage returns a paginated list of project IDs that a component is part of.
func (c *Client) GetNextComponentProjectPage(request ComponentProjectPage) (ComponentProjectPage, error) {
	var projectIds []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ComponentBucket))
		cb := b.Bucket([]byte(request.ComponentID))
		cpb := cb.Bucket([]byte(ComponentProjectBucket))

		if cpb == nil {
			return nil
		}

		c := cpb.Cursor()

		if request.NextPage != nil {
			c.Seek([]byte(*request.NextPage))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			projectIds = append(projectIds, string(k))
		}

		return nil
	})

	return ComponentProjectPage{ProjectIds: projectIds, NextPage: nil}, err
}
