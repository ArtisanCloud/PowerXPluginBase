package models

import (
	"fmt"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TaskActivity 任务活动记录模型
type TaskActivity struct {
	BaseModel

	// 关联信息
	TaskID uint  `gorm:"column:task_id;not null;index" json:"task_id"`
	UserID int64 `gorm:"column:user_id;not null;index" json:"user_id"`

	// 活动信息
	ActivityType ActivityType `gorm:"column:activity_type;type:varchar(20);not null;index" json:"activity_type"`
	Field        string       `gorm:"column:field;type:varchar(50)" json:"field,omitempty"`  // 变更的字段名
	OldValue     string       `gorm:"column:old_value;type:text" json:"old_value,omitempty"` // 旧值
	NewValue     string       `gorm:"column:new_value;type:text" json:"new_value,omitempty"` // 新值
	Comment      string       `gorm:"column:comment;type:text" json:"comment,omitempty"`     // 备注/评论

	// 时间记录
	TimeSpent    *int       `gorm:"column:time_spent" json:"time_spent,omitempty"`       // 花费时间（分钟）
	ActivityDate *time.Time `gorm:"column:activity_date" json:"activity_date,omitempty"` // 活动日期

	// 元数据
	Meta datatypes.JSON `gorm:"column:meta;type:jsonb" json:"meta,omitempty"`

	// 关联模型
	Task *Task `json:"task,omitempty" gorm:"foreignKey:TaskID"`
}

// ActivityType 活动类型枚举
type ActivityType string

const (
	ActivityTypeCreated      ActivityType = "created"       // 创建任务
	ActivityTypeUpdated      ActivityType = "updated"       // 更新任务
	ActivityTypeStatusChange ActivityType = "status_change" // 状态变更
	ActivityTypeAssigned     ActivityType = "assigned"      // 分配任务
	ActivityTypeUnassigned   ActivityType = "unassigned"    // 取消分配
	ActivityTypeCommented    ActivityType = "commented"     // 添加评论
	ActivityTypeTimeLogged   ActivityType = "time_logged"   // 记录工时
	ActivityTypeLabelAdded   ActivityType = "label_added"   // 添加标签
	ActivityTypeLabelRemoved ActivityType = "label_removed" // 移除标签
	ActivityTypeAttachment   ActivityType = "attachment"    // 添加附件
	ActivityTypeLinked       ActivityType = "linked"        // 关联任务
	ActivityTypeUnlinked     ActivityType = "unlinked"      // 取消关联
	ActivityTypeDeleted      ActivityType = "deleted"       // 删除任务
)

// TableName 返回表名
func (ta *TaskActivity) TableName() string {
	return PowerXSchema + "." + TableTaskActivity
}

// GetTableName 获取表名
func (ta *TaskActivity) GetTableName(needFull bool) string {
	if needFull {
		return ta.TableName()
	}
	return TableTaskActivity
}

// BeforeCreate GORM 钩子
func (ta *TaskActivity) BeforeCreate(tx *gorm.DB) error {
	// 设置默认活动日期
	if ta.ActivityDate == nil {
		now := time.Now()
		ta.ActivityDate = &now
	}

	return nil
}

// 业务方法

// IsStatusChange 检查是否为状态变更活动
func (ta *TaskActivity) IsStatusChange() bool {
	return ta.ActivityType == ActivityTypeStatusChange
}

// IsAssignment 检查是否为分配相关活动
func (ta *TaskActivity) IsAssignment() bool {
	return ta.ActivityType == ActivityTypeAssigned || ta.ActivityType == ActivityTypeUnassigned
}

// IsComment 检查是否为评论活动
func (ta *TaskActivity) IsComment() bool {
	return ta.ActivityType == ActivityTypeCommented
}

// IsTimeLogging 检查是否为工时记录活动
func (ta *TaskActivity) IsTimeLogging() bool {
	return ta.ActivityType == ActivityTypeTimeLogged
}

// GetDisplayMessage 获取活动显示消息
func (ta *TaskActivity) GetDisplayMessage() string {
	switch ta.ActivityType {
	case ActivityTypeCreated:
		return "创建了任务"
	case ActivityTypeUpdated:
		if ta.Field != "" {
			return fmt.Sprintf("更新了 %s", ta.Field)
		}
		return "更新了任务"
	case ActivityTypeStatusChange:
		if ta.OldValue != "" && ta.NewValue != "" {
			return fmt.Sprintf("状态从 %s 变更为 %s", ta.OldValue, ta.NewValue)
		}
		return "变更了状态"
	case ActivityTypeAssigned:
		return "分配了任务"
	case ActivityTypeUnassigned:
		return "取消了任务分配"
	case ActivityTypeCommented:
		return "添加了评论"
	case ActivityTypeTimeLogged:
		if ta.TimeSpent != nil {
			hours := *ta.TimeSpent / 60
			minutes := *ta.TimeSpent % 60
			if hours > 0 {
				return fmt.Sprintf("记录了 %d 小时 %d 分钟工时", hours, minutes)
			}
			return fmt.Sprintf("记录了 %d 分钟工时", minutes)
		}
		return "记录了工时"
	case ActivityTypeLabelAdded:
		if ta.NewValue != "" {
			return fmt.Sprintf("添加了标签 %s", ta.NewValue)
		}
		return "添加了标签"
	case ActivityTypeLabelRemoved:
		if ta.OldValue != "" {
			return fmt.Sprintf("移除了标签 %s", ta.OldValue)
		}
		return "移除了标签"
	case ActivityTypeAttachment:
		return "添加了附件"
	case ActivityTypeLinked:
		return "关联了任务"
	case ActivityTypeUnlinked:
		return "取消了任务关联"
	case ActivityTypeDeleted:
		return "删除了任务"
	default:
		return string(ta.ActivityType)
	}
}

// GetTimeSpentFormatted 获取格式化的时间花费
func (ta *TaskActivity) GetTimeSpentFormatted() string {
	if ta.TimeSpent == nil {
		return ""
	}

	hours := *ta.TimeSpent / 60
	minutes := *ta.TimeSpent % 60

	if hours > 0 {
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dm", minutes)
}

// IsRecent 检查活动是否为最近发生（24小时内）
func (ta *TaskActivity) IsRecent() bool {
	if ta.ActivityDate == nil {
		return false
	}
	return time.Since(*ta.ActivityDate) <= 24*time.Hour
}

// Validate 验证任务活动数据
func (ta *TaskActivity) Validate() error {
	if ta.TaskID == 0 {
		return fmt.Errorf("task_id is required")
	}

	if ta.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}

	if !ta.ActivityType.IsValid() {
		return fmt.Errorf("invalid activity type: %s", ta.ActivityType)
	}

	if ta.TimeSpent != nil && *ta.TimeSpent < 0 {
		return fmt.Errorf("time spent cannot be negative")
	}

	// 工时记录必须有 TimeSpent
	if ta.ActivityType == ActivityTypeTimeLogged && (ta.TimeSpent == nil || *ta.TimeSpent == 0) {
		return fmt.Errorf("time logged activity must have positive time spent")
	}

	// 评论活动必须有评论内容
	if ta.ActivityType == ActivityTypeCommented && ta.Comment == "" {
		return fmt.Errorf("comment activity must have comment content")
	}

	return nil
}

// 静态方法

// NewTaskCreatedActivity 创建任务创建活动
func NewTaskCreatedActivity(taskID uint, userID int64) *TaskActivity {
	return &TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: ActivityTypeCreated,
	}
}

