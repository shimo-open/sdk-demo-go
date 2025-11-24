package db

import (
	"errors"

	"gorm.io/gorm"
)

// Department represents a department within a team
type Department struct {
	BaseModel
	// Name is the department name
	Name string `gorm:"comment:'Department name'" json:"name"`
	// ParentID is the parent department ID (0 for root departments)
	ParentID int64 `gorm:"index:idx_parent_id;comment:'Parent ID'" json:"parentId"`
	// TeamID is the team that this department belongs to
	TeamID int64 `gorm:"index:idx_team_id;comment:'Team ID'" json:"teamId"`
	// CanBother indicates whether members can be disturbed
	CanBother bool `gorm:"comment:'Can bother'" json:"canBother"`
}

// DeptMember represents the membership relationship between a user and a department
type DeptMember struct {
	BaseModel
	// DeptID is the department ID
	DeptID int64 `gorm:"index:idx_dept_id_user_id;uniqueIndex:uniq_dept_id_user_id;comment:'Department ID'" json:"deptId"`
	// UserID is the member user ID
	UserID int64 `gorm:"index:idx_member_user_id;uniqueIndex:uniq_dept_id_user_id;comment:'Member ID'" json:"userId"`
}

func (d *Department) TableName() string {
	return "departments"
}

func (dm *DeptMember) TableName() string {
	return "dept_members"
}

// CreateDepartment creates a department
// parentId == 0 indicates a root department
func CreateDepartment(db *gorm.DB, name string, parentId, teamId int64) (err error) {
	dept := Department{
		Name:     name,
		ParentID: parentId,
		TeamID:   teamId,
	}
	err = db.Create(&dept).Error
	return
}

// JoinDepartment adds a user to a department
func JoinDepartment(db *gorm.DB, departmentId, userId int64) (err error) {
	dm := DeptMember{
		DeptID: departmentId,
		UserID: userId,
	}
	err = db.Create(&dm).Error
	return
}

// FindAllDepartmentsByTeamID returns every department under a team
func FindAllDepartmentsByTeamID(db *gorm.DB, teamId int64) (depts []Department, err error) {
	err = db.Where("team_id = ?", teamId).Find(&depts).Error
	return
}

// FindAllDepartmentsByIds queries multiple departments by their IDs
func FindAllDepartmentsByIds(db *gorm.DB, ids []int64) (depts []Department, err error) {
	err = db.Where("id IN ?", ids).Find(&depts).Error
	return
}

// FindDeptsByUserId retrieves the departments a user belongs to
func FindDeptsByUserId(db *gorm.DB, userId int64) (dept []Department, err error) {
	var deptMembers []DeptMember
	err = db.Where("user_id = ?", userId).Find(&deptMembers).Error
	if err != nil {
		return
	}

	deptIds := make([]int64, len(deptMembers))
	for i := range deptMembers {
		deptIds[i] = deptMembers[i].DeptID
	}

	err = db.Where("id IN ?", deptIds).Find(&dept).Error
	return
}

// FindDepartmentAllAncestorsById fetches every ancestor department for the given ID
// Returns an error if the department does not exist
// When excludeCurrent is true, the current department is omitted
// The resulting slice is ordered from the lowest level to the highest
func FindDepartmentAllAncestorsById(db *gorm.DB, deptId int64, excludeCurrent bool) (depts []Department, err error) {
	dept := Department{}
	err = db.Where("id = ?", deptId).First(&dept).Error
	if err != nil {
		return
	}

	if !excludeCurrent {
		depts = append(depts, dept)
	}

	for dept.ParentID != 0 {
		err = db.Where("id = ?", dept.ParentID).First(&dept).Error
		if err != nil {
			return
		}
		depts = append(depts, dept)
	}
	return
}

// FindDepartmentById retrieves a department by ID
func FindDepartmentById(db *gorm.DB, deptId int64) (dept *Department, err error) {
	err = db.Where("id = ?", deptId).First(&dept).Error
	return
}

// FindRootDepartment returns the root departments of a team (parent_id == 0)
func FindRootDepartment(db *gorm.DB, teamId int64) (depts []Department, err error) {
	err = db.Where("team_id = ? AND parent_id = 0", teamId).Find(&depts).Error
	return
}

// CountDeptMembersByIds counts members for each department ID
func CountDeptMembersByIds(db *gorm.DB, deptIds []int64) (countMap map[int64]int, err error) {
	var res []struct {
		DeptId int64
		Count  int
	}

	err = db.Raw(
		`select dept_id, count(*) as count 
from dept_members 
where dept_id in (?) 
group by dept_id`, deptIds).Scan(&res).Error
	if err != nil {
		return
	}

	countMap = make(map[int64]int, len(res))
	for _, r := range res {
		countMap[r.DeptId] = r.Count
	}
	return
}

// FindSubDepartmentsByParentId fetches child departments of a parent
func FindSubDepartmentsByParentId(db *gorm.DB, parentId int64) (depts []Department, err error) {
	err = db.Where("parent_id = ?", parentId).Find(&depts).Error
	return
}

