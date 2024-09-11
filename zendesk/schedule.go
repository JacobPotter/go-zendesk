package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type ScheduleInterval struct {
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
}

type Schedule struct {
	Id        int                `json:"id"`
	Intervals []ScheduleInterval `json:"intervals,omitempty"`
	Name      string             `json:"name"`
	TimeZone  string             `json:"time_zone"`
	CreatedAt *time.Time         `json:"created_at,omitempty"`
	UpdatedAt *time.Time         `json:"updated_at,omitempty"`
}

type ScheduleAPI interface {
	CreateSchedule(ctx context.Context, schedule Schedule) (Schedule, error)
	GetSchedule(ctx context.Context, scheduleId int64) (Schedule, error)
	UpdateSchedule(ctx context.Context, scheduleId int64, schedule Schedule) (Schedule, error)
	DeleteSchedule(ctx context.Context, scheduleId int64) error
}

func (z *Client) CreateSchedule(ctx context.Context, schedule Schedule) (Schedule, error) {
	var data, result struct {
		Schedule Schedule `json:"schedule"`
	}

	data.Schedule = schedule

	body, err := z.post(ctx, "/schedules.json", data)
	if err != nil {
		return Schedule{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Schedule{}, err
	}

	return result.Schedule, nil

}

func (z *Client) GetSchedule(ctx context.Context, scheduleId int64) (Schedule, error) {
	var result struct {
		Schedule Schedule `json:"schedule"`
	}

	body, err := z.get(ctx, fmt.Sprintf("/schedules/%d.json", scheduleId))
	if err != nil {
		return Schedule{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Schedule{}, err
	}

	return result.Schedule, nil

}

func (z *Client) UpdateSchedule(ctx context.Context, scheduleId int64, schedule Schedule) (Schedule, error) {
	var data, result struct {
		Schedule Schedule `json:"schedule"`
	}

	data.Schedule = schedule

	body, err := z.put(ctx, fmt.Sprintf("/schedules/%d.json", scheduleId), data)
	if err != nil {
		return Schedule{}, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return Schedule{}, err
	}

	return result.Schedule, nil
}

func (z *Client) DeleteSchedule(ctx context.Context, scheduleId int64) error {
	err := z.delete(ctx, fmt.Sprintf("/schedules/%d.json", scheduleId))
	if err != nil {
		return err
	}
	return nil
}