// NewStatusChangeActivity 创建状态变更活动
func NewStatusChangeActivity(taskID uint, userID int64, oldStatus, newStatus TaskStatus) *TaskActivity {
	return &TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: ActivityTypeStatusChange,
		Field:        "status",
		OldValue:     string(oldStatus),
		NewValue:     string(newStatus),
	}
}

// NewAssignmentActivity 创建分配活动
func NewAssignmentActivity(taskID uint, userID int64, assigneeID *int64) *TaskActivity {
	activity := &TaskActivity{
		TaskID: taskID,
		UserID: userID,
		Field:  "assignee",
	}

	if assigneeID != nil {
		activity.ActivityType = ActivityTypeAssigned
		activity.NewValue = fmt.Sprintf("%d", *assigneeID)
	} else {
		activity.ActivityType = ActivityTypeUnassigned
	}

	return activity
}

// NewCommentActivity 创建评论活动
func NewCommentActivity(taskID uint, userID int64, comment string) *TaskActivity {
	return &TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: ActivityTypeCommented,
		Comment:      comment,
	}
}

// NewTimeLogActivity 创建工时记录活动
func NewTimeLogActivity(taskID uint, userID int64, timeSpent int, comment string) *TaskActivity {
	return &TaskActivity{
		TaskID:       taskID,
		UserID:       userID,
		ActivityType: ActivityTypeTimeLogged,
		TimeSpent:    &timeSpent,
		Comment:      comment,
	}
}

// NewLabelActivity 创建标签活动
func NewLabelActivity(taskID uint, userID int64, label string, added bool) *TaskActivity {
	activity := &TaskActivity{
		TaskID: taskID,
		UserID: userID,
		Field:  "labels",
	}

	if added {
		activity.ActivityType = ActivityTypeLabelAdded
		activity.NewValue = label
	} else {
		activity.ActivityType = ActivityTypeLabelRemoved
		activity.OldValue = label
	}

	return activity
}

// 枚举验证方法

// IsValid 验证活动类型是否有效
func (at ActivityType) IsValid() bool {
	switch at {
	case ActivityTypeCreated, ActivityTypeUpdated, ActivityTypeStatusChange,
		ActivityTypeAssigned, ActivityTypeUnassigned, ActivityTypeCommented,
		ActivityTypeTimeLogged, ActivityTypeLabelAdded, ActivityTypeLabelRemoved,
		ActivityTypeAttachment, ActivityTypeLinked, ActivityTypeUnlinked, ActivityTypeDeleted:
		return true
	default:
		return false
	}
}

// String 实现 Stringer 接口
func (at ActivityType) String() string {
	return string(at)
}
