package v1

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/WishZ/go-todo-service/pkg/enum"
	"time"

	v1 "github.com/WishZ/go-todo-service/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// toDoServiceServer是v1.ToDoServiceServer proto接口的实现
type toDoServiceServer struct {
	db *sql.DB
}

func (s toDoServiceServer) Create(ctx context.Context, request *v1.CreateRequest) (*v1.CreateResponse, error) {
	//检查版本
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}

	//获取db连接
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(request.ToDo.Reminder)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	//插入数据
	res, err := c.ExecContext(ctx, "INSERT INTO m_todo(`title`,`description`,`reminder`) VALUES (?,?,?)", request.ToDo.Title, request.ToDo.Description, reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to insert into ToDo-> "+err.Error())
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve id for created ToDo-> "+err.Error())
	}
	return &v1.CreateResponse{
		Api: enum.ApiVersion,
		Id:  id,
	}, nil
}

func (s toDoServiceServer) Read(ctx context.Context, request *v1.ReadRequest) (*v1.ReadResponse, error) {
	//检查版本
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}

	//获取db连接
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "SELECT id,title,description,reminder FROM m_todo WHERE `id` = ?", request.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from m_todo->"+err.Error())
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("m_todo with id = '%d' is not found", request.Id))
	}

	var td v1.ToDo
	var reminder time.Time

	if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "failed tod retrieve data from ToDo row =>"+err.Error())
	}
	td.Reminder, err = ptypes.TimestampProto(reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
	}
	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("found multiple ToDo rows with ID='%d'", request.Id))
	}

	return &v1.ReadResponse{
		Api:  enum.ApiVersion,
		Todo: &td,
	}, nil
}

func (s toDoServiceServer) Update(ctx context.Context, request *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	// 检查服务器是否支持客户端请求的API版本
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}
	// 从池中获取sql连接
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(request.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "reminder field has invalid format-> "+err.Error())
	}

	res, err := c.ExecContext(ctx, "UPDATE m_todo SET `title`=?,`Description` = ?,`reminder` = ? WHERE `id` = ?",
		request.ToDo.Title, request.ToDo.Description, reminder, request.ToDo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "数据更新失败："+err.Error())
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "获取更新行数失败："+err.Error())
	}

	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("更新的ID:%d不存在", request.ToDo.Id))
	}

	return &v1.UpdateResponse{
		Api:     enum.ApiVersion,
		Updated: rows,
	}, nil
}

func (s toDoServiceServer) Delete(ctx context.Context, request *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	// 检查服务器是否支持客户端请求的API版本
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}
	// 从池中获取sql连接
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	// 删除ToDo
	res, err := c.ExecContext(ctx, "DELETE FROM m_todo WHERE `id`=?", request.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to delete ToDo-> "+err.Error())
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve rows affected value-> "+err.Error())
	}
	if rows == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", request.Id))
	}
	return &v1.DeleteResponse{
		Api:     enum.ApiVersion,
		Deleted: rows,
	}, nil
}

func (s toDoServiceServer) ReadAll(ctx context.Context, request *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	// 检查服务器是否支持客户端请求的API版本
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}
	// 从池中获取sql连接
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	// 获取ToDo列表
	rows, err := c.QueryContext(ctx, "SELECT id,title,description,reminder FROM m_todo")
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to select from ToDo-> "+err.Error())
	}
	defer rows.Close()

	var reminder time.Time
	var list []*v1.ToDo

	for rows.Next() {
		td := new(v1.ToDo)
		if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
			return nil, status.Error(codes.Unknown, "failed to retrieve field values from ToDo row-> "+err.Error())
		}

		td.Reminder, err = ptypes.TimestampProto(reminder)

		if err != nil {
			return nil, status.Error(codes.Unknown, "reminder field has invalid format-> "+err.Error())
		}
		list = append(list, td)
	}
	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "failed to retrieve data from ToDo-> "+err.Error())
	}
	return &v1.ReadAllResponse{
		Api:   enum.ApiVersion,
		ToDos: list,
	}, nil
}

// NewToDoServiceServer创建ToDo服务
func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{db}
}

// checkAPI检查服务器是否支持客户端请求的API版本
func (s *toDoServiceServer) checkApi(api string) error {
	if len(api) > 0 {
		if enum.ApiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", enum.ApiVersion, api)
		}
	}
	return nil
}

//获取数据库连接
func (s *toDoServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}

	return c, nil
}
