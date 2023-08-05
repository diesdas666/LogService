package mapper

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ModelIdToRepoId(ID string) (primitive.ObjectID, error) {
	repoID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("failed to parse id: %s", ID)
	}
	return repoID, nil
}

func RepoIdToModelId(ID primitive.ObjectID) string {
	return ID.Hex()
}
