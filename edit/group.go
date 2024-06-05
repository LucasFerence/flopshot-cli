package edit

const GroupType = "group"

type Group struct {
	Id               string        `json:"_id" mapstructure:"_id"`
	ExternalId       string        `json:"externalId"`
	Name             string        `json:"name"`
	ScheduleOffsetMs int           `json:"scheduleOffsetMs"`
	Provider         ReferenceType `json:"provider" editRef:"provider"`
}

func (group Group) Label() string {
	return group.Name
}

func init() {
	RegisterType(GroupType, Group{})
}
