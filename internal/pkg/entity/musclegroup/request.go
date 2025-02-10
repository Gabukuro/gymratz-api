package musclegroup

import "github.com/Gabukuro/gymratz-api/internal/pkg/entity/base"

type (
	CreateMuscleGroupRequest struct {
		Name string `json:"name"`
	}

	UpdateMuscleGroupRequest struct {
		Name string `json:"name"`
	}

	ListMuscleGroupsQueryParams struct {
		base.ListQueryParams
		Name string `query:"name"`
	}
)
