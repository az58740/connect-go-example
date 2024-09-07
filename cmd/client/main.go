package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	usersv1 "github.com/az5840/connect-go-example/gen/users/v1"        // generated by protoc-gen-go
	"github.com/az5840/connect-go-example/gen/users/v1/usersv1connect" // generated by protoc-gen-connect-go
	"google.golang.org/protobuf/types/known/emptypb"

	"crypto/tls"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
)

func createUsers(client usersv1connect.UsersServiceClient) {
	log.Printf("Traversing %d users", len(users))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream := client.CreateUser(ctx)
	for i := range users {
		if err := stream.Send(users[i]); err != nil {
			log.Fatalf("client.CreateUser: stream.Send(%v) failed: %v", &users[i], err)
		}
	}
	if _, err := stream.CloseAndReceive(); err != nil {
		log.Fatalf("client.CreatePerson failed: %v", err)
	}
	log.Println("Users created successfully!")

}

func usersList(client usersv1connect.UsersServiceClient) {
	fmt.Println("***Users List***")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.GetUsers(ctx, &connect.Request[emptypb.Empty]{})
	if err != nil {
		fmt.Println(connect.NewError(connect.CodeInternal, err))
	}
	for stream.Receive() {
		fmt.Println(stream.Msg())
	}

}
func getUser(client usersv1connect.UsersServiceClient, id int32) {
	req := connect.NewRequest(&usersv1.GetUserRequest{
		Id: id,
	})
	req.Header().Set("Acme-Tenant-Id", "1234")
	res, err := client.GetUser(context.Background(), req)
	if err != nil {
		fmt.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println(connectErr.Message())
			fmt.Println(connectErr.Details())
		}
		return
	}
	fmt.Println(res.Msg)
}
func main() {
	client := usersv1connect.NewUsersServiceClient(
		newInsecureClient(),
		"http://localhost:8080",
		connect.WithHTTPGet(),
	)
	//Get User
	getUser(client, 5)
	createUsers(client)
	//Get All users
	usersList(client)
}
func newInsecureClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
		},
		// Don't forget timeouts!
		Timeout: 15 * time.Second,
	}
}

var users = []*usersv1.CreateUserRequest{
	{Name: "Ali", Id: 7, Age: 52, Status: usersv1.UserStatus_USER_STATUS_AVAILABLE},
	{Name: "Zahra", Id: 8, Age: 52, Status: usersv1.UserStatus_USER_STATUS_AVAILABLE},
	{Name: "Elen", Id: 9, Age: 52, Status: usersv1.UserStatus_USER_STATUS_AVAILABLE},
	{Name: "Mohamad", Id: 10, Age: 52, Status: usersv1.UserStatus_USER_STATUS_AVAILABLE},
}
