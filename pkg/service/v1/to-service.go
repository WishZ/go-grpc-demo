package v1

import (
	"context"
	"database/sql"
	v1 "github.com/WishZ/go-grpc-demo/api/proto/v1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	apiVersion = "v1"
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
		Api: apiVersion,
		Id:  id,
	}, nil
}

func (s toDoServiceServer) Read(ctx context.Context, request *v1.ReadRequest) (*v1.ReadResponse, error) {
	panic("implement me")
}

func (s toDoServiceServer) Update(ctx context.Context, request *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	panic("implement me")
}

func (s toDoServiceServer) Delete(ctx context.Context, request *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	panic("implement me")
}

func (s toDoServiceServer) ReadAll(ctx context.Context, request *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	panic("implement me")
}

// NewToDoServiceServer创建ToDo服务
func NewToDoServiceServer(db *sql.DB) v1.ToDoServiceServer {
	return &toDoServiceServer{db}
}

// checkAPI检查服务器是否支持客户端请求的API版本
func (s *toDoServiceServer) checkApi(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
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
