package imageOperation

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
)

type modifiedImage struct {
	ImageName    string `json:"imageName"`
	ImageVersion string `json:"imageVersion"`
	ImageSize    int64  `json:"imageSize"`
}

func ListImages() []modifiedImage {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Println(err)
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	imageList, err := client.ImageService().List(ctx)

	var modifiedImageList []modifiedImage
	for _, image := range imageList {
		oldName := image.Name
		var name string
		var version string
		if strings.Contains(oldName, "@") {
			name = oldName[:strings.LastIndex(oldName, "@")]
			version = ""
		} else {
			if strings.Contains(oldName, "sha256") {
				name = oldName
				version = ""
			} else {
				lastIndex := strings.LastIndex(oldName, ":")
				if lastIndex != -1 {
					name = oldName[:lastIndex]
					version = oldName[lastIndex+1:]
				}
			}
		}
		size := image.Target.Size
		modifyingImage := modifiedImage{name, version, size}
		modifiedImageList = append(modifiedImageList, modifyingImage)
	}

	if err != nil {
		log.Println(err)
	}

	return modifiedImageList
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
