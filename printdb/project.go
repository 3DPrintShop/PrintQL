package printdb

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

// ProjectMetadata describes traits about a project.
type ProjectMetadata struct {
	Public bool `json:"public"`
}

// ProjectComponentPage is a paginated list of component IDs for a project.
type ProjectComponentPage struct {
	projectID string
	nextPage  *string
	// Ids of the components
	ComponentIds []string
}

// ProjectMediaPage is a list of media ids for a project.
type ProjectMediaPage struct {
	projectID string
	nextPage  *string
	// The IDs of the media for a project
	MediaIds []string
}

// A Project is a set of data that describes something you can make.
type Project struct {
	ID         string
	Name       string
	Metadata   ProjectMetadata
	Components ProjectComponentPage
	Images     ProjectMediaPage
}

// GetProject retrieves a project based on it's ID.
func (c *Client) GetProject(projectID string) (Project, error) {
	var project Project

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		if b == nil {
			return fmt.Errorf("database not setup")
		}
		pv := b.Bucket([]byte(projectID))

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
			ID:       projectID,
			Name:     string(pv.Get([]byte(Name))),
			Metadata: metadata,
		}

		return nil
	})

	if err != nil {
		return project, err
	}

	components, err := c.GetComponentsForProject(projectID)

	if err != nil {
		return project, err
	}

	project.Components = components
	images, err := c.GetImagesForProject(projectID)
	project.Images = images

	return project, err
}

// A ProjectPage is a list of project IDs as well as the structure needed to retrieve the next page.
type ProjectPage struct {
	// The list of project ids.
	ProjectIDs []string
	nextKey    *string
}

// GetProjects retrieves a paginated list of project IDs.
func (c *Client) GetProjects(pageID *string) (ProjectPage, error) {
	var projectIds []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		c := b.Cursor()

		if pageID != nil {
			c.Seek([]byte(*pageID))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			projectIds = append(projectIds, string(k))
		}

		return nil
	})

	return ProjectPage{ProjectIDs: projectIds, nextKey: nil}, err
}

// GetImagesForProject retrieves a paginated list of images for a project.
func (c *Client) GetImagesForProject(projectID string) (ProjectMediaPage, error) {
	return c.GetNextProjectImagePage(ProjectMediaPage{
		projectID: projectID,
	})
}

// GetNextProjectImagePage retrieves the next page of image IDs for a project.
func (c *Client) GetNextProjectImagePage(request ProjectMediaPage) (ProjectMediaPage, error) {
	var mediaIDs []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		pb := b.Bucket([]byte(request.projectID))
		pcb := pb.Bucket([]byte(ProjectImagesBucket))

		if pcb == nil {
			return nil
		}

		c := pcb.Cursor()

		if request.nextPage != nil {
			c.Seek([]byte(*request.nextPage))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			//Should probably put something in here to restrict to images
			mediaIDs = append(mediaIDs, string(k))
		}

		return nil
	})

	return ProjectMediaPage{MediaIds: mediaIDs, nextPage: nil}, err
}

// GetComponentsForProject retrieves a paginated list of component IDs for a project.
func (c *Client) GetComponentsForProject(projectID string) (ProjectComponentPage, error) {
	return c.GetNextProjectComponentPage(ProjectComponentPage{
		projectID: projectID,
	})
}

// GetNextProjectComponentPage gets the next page of component ids for a project.
func (c *Client) GetNextProjectComponentPage(request ProjectComponentPage) (ProjectComponentPage, error) {
	var componentIds []string

	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		pb := b.Bucket([]byte(request.projectID))
		pcb := pb.Bucket([]byte(ProjectComponentsBucket))

		c := pcb.Cursor()

		if request.nextPage != nil {
			c.Seek([]byte(*request.nextPage))
		}

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			componentIds = append(componentIds, string(k))
		}

		return nil
	})

	return ProjectComponentPage{ComponentIds: componentIds, nextPage: nil}, err
}

// CreateProject creates a new project and saves it returning an ID to that project.
func (c *Client) CreateProject(projectName string) (string, error) {
	id := uuid.New()
	defaultProjectMetadata := ProjectMetadata{
		Public: false,
	}
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		pb, err := b.CreateBucket([]byte(id.String()))

		if err != nil {
			return err
		}

		pcb, err := pb.CreateBucket([]byte(ProjectComponentsBucket))

		if err != nil {
			return err
		}

		if pcb == nil {
			return fmt.Errorf("Failed to create bucket: " + ProjectComponentsBucket)
		}

		err = pb.Put([]byte(Name), []byte(projectName))

		if err != nil {
			return err
		}

		jsonString, err := json.Marshal(defaultProjectMetadata)

		if err != nil {
			return err
		}

		return pb.Put([]byte(Metadata), jsonString)
	})

	return id.String(), err
}

// AssociateComponentWithProject associates a component with a project.
func (c *Client) AssociateComponentWithProject(projectID string, componentID string) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		pb := b.Bucket([]byte(projectID))
		pc := pb.Bucket([]byte(ProjectComponentsBucket))

		err := pc.Put([]byte(componentID), []byte("{}"))

		if err != nil {
			return err
		}

		b = tx.Bucket([]byte(ComponentBucket))
		cb := b.Bucket([]byte(componentID))
		cp := cb.Bucket([]byte(ComponentProjectBucket))

		err = cp.Put([]byte(projectID), []byte("{}"))

		return err
	})

	return err
}

// AssociateImageWithProjectRequest is the data needed to associate an image with a project.
type AssociateImageWithProjectRequest struct {
	// ProjectID is the id for the project to be associated to
	ProjectID string
	// ImageID is the id of the image being associated
	ImageID string
	// Type is the type stored with the relationship so they can be filtered for use in things like galleries vs the card.
	Type string
}

// AssociateImageWithProject associates an image with a project.
func (c *Client) AssociateImageWithProject(request AssociateImageWithProjectRequest) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ProjectBucket))
		pb := b.Bucket([]byte(request.ProjectID))
		pc := pb.Bucket([]byte(ProjectImagesBucket))

		if pc == nil {
			pc, _ = pb.CreateBucket([]byte(ProjectImagesBucket))
		}

		err := pc.Put([]byte(request.ImageID), []byte(request.Type))
		return err
	})

	return err
}
