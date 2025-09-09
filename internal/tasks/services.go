package tasks

import (
	"context"
	"coreflow/internal/db"

	"coreflow/gen/postgres/public/model"
	"coreflow/gen/postgres/public/table"

	"github.com/go-jet/jet/v2/postgres"
)

func NewServices(db *db.DB) *Services {
	return &Services{db: db}
}

type Services struct {
	db *db.DB
}

type TaskAssignee struct {
	ID    int32  `json:"id"`
	Email string `json:"email"`
}

type TaskAttachment struct {
	ID       int32  `json:"id"`
	Mimetype string `json:"mimetype"`
	Filename string `json:"filename"`
}

type TaskWithRelations struct {
	ID          int32            `json:"id"`
	Status      model.TaskStatus `json:"status"`
	Name        string           `json:"name"`
	Assignee    *TaskAssignee    `json:"assignee"`
	Attachments []TaskAttachment `json:"attachments"`
}

func (s *Services) ListTasks(ctx context.Context, cursor int32, limit int32) ([]TaskWithRelations, error) {
	t := table.Tasks.AS("t")
	u := table.Users.AS("u")
	ta := table.TasksAttachments.AS("ta")

	query := postgres.SELECT(
		t.ID.AS("task_id"),
		t.Status.AS("task_status"),
		t.Name.AS("task_name"),
		u.ID.AS("user_id"),
		u.Email.AS("user_email"),
		ta.ID.AS("attachment_id"),
		ta.Mimetype.AS("attachment_mimetype"),
		ta.Filename.AS("attachment_filename"),
	).FROM(
		t.LEFT_JOIN(u, u.ID.EQ(t.AssigneeID)).
			LEFT_JOIN(ta, ta.TaskID.EQ(t.ID)),
	).WHERE(
		t.ID.GT(postgres.Int(int64(cursor))),
	).ORDER_BY(
		t.ID, ta.ID,
	)

	sql, args := query.Sql()

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	taskMap := make(map[int32]*TaskWithRelations)
	taskOrder := []int32{}

	for rows.Next() {
		var (
			taskID         int32
			taskStatus     model.TaskStatus
			taskName       string
			userID         *int32
			userEmail      *string
			attachmentID   *int32
			attachmentMime *string
			attachmentFile *string
		)

		err := rows.Scan(
			&taskID, &taskStatus, &taskName,
			&userID, &userEmail,
			&attachmentID, &attachmentMime, &attachmentFile,
		)
		if err != nil {
			return nil, err
		}

		if len(taskMap) >= int(limit) && taskMap[taskID] == nil {
			break
		}

		task, exists := taskMap[taskID]
		if !exists {
			task = &TaskWithRelations{
				ID:          taskID,
				Status:      taskStatus,
				Name:        taskName,
				Attachments: []TaskAttachment{},
			}

			if userID != nil && userEmail != nil {
				task.Assignee = &TaskAssignee{
					ID:    *userID,
					Email: *userEmail,
				}
			}

			taskMap[taskID] = task
			taskOrder = append(taskOrder, taskID)
		}

		if attachmentID != nil && attachmentMime != nil && attachmentFile != nil {
			exists := false
			for _, existing := range task.Attachments {
				if existing.ID == *attachmentID {
					exists = true
					break
				}
			}

			if !exists {
				attachment := TaskAttachment{
					ID:       *attachmentID,
					Mimetype: *attachmentMime,
					Filename: *attachmentFile,
				}
				task.Attachments = append(task.Attachments, attachment)
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	result := make([]TaskWithRelations, 0, len(taskOrder))

	for _, id := range taskOrder {
		if task := taskMap[id]; task != nil {
			result = append(result, *task)
		}
	}

	return result, nil
}