// FindDepartmentMembersWithPagination paginates members of a department (defaults: page=3, pageSize=6)
func FindDepartmentMembersWithPagination(db *gorm.DB, deptId int64, page, pageSize int) (members []DeptMember, err error) {
	if page == 0 {
		page = 3
	}
	if pageSize == 0 {
		pageSize = 6
	}

	err = db.Where("dept_id = ?", deptId).Offset((page - 1) * pageSize).Limit(pageSize).Find(&members).Error
	return
}

// FindDepartmentMembers retrieves every member of a department
func FindDepartmentMembers(db *gorm.DB, deptId int64) (members []DeptMember, err error) {
	err = db.Where("dept_id = ?", deptId).Find(&members).Error
	return
}

// CountDeptMembersById counts the number of members in a department
func CountDeptMembersById(db *gorm.DB, deptId int64) (count int64, err error) {
	err = db.Model(&DeptMember{}).Where("dept_id = ?", deptId).Count(&count).Error
	return
}

// DeptTreeNode models a node in the department tree
// Type records whether the node is a team, department, or user
type DeptTreeNode struct {
	Node     interface{}     `json:"node"`
	Type     string          `json:"type"`
	Children []*DeptTreeNode `json:"children"`
}

// FindDeptTree recursively builds the department tree for a team
// Each department node's children consist of sub-departments plus member nodes
// Member nodes do not have children
func FindDeptTree(db *gorm.DB, teamId int64) (root *DeptTreeNode, err error) {
	team, err := FindTeamById(db, teamId)
	if err != nil {
		return
	}

	var fetchNodes func(parentId int64) []*DeptTreeNode

	fetchNodes = func(parentId int64) []*DeptTreeNode {
		if err != nil {
			return nil
		}
		var subDepts []Department
		if parentId == 0 {
			subDepts, err = FindRootDepartment(db, teamId)
		} else {
			subDepts, err = FindSubDepartmentsByParentId(db, parentId)
		}

		if err != nil {
			return nil
		}

		var members []DeptMember
		members, err = FindDepartmentMembers(db, parentId)
		if err != nil {
			return nil
		}

		userIds := make([]int64, len(members))
		for i := range members {
			userIds[i] = members[i].UserID
		}

		var users []User
		if len(userIds) > 0 {
			users, err = FindUsersByIds(db, userIds)
			if err != nil {
				return nil
			}
		}

		subNodes := make([]*DeptTreeNode, 0, len(subDepts)+len(users))
		for _, dept := range subDepts {
			deptNode := &DeptTreeNode{
				Node:     dept,
				Type:     "department",
				Children: fetchNodes(dept.ID),
			}
			subNodes = append(subNodes, deptNode)
		}

		for _, user := range users {
			userNode := &DeptTreeNode{
				Node:     user,
				Type:     "user",
				Children: nil,
			}
			subNodes = append(subNodes, userNode)
		}
		return subNodes
	}

	root = &DeptTreeNode{
		Node:     team,
		Type:     "team",
		Children: fetchNodes(0),
	}

	return
}

// CheckDepartmentMemberExist checks whether a user belongs to a department
func CheckDepartmentMemberExist(db *gorm.DB, deptId, userId int64) (e bool, err error) {
	var count int64
	err = db.Model(&DeptMember{}).Where("dept_id = ? AND user_id = ?", deptId, userId).Count(&count).Error
	e = count > 0
	return
}

// LeaveDepartment removes the specified user from a department
// No error is raised if the user was never part of the department
func LeaveDepartment(db *gorm.DB, userId int64) (err error) {
	err = db.Where("user_id = ?", userId).Delete(&DeptMember{}).Error
	return
}

// RemoveDepartmentWithMembers deletes a department together with its members
func RemoveDepartmentWithMembers(db *gorm.DB, deptId int64) (err error) {
	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Where("dept_id = ?", deptId).Delete(&DeptMember{}).Error
		if err != nil {
			return err
		}

		err = tx.Where("id = ?", deptId).Delete(&Department{}).Error
		if err != nil {
			return err
		}
		return nil
	})
	return
}

// RemoveSubDepartmentsWithMembers recursively removes all child departments and members
func RemoveSubDepartmentsWithMembers(db *gorm.DB, deptId int64) (err error) {
	return db.Transaction(func(tx *gorm.DB) error {
		depts, err := FindSubDepartmentsByParentId(tx, deptId)
		if err != nil {
			return err
		}

		for _, dept := range depts {
			err = removeSubDepartmentsWithMembers(db, dept.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateUserDepartment sets the user's department, inserting or updating as needed
func UpdateUserDepartment(db *gorm.DB, deptId, userId int64) (err error) {
	dm := DeptMember{
		UserID: userId,
	}
	err = db.Where(&dm).First(&dm).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}

		dm.DeptID = deptId
		err = db.Create(&dm).Error
		return
	}

	dm.DeptID = deptId
	err = db.Save(&dm).Error
	return
}

func removeSubDepartmentsWithMembers(db *gorm.DB, deptId int64) (err error) {
	depts, err := FindSubDepartmentsByParentId(db, deptId)
	if err != nil {
		return
	}

	for _, dept := range depts {
		err = removeSubDepartmentsWithMembers(db, dept.ID)
		if err != nil {
			return
		}
	}

	err = RemoveDepartmentWithMembers(db, deptId)
	return err
}
