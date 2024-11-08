package groups

import (
	"group-service/pkg/utils"

	"github.com/labstack/echo/v4"

	DAO "group-service/internal/database/DAOs"
	"group-service/internal/database/models"
)

var logger = utils.GetLogger()

func list(c echo.Context) error {
	data,err := DAO.Group.FindAll()
	if err != nil {
		logger.Errorf("Error creating group: %v", err)
		return c.JSON(500, "Error creating group")
	}

	return c.JSON(200, data)
}

func create(c echo.Context) error {
	data := c.Get("validatedData").(*GroupCreateSerizliser)
	group := models.Group{Name: data.Name}
	err := DAO.Group.Create(&group)
	if err != nil {
		logger.Errorf("Error creating group: %v", err)
		return c.JSON(500, "Error creating group")
	}

	groupMembers := []models.GroupMember{}
	for _, memeber := range data.Members {
		if !models.MemberRoles[memeber.Role] {
			return c.JSON(400, "Invalid role")
		}

		groupMembers = append(groupMembers, models.GroupMember{
			SourceId: memeber.ID,
			Role:     memeber.Role,
			GroupId:  group.ID,
		})
	}

	err = DAO.GroupMember.Create(&groupMembers)
	if err != nil {
		logger.Errorf("Error creating group members: %v", err)
		return c.JSON(500, "Error creating group members")
	}

	return c.JSON(201, map[string]any{
		"group":   group,
		"members": groupMembers,
	})
}

func getGroup(c echo.Context) error {
	id := c.Get("validatedData").(*GroupGetSerizliser).ID
	group, err := DAO.Group.FindByID(id, "GroupMembers")
	if err != nil {
		return c.JSON(404, "Group not found")
	}

	return c.JSON(200, group)
}

func addMembers(c echo.Context) error {
	data := c.Get("validatedData").(*AddMembersSerizliser)
	groupMembers := []models.GroupMember{}
	for _, memeber := range data.Members {
		if !models.MemberRoles[memeber.Role] {
			return c.JSON(400, "Invalid role")
		}

		groupMembers = append(groupMembers, models.GroupMember{
			SourceId: memeber.ID,
			Role:     memeber.Role,
			GroupId:  data.ID,
		})
	}

	err := DAO.GroupMember.Create(&groupMembers)
	if err != nil {
		logger.Errorf("Error creating group members: %v", err)
		return c.JSON(500, "Error creating group members")
	}

	return c.JSON(201, groupMembers)
}

func removeMembers(c echo.Context) error {
	data := c.Get("validatedData").(*removeMembersSerizliser)
	err := DAO.GroupMember.Delete(data.IDS...)
	if err != nil {
		logger.Errorf("Error removing group members: %v", err)
		return c.JSON(500, "Error removing group members")
	}

	return c.JSON(200, "Members removed")
}

func deleteGroup(c echo.Context) error {
	data := c.Get("validatedData").(*IDParam)

	members, err := DAO.GroupMember.FindAll(map[string]interface{}{
		"group_id": data.ID,
	})
	if err != nil {
		logger.Errorf("Error removing group : %v", err)
		return c.JSON(500, "Error removing group ")
	}

	memberIds := []uint{}
	for _, member := range members {
		memberIds = append(memberIds, member.ID)
	}

	err = DAO.GroupMember.Delete(memberIds...)
	if err != nil {
		logger.Errorf("Error removing group : %v", err)
		return c.JSON(500, "Error removing group ")
	}

	err = DAO.Group.Delete(data.ID)
	if err != nil {
		logger.Errorf("Error removing group : %v", err)
		return c.JSON(500, "Error removing group ")
	}

	return c.JSON(200, "group removed")
}
