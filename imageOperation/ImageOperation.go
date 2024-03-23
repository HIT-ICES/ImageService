package imageOperation

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/platforms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
)

type DeleteImages struct {
	Image string `json:"image"`
}

type modifiedImage struct {
	ImageName    string `json:"imageName"`
	ImageVersion string `json:"imageVersion"`
	ImageSize    int64  `json:"imageSize"`
}

type deleteImageResponse struct {
	Image   string `json:"image"`
	Success string `json:"success"`
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
		//contentStore := client.ContentStore()
		//if contentStore == nil {
		//	log.Fatal("ContentStore is nil")
		//}
		//size := image.Target.Size
		size, err := image.Size(ctx, client.ContentStore(), platforms.All)
		if err != nil {
			log.Println("Error while getting image size:", err)
			continue
		}

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

	//image, err := client.ImageService().Get(ctx, imageName)
	//if err != nil {
	//	return err.Error()
	//}
	//
	//createImage, err := client.ImageService().Create(ctx, image)
	//if err != nil {
	//	return err.Error()
	//}

	image, err := client.Pull(ctx, imageName, containerd.WithPullUnpack)
	if err != nil {
		return err.Error()
	}

	return "Pull " + image.Name() + " successfully"
}

func DeleteImage(deleteImages *[]DeleteImages) *[]deleteImageResponse {
	client, err := containerd.New("/run/containerd/containerd.sock")
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	var responses []deleteImageResponse

	for _, image := range *deleteImages {
		err = client.ImageService().Delete(ctx, image.Image, images.SynchronousDelete())
		if err != nil {
			responses = append(responses, deleteImageResponse{
				Image:   image.Image,
				Success: err.Error(),
			})
		} else {
			responses = append(responses, deleteImageResponse{
				Image:   image.Image,
				Success: "delete successfully",
			})
		}
	}

	return &responses
}
