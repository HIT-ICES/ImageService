package imageOperation

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func ListImages() []images.Image {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Println(err)
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	imageList, err := client.ImageService().List(ctx)
	if err != nil {
		log.Println(err)
	}

	return imageList
}

func PullImage(imageName string) string {
	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	client, err := containerd.New("/run/containerd/containerd.sock",
		containerd.WithDialOpts(dialOpts))
	if err != nil {
		return err.Error()
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	image, err := client.Pull(ctx, imageName, containerd.WithPullUnpack)
	if err != nil {
		return err.Error()
	}

	return "Pull " + image.Name() + " successfully"
}

func DeleteImage(imageName string) string {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err.Error()
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	err = client.ImageService().Delete(ctx, imageName, func(ctx context.Context, options *images.DeleteOptions) error {
		options.Synchronous = true
		return nil
	})
	if err != nil {
		return err.Error()
	}

	return "Delete successfully"
}
