package main

import (
	"context"
	"errors"
	"log"
	"os"

	pb "github.com/aaronflower/shippy-service-consignment/proto/consignment"
	userService "github.com/aaronflower/shippy-service-user/proto/auth"
	vesselProto "github.com/aaronflower/shippy-service-vessel/proto/vessel"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
)

const (
	defaultHost = "localhost:27017"
	port        = ":50051"
)

// AuthWrapper is a high-order function which takes a HandlerFunc and returns a function,
// which takes a context, request and response interface. The token is extracted from the
// context set in our consignment-cli, that token is then sent over to ther user service to be validate.
// If valid, the call is passed along to the handler. If not, an error is returned.
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		// Note this is now uppercase(not entirely sure why this is ...)
		token := meta["Token"]
		log.Println("Authentication with token: ", token)
		authClient := userService.NewAuthClient("go.micro.srv.user", client.DefaultClient)
		_, err := authClient.ValidateToken(context.Background(), &userService.Token{
			Token: token,
		})
		if err != nil {
			log.Fatal(err)
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}

func main() {
	// Database host from the environment variables
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}

	session, err := CreateSession(host)
	defer session.Close()

	if err != nil {
		log.Panicf("Could not connect to datastore with host %s - %v", host, err)
	}

	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition.
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		// Our auth middleware
		micro.WrapHandler(AuthWrapper),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	// Init will parse the command line flags.
	srv.Init()

	// Register our service with the gRPC server, this will tie our implementation into
	// the auto-generated interface code for our portobuf definition.
	// Register service
	pb.RegisterShippingServiceHandler(srv.Server(), &service{session, vesselClient})

	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
