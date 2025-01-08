package gRPC

import (
	"context"
	"fmt"
	"strconv"

	"group-service/pkg/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	DAO "group-service/internal/database/DAOs"
	groupGen "group-service/internal/gRPC/group-service-grpc-gen"
)

var logger = utils.GetLogger()

type GroupServiceServer struct {
	groupGen.UnimplementedGroupServiceServer
}

func (s GroupServiceServer) GetGroupUsers(ctx context.Context, req *groupGen.GetGroupUsersRequest) (*groupGen.GetGroupUsersResponse, error) {
	logger.Infof("GetGroupUsers: %s", req.GroupId)

	groupId, err := strconv.Atoi(req.GroupId)
	if err != nil {
		logger.Errorf("Group not found: %v\n", err)
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Group not found: %s", err.Error()),
		)
	}

	group, err := DAO.Group.FindByID(groupId, "GroupMembers")
	if err != nil {
		logger.Errorf("Group not found: %v\n", err)
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Group not found: %s", err.Error()),
		)
	}

	userIds := []string{}

	for _, member := range group.GroupMembers {
		userIds = append(userIds, member.SourceId)
	}

	return &groupGen.GetGroupUsersResponse{
		Status:  "OK",
		UserIds: userIds,
	}, nil
}
