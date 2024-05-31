package edit

const GroupType = "group"

type Group struct {
	Id               string        `json:"_id" mapstructure:"_id"`
	Name             string        `json:"name"`
	Provider         ReferenceType `json:"provider"`
	ScheduleOffsetMs int           `json:"scheduleOffsetMs"`
}

func (group Group) Label() string {
	return group.Name
}

func init() {
	RegisterType(GroupType, Group{})
}
