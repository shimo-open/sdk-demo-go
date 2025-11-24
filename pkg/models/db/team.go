package db

import (
	"errors"

	"gorm.io/gorm"
)

// Team represents a team in the system
type Team struct {
	BaseModel
	// Name is the team name
	Name string `gorm:"comment:'Team name'" json:"name"`
}

// TeamRole represents a user's role in a team
type TeamRole struct {
	BaseModel
	// TeamID is the team ID
	TeamID int64 `gorm:"uniqueIndex:uniq_team_id_user_id;comment:'Team ID'" json:"teamId"`
	// UserID is the user ID
	UserID int64 `gorm:"uniqueIndex:uniq_team_id_user_id;index:idx_team_user_id;comment:'User ID'" json:"userId"`
	// Role is the user's role in the team (creator/manager/member)
	Role string `gorm:"check:role IN ('creator','manager','member');comment:'Role (creator/manager/member)'"  json:"role"`
}

const (
	CREATOR = "creator"
	MANAGER = "manager"
	MEMBER  = "member"
)

func (t *Team) TableName() string {
	return "teams"
}

func (tr *TeamRole) TableName() string {
	return "team_role"
}

// CreateTeam creates a team from the instance and creator ID
func CreateTeam(db *gorm.DB, team *Team, userId int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(team).Error
		if err != nil {
			return err
		}

		err = tx.Create(&TeamRole{
			TeamID: team.ID,
			UserID: userId,
			Role:   CREATOR,
		}).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// FindTeamById retrieves a team by ID
func FindTeamById(db *gorm.DB, teamId int64) (team *Team, err error) {
	err = db.Where("id = ?", teamId).First(&team).Error
	return
}

// JoinTeam adds a user to a team.
// Returns an error if the specified team ID does not exist
func JoinTeam(db *gorm.DB, teamId int64, userId int64) (err error) {
	_, err = FindTeamById(db, teamId)
	if err != nil {
		return
	}

	tr := &TeamRole{
		TeamID: teamId,
		UserID: userId,
		Role:   MEMBER,
	}
	err = db.Create(tr).Error
	return
}

// LeaveTeam soft-deletes the user-team relation (i.e., leaves the team).
// Returns an error if the user is not part of the team
// Returns an error if the user is the team creator
func LeaveTeam(db *gorm.DB, teamId int64, userId int64) (err error) {
	tr := &TeamRole{
		TeamID: teamId,
		UserID: userId,
	}

	check := &TeamRole{}
	err = db.Where(tr).First(check).Error
	if err != nil {
		return
	}

	if check.Role == CREATOR {
		return errors.New("creator can't leave team")
	}

	err = db.Where(tr).Delete(&TeamRole{}).Error
	return
}

// FindAllTeams returns every team
func FindAllTeams(db *gorm.DB) (teams []Team, err error) {
	err = db.Find(&teams).Error
	return
}

// FindTeamCreator fetches the creator ID for a team
func FindTeamCreator(db *gorm.DB, teamId int64) (userId int64, err error) {
	tr := &TeamRole{
		TeamID: teamId,
		Role:   CREATOR,
	}
	err = db.Where(tr).First(tr).Error
	return tr.UserID, err
}

// TransferTeam transfers the team ownership
// Returns an error if the team does not exist
// Returns an error if the new owner is not part of the team
func TransferTeam(db *gorm.DB, teamId int64, newCreatorId int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		tr := &TeamRole{
			TeamID: teamId,
			Role:   CREATOR,
		}
		res := tx.Model(&tr).Where(&tr).Update("role", MEMBER)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return errors.New("specified team not exists, check team id")
		}

		ntr := &TeamRole{
			TeamID: teamId,
			UserID: newCreatorId,
		}
		res = tx.Model(&ntr).Where(&ntr).Update("role", CREATOR)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return errors.New("no such user in the team, check new creator id")
		}
		return nil
	})
}

type TeamWithCreator struct {
	Team
	CreatorID int64    `json:"creatorId"`
	Creator   User     `json:"creator"`
	Role      TeamRole `json:"role"`
}

// FindTeamsByUserId gets every team for the given user ID
// Each returned team includes detail fields such as creatorId, creator, and role
func FindTeamsByUserId(db *gorm.DB, userId int64) (twc []TeamWithCreator, err error) {
	var results []struct {
		Team
		TeamRole
		User
	}

	err = db.Raw(`
select teams.*, team_role.*, users.*
from teams
join team_role on teams.id = team_role.team_id
join users on team_role.user_id = users.id
where team_role.deleted_at = 0 and (team_role.user_id = ? or team_role.role = 'creator')
`, userId).Scan(&results).Error

	creators := map[int64]User{}
	for i := range results {
		if results[i].Role != CREATOR {
			continue
		}
		creators[results[i].Team.ID] = results[i].User
	}

	twc = make([]TeamWithCreator, 0)

	for i := range results {
		if results[i].UserID != userId {
			continue
		}
		var t TeamWithCreator
		t.CreatorID = creators[results[i].TeamID].ID
		t.Creator = creators[results[i].TeamID]
		t.Team = results[i].Team
		t.Role = results[i].TeamRole
		twc = append(twc, t)
	}
	return
}

// FindTeamAllMembersByTeamId retrieves every member ID under a team
func FindTeamAllMembersByTeamId(db *gorm.DB, teamId int64) (userIds []int64, err error) {
	var trs []TeamRole
	err = db.Where("team_id = ?", teamId).Find(&trs).Error
	if err != nil {
		return
	}

	userIds = make([]int64, len(trs))
	for i := range trs {
		userIds[i] = trs[i].UserID
	}
	return
}

// FindTeamRole fetches a member's role within a team
func FindTeamRole(db *gorm.DB, teamId int64, userId int64) (tr *TeamRole, err error) {
	err = db.Where("team_id = ? AND user_id = ?", teamId, userId).First(&tr).Error
	return
}

// FindTeamRoleWithLimit lists team members with a limit (default 10)
func FindTeamRoleWithLimit(db *gorm.DB, teamId int64, limit int) (tr *TeamRole, err error) {
	if limit == 0 {
		limit = 10
	}
	err = db.Where("team_id = ?", teamId).Limit(limit).Find(&tr).Error
	return
}

// FindTeamWithPagination lists team members with pagination (default: 3 pages, 6 per page)
func FindTeamWithPagination(db *gorm.DB, teamId int64, page int, pageSize int) (trs []TeamRole, err error) {
	if page == 0 {
		page = 3
	}
	if pageSize == 0 {
		pageSize = 6
	}
	err = db.Where("team_id = ?", teamId).Offset((page - 1) * pageSize).Limit(pageSize).Find(&trs).Error
	return
}

// FindTeamRoles fetches roles for multiple user IDs
func FindTeamRoles(db *gorm.DB, teamId int64, userIds []int64) (trs []TeamRole, err error) {
	err = db.Where("user_id IN ?", userIds).Find(&trs).Error
	return
}

// CountTeamMembers counts how many members a team has
func CountTeamMembers(db *gorm.DB, teamId int64) (cnt int64, err error) {
	err = db.Model(&TeamRole{}).Where("team_id", teamId).Count(&cnt).Error
	return
}
