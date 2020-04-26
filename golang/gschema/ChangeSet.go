package gschema

type ChangeType int

const (
	Attribute_Added  ChangeType = 1
	Attribute_Delete ChangeType = 2
	Attribute_Change ChangeType = 3
)

type ChangeSet struct {
	changedInstance interface{}
	changes         []*ChangeAttribute
}

type ChangeAttribute struct {
	changeType ChangeType
	attribute  *Attribute
	oldValue   interface{}
	newValue   interface{}
}

func (changeSet *ChangeSet) ChangedInstance() interface{} {
	return changeSet.changedInstance
}

func (changeSet *ChangeSet) Changes() []*ChangeAttribute {
	return changeSet.changes
}

func NewChangeSet(changedInstance interface{}) *ChangeSet {
	changeSet := &ChangeSet{}
	changeSet.changedInstance = changedInstance
	changeSet.changes = make([]*ChangeAttribute, 0)
	return changeSet
}

func (n *ChangeSet) AddChangeAttribute(changeType ChangeType, attribute *Attribute, oldValue interface{}, newValue interface{}) {
	changeAttribute := &ChangeAttribute{}
	changeAttribute.attribute = attribute
	changeAttribute.oldValue = oldValue
	changeAttribute.newValue = newValue
	changeAttribute.changeType = changeType
	n.changes = append(n.changes, changeAttribute)
}

func (pn *ChangeAttribute) Type() ChangeType {
	return pn.changeType
}

func (pn *ChangeAttribute) Attribute() *Attribute {
	return pn.attribute
}

func (pn *ChangeAttribute) OldValue() interface{} {
	return pn.oldValue
}

func (pn *ChangeAttribute) NewValue() interface{} {
	return pn.newValue
}
