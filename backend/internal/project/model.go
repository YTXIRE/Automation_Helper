package project

type Project struct {
	ID                string `json:"id" bson:"_id"`
	Name              string `json:"name" bson:"name"`
	TestRailProjectID string `json:"testrail_project_id" bson:"testrail_project_id"`
}
