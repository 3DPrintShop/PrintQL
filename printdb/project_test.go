package printdb

import (
	"github.com/boltdb/bolt"
	"testing"
)

func TestClient_CreateProject(t *testing.T) {
	context, err := setup()
	defer teardown(context)

	if err != nil {
		t.Error(err)
	}

	client, err := NewClient(context.db)

	if err != nil {
		t.Error(err)
		return
	}

	projectId, err := client.CreateProject(NewProjectRequest{
		Name: TestName,
	})

	if err != nil {
		t.Error(err)
		return
	}

	context.db.View(func(tx *bolt.Tx) error {
		pb := tx.Bucket([]byte(ProjectBucket))
		if pb == nil {
			t.Error("Project bucket doesn't sexist within bolt database")
		}

		p := pb.Bucket([]byte(projectId))
		if p == nil {
			t.Error("Bucket wasn't created for projectId: " + projectId)
		}

		pcb := p.Bucket([]byte(ProjectComponentsBucket))
		if pcb == nil {
			t.Error()
		}

		metadataString := string(p.Get([]byte(Metadata)))
		if metadataString != "{\"public\":false}" {
			t.Errorf("Expected metadata to be {}, but was %s", metadataString)
		}

		return nil
	})
}

func TestClient_AddComponentToProject(t *testing.T) {
	context, err := setup()
	defer teardown(context)

	if err != nil {
		t.Error(err)
	}

	client, err := NewClient(context.db)

	if err != nil {
		t.Error(err)
		return
	}

	projectId, err := client.CreateProject(NewProjectRequest{
		Name: TestName,
	})

	if err != nil {
		t.Error(err)
		return
	}

	componentId, err := client.CreateComponent(NewComponentRequest{
		Name: TestComponentName,
		Type: TestComponentType,
	})

	err = client.AssociateComponentWithProject(AssociateComponentWithProjectRequest{
		ComponentId: componentId,
		ProjectId:   projectId,
	})

	if err != nil {
		t.Error(err)
	}

	context.db.View(func(tx *bolt.Tx) error {
		pb := tx.Bucket([]byte(ProjectBucket))
		if pb == nil {
			t.Error("Project bucket doesn't sexist within bolt database")
		}

		p := pb.Bucket([]byte(projectId))
		if p == nil {
			t.Error("Bucket wasn't created for projectId: " + projectId)
		}

		pcb := p.Bucket([]byte(ProjectComponentsBucket))
		if pcb == nil {
			t.Error()
		}

		value := string(pcb.Get([]byte(componentId)))

		if value != "{}" {
			t.Error("Failed to retrieve component association from expected bucket")
		}

		return nil
	})
}

func TestClient_GetProject(t *testing.T) {
	context, err := setup()
	defer teardown(context)

	if err != nil {
		t.Error(err)
	}

	client, err := NewClient(context.db)

	if err != nil {
		t.Error(err)
		return
	}

	projectId, err := client.CreateProject(NewProjectRequest{
		Name: TestName,
	})

	if err != nil {
		t.Error(err)
		return
	}

	compId, err := client.CreateComponent(NewComponentRequest{
		Type: "STL",
		Name: "Component!",
	})

	err = client.AssociateComponentWithProject(AssociateComponentWithProjectRequest{
		ComponentId: compId,
		ProjectId:   projectId,
	})

	if err != nil {
		t.Error(err)
	}

	project, err := client.Project(projectId)

	if err != nil {
		t.Error(err)
	}

	if project.ID != projectId {
		t.Error("Project id doesn't match id passed in")
	}

	if project.Metadata.Public != false {
		t.Error("default public state not set on creation")
	}

	if project.Name != TestName {
		t.Errorf("Name was incorrectly set to: %s", project.Name)
	}

	if len(project.Components.ComponentIds) != 1 {
		t.Errorf("Wrong number of components returned expected 1, got %d", len(project.Components.ComponentIds))
	}
}
