package main

import (
	"context"
	"flag"
	v1 "github.com/WishZ/go-grpc-demo/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"time"
)

const apiVersion = "v1"

func main() {
	//获取配置
	address := flag.String("server", "", "gRPC server in format host:port")
	flag.Parse()

	//建立连接
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect:%v", err)
	}
	defer conn.Close()
	c := v1.NewToDoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)

	req1 := v1.CreateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Title:       "title (" + pfx + ")",
			Description: "description (" + pfx + ")",
			Reminder:    reminder,
		},
	}
	res1, err := c.Create(ctx, &req1)
	log.Printf("Create result: <%+v>\n\n", res1)
	id := res1.Id
	// Read
	req2 := v1.ReadRequest{
		Api: apiVersion,
		Id:  id,
	}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatalf("Read failed: %v", err)
	}
	log.Printf("Read result: <%+v>\n\n", res2)
	// Update
	req3 := v1.UpdateRequest{
		Api: apiVersion,
		ToDo: &v1.ToDo{
			Id:          res2.Todo.Id,
			Title:       res2.Todo.Title,
			Description: res2.Todo.Description + " + updated",
			Reminder:    res2.Todo.Reminder,
		},
	}
	res3, err := c.Update(ctx, &req3)
	if err != nil {
		log.Fatalf("Update failed: %v", err)
	}
	log.Printf("Update result: <%+v>\n\n", res3)
	// Call ReadAll
	req4 := v1.ReadAllRequest{
		Api: apiVersion,
	}
	res4, err := c.ReadAll(ctx, &req4)
	if err != nil {
		log.Fatalf("ReadAll failed: %v", err)
	}
	log.Printf("ReadAll result: <%+v>\n\n", res4)
	// Delete
	req5 := v1.DeleteRequest{
		Api: apiVersion,
		Id:  id,
	}
	res5, err := c.Delete(ctx, &req5)
	if err != nil {
		log.Fatalf("Delete failed: %v", err)
	}
	log.Printf("Delete result: <%+v>\n\n", res5)
}
