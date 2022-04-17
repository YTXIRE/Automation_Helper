package project

type Project struct {
	ID                string `json:"id" bson:"_id,omitempty"`
	Name              string `json:"name" bson:"name"`
	TestRailProjectID string `json:"testrail_project_id" bson:"testrail_project_id"`
	CreatedAt         int64  `json:"created_at" bson:"created_at"`
	UpdatedAt         int64  `json:"updated_at" bson:"updated_at"`
}
