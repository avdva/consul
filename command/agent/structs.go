package agent

import (
	"github.com/hashicorp/consul/consul/structs"
	"github.com/hashicorp/consul/types"
)

// ServiceDefinition is used to JSON decode the Service definitions
type ServiceDefinition struct {
	ID                string
	Name              string
	Tags              []string
	DynamicTags       DynamicTags
	Address           string
	Port              int
	Check             CheckType
	Checks            CheckTypes
	Token             string
	EnableTagOverride bool
}

func (s *ServiceDefinition) NodeService() *structs.NodeService {
	ns := &structs.NodeService{
		ID:                s.ID,
		Service:           s.Name,
		Tags:              s.Tags,
		Address:           s.Address,
		Port:              s.Port,
		EnableTagOverride: s.EnableTagOverride,
	}
	if ns.ID == "" && ns.Service != "" {
		ns.ID = ns.Service
	}
	return ns
}

func (s *ServiceDefinition) CheckTypes() (checks CheckTypes) {
	if s.Check.Valid() {
		checks = append(checks, &s.Check)
	}
	for _, check := range s.Checks {
		if check.Valid() {
			checks = append(checks, check)
		}
	}
	return
}

func (s *ServiceDefinition) TagsCheckTypes() (tags DynamicTags) {
	for _, check := range s.DynamicTags {
		if check.Valid() {
			tags = append(tags, check)
		}
	}
	return
}

// DynamicTag is used to JSON decode the Dynamic tag definitions
type DynamicTag struct {
	Name      string
	CheckType `mapstructure:",squash"`
}

type DynamicTags []*DynamicTag

func (d DynamicTags) Names() (names []string) {
	for _, tag := range d {
		names = append(names, tag.Name)
	}
	return
}

func (d DynamicTags) CheckTypes() (checks CheckTypes) {
	for _, tag := range d {
		checks = append(checks, &tag.CheckType)
	}
	return
}

// CheckDefinition is used to JSON decode the Check definitions
type CheckDefinition struct {
	ID        types.CheckID
	Name      string
	Notes     string
	ServiceID string
	Token     string
	Status    string
	CheckType `mapstructure:",squash"`
}

func (c *CheckDefinition) HealthCheck(node string) *structs.HealthCheck {
	health := &structs.HealthCheck{
		Node:      node,
		CheckID:   c.ID,
		Name:      c.Name,
		Status:    structs.HealthCritical,
		Notes:     c.Notes,
		ServiceID: c.ServiceID,
	}
	if c.Status != "" {
		health.Status = c.Status
	}
	if health.CheckID == "" && health.Name != "" {
		health.CheckID = types.CheckID(health.Name)
	}
	return health
}

// persistedService is used to wrap a service definition and bundle it
// with an ACL token so we can restore both at a later agent start.
type persistedService struct {
	Token   string
	Service *structs.NodeService
}
