package printdb

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

type ProjectMetadata struct {
	Public bool `json:"public"`
}

type ProjectComponentPage struct {
	ProjectId    string
	NextPage     *string
	ComponentIds []string
}

type ProjectMediaPage struct {
	ProjectId string
	NextPage  *string
	MediaIds  []string
}

type Project struct {
	ID         string
	Name       string
	Metadata   ProjectMetadata
	Components ProjectComponentPage
	Images     ProjectMediaPage
}

func (c *Client) Project(projectId string) (Project, error) {
	var project Project

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		if b == nil {
			return fmt.Errorf("database not setup")
		}
		pv := b.Bucket([]byte(projectId))

		if pv == nil {
			return fmt.Errorf("project by that id does not exist")
		}

		var metadata ProjectMetadata

		metadataBytes := pv.Get([]byte(Metadata))
		if metadataBytes != nil {
			err := json.Unmarshal(metadataBytes, &metadata)
			if err != nil {
				return err
			}
		} else {
			metadata = ProjectMetadata{Public: false}
		}

		project = Project{
			ID:       projectId,
			Name:     string(pv.Get([]byte(Name))),
			Metadata: metadata,
		}

		return nil
	})

	if err != nil {
		return project, err
	}

	components, err := c.GetComponentsForProject(projectId)
	project.Components = components
	images, err := c.GetImagesForProject(projectId)
	project.Images = images

	return project, err
}

type ProjectPage struct {
	ProjectIds []string
	NextKey    *string
}

func (c *Client) Projects(pageId *string) (ProjectPage, error) {
	var projectIds []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		c := b.Cursor()

		if pageId != nil {
			c.Seek([]byte(*pageId))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			fmt.Println(string(k))
			projectIds = append(projectIds, string(k))
		}

		return nil
	})

	return ProjectPage{ProjectIds: projectIds, NextKey: nil}, err
}

type NewProjectRequest struct {
	Name string
}

func (c *Client) GetImagesForProject(projectID string) (ProjectMediaPage, error) {
	return c.GetNextProjectImagePage(ProjectMediaPage{
		ProjectId: projectID,
	})
}

func (c *Client) GetNextProjectImagePage(request ProjectMediaPage) (ProjectMediaPage, error) {
	var mediaIDs []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		pb := b.Bucket([]byte(request.ProjectId))
		pcb := pb.Bucket([]byte(ProjectImagesBucket))

		if pcb == nil {
			fmt.Printf("Failed to get pcb for project")
			return nil
		}

		c := pcb.Cursor()

		if request.NextPage != nil {
			c.Seek([]byte(*request.NextPage))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			//Should probably put something in here to restrict to images
			fmt.Println(string(k))
			mediaIDs = append(mediaIDs, string(k))
		}

		return nil
	})

	return ProjectMediaPage{MediaIds: mediaIDs, NextPage: nil}, err
}

func (c *Client) GetComponentsForProject(projectId string) (ProjectComponentPage, error) {
	return c.GetNextProjectComponentPage(ProjectComponentPage{
		ProjectId: projectId,
	})
}

func (c *Client) GetNextProjectComponentPage(request ProjectComponentPage) (ProjectComponentPage, error) {
	var componentIds []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		pb := b.Bucket([]byte(request.ProjectId))
		pcb := pb.Bucket([]byte(ProjectComponentsBucket))

		if pcb == nil {
			fmt.Printf("Failed to get pcb for project")
			return nil
		}

		c := pcb.Cursor()

		if request.NextPage != nil {
			c.Seek([]byte(*request.NextPage))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			fmt.Println(string(k))
			componentIds = append(componentIds, string(k))
		}

		return nil
	})

	return ProjectComponentPage{ComponentIds: componentIds, NextPage: nil}, err
}

func (c *Client) CreateProject(request NewProjectRequest) (string, error) {
	id := uuid.New()
	defaultProjectMetadata := ProjectMetadata{
		Public: false,
	}
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		pb, err := b.CreateBucket([]byte(id.String()))

		if err != nil {
			fmt.Printf("Error creating project bucket")
			return err
		}

		pcb, err := pb.CreateBucket([]byte(ProjectComponentsBucket))

		if err != nil {
			return err
		}

		if pcb == nil {
			return fmt.Errorf("Failed to create bucket: " + ProjectComponentsBucket)
		}

		pb.Put([]byte(Name), []byte(request.Name))
		json, err := json.Marshal(defaultProjectMetadata)
		return pb.Put([]byte(Metadata), json)
	})

	return id.String(), err
}

type AssociateComponentWithProjectRequest struct {
	ProjectId   string
	ComponentId string
}

func (c *Client) AssociateComponentWithProject(request AssociateComponentWithProjectRequest) error {
	fmt.Printf("Creating association for %s and %s\n", request.ProjectId, request.ComponentId)
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		pb := b.Bucket([]byte(request.ProjectId))
		pc := pb.Bucket([]byte(ProjectComponentsBucket))

		if pc == nil {
			pc, _ = pb.CreateBucket([]byte(ProjectComponentsBucket))
		}

		err := pc.Put([]byte(request.ComponentId), []byte("{}"))

		if err != nil {
			return err
		}

		b = tx.Bucket([]byte(ComponentBucket))
		cb := b.Bucket([]byte(request.ComponentId))
		cp := cb.Bucket([]byte(ComponentProjectBucket))

		if cp == nil {
			cp, _ = cb.CreateBucket([]byte(ComponentProjectBucket))
		}

		err = cp.Put([]byte(request.ProjectId), []byte("{}"))

		return err
	})

	return err
}

type AssociateImageWithProjectRequest struct {
	ProjectId string
	ImageId   string
	Type      string
}

func (c *Client) AssociateImageWithProject(request AssociateImageWithProjectRequest) error {
	fmt.Printf("Creating association for %s and image %s\n", request.ProjectId, request.ImageId)
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		pb := b.Bucket([]byte(request.ProjectId))
		pc := pb.Bucket([]byte(ProjectImagesBucket))

		if pc == nil {
			pc, _ = pb.CreateBucket([]byte(ProjectImagesBucket))
		}

		err := pc.Put([]byte(request.ImageId), []byte(request.Type))
		return err
	})

	return err
}
