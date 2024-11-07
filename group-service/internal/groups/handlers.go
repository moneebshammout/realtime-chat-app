package groups

import (
	"group-service/pkg/utils"

	"github.com/labstack/echo/v4"

	DAO "group-service/internal/database/DAOs"
	"group-service/internal/database/models"
)

var logger = utils.GetLogger()

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
			Group:    group,
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
	group := models.Group{}
	DAO.Group.FindByID(id, &group, "GroupMembers")
	if group.ID == 0 {
		return c.JSON(404, "Group not found")
	}

	return c.JSON(200, group)
}

func deleteGroup(c echo.Context) error {
	return nil
}

func addMembers(c echo.Context) error {
	return nil
}

func removeMembers(c echo.Context) error {
	return nil
}
