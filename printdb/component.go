package printdb

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

type Component struct {
	ID       string
	Name     string
	Type     string
	Projects ComponentProjectPage
}

type ComponentPage struct {
	ComponentIds []string
	NextKey      *string
}

type ComponentProjectPage struct {
	ComponentId string
	NextPage    *string
	ProjectIds  []string
}

func (c *Client) Component(componentId string) (Component, error) {
	var component Component

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ComponentBucket))
		if b == nil {
			return fmt.Errorf("database not setup")
		}
		pv := b.Bucket([]byte(componentId))

		if pv == nil {
			return fmt.Errorf("project by that id does not exist")
		}

		component = Component{
			ID:   componentId,
			Name: string(pv.Get([]byte(Name))),
			Type: string(pv.Get([]byte(Type))),
		}

		return nil
	})

	projects, err := c.GetProjectsForComponent(componentId)
	component.Projects = projects

	return component, err
}

func (c *Client) Components(pageId *string) (ComponentPage, error) {
	var componentIds []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ComponentBucket))
		c := b.Cursor()

		if pageId != nil {
			c.Seek([]byte(*pageId))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			fmt.Println(string(k))
			componentIds = append(componentIds, string(k))
		}

		return nil
	})

	return ComponentPage{ComponentIds: componentIds, NextKey: nil}, err
}

type NewComponentRequest struct {
	Name string
	Type string
}

func (c *Client) CreateComponent(request NewComponentRequest) (string, error) {
	id := uuid.New()
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ComponentBucket))
		cb, err := b.CreateBucket([]byte(id.String()))
		if err != nil {
			fmt.Printf("Error creating component bucket")
			return err
		}

		_, err = cb.CreateBucket([]byte(ComponentProjectBucket))

		if err != nil {
			fmt.Printf("Error creating component project's bucket")
			return err
		}

		cb.Put([]byte(Name), []byte(request.Name))
		cb.Put([]byte(Type), []byte(request.Type))
		return nil
	})

	return id.String(), err
}

func (c *Client) GetProjectsForComponent(projectId string) (ComponentProjectPage, error) {
	return c.GetNextComponentProjectPage(ComponentProjectPage{
		ComponentId: projectId,
	})
}

func (c *Client) GetNextComponentProjectPage(request ComponentProjectPage) (ComponentProjectPage, error) {
	var projectIds []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ComponentBucket))
		cb := b.Bucket([]byte(request.ComponentId))
		cpb := cb.Bucket([]byte(ComponentProjectBucket))

		if cpb == nil {
			fmt.Printf("Failed to get pcb for project")
			return nil
		}

		c := cpb.Cursor()

		if request.NextPage != nil {
			c.Seek([]byte(*request.NextPage))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			fmt.Println(string(k))
			projectIds = append(projectIds, string(k))
		}

		return nil
	})

	return ComponentProjectPage{ProjectIds: projectIds, NextPage: nil}, err
}
