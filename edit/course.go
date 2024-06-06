package edit

const CourseType = "course"

type Course struct {
	Id         string        `json:"_id" mapstructure:"_id"`
	ExternalId string        `json:"externalId"`
	Name       string        `json:"name"`
	Provider   ReferenceType `json:"provider" editRef:"provider"`
	Group      ReferenceType `json:"group" editRef:"group"`
}

func (course Course) Label() string {
	return course.Name
}

func init() {
	RegisterType(CourseType, Course{})
}
