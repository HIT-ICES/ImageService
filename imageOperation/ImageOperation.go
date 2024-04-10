package imageOperation

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/platforms"
	"github.com/containerd/containerd/remotes/docker"
	"log"
	"net/http"
	"strings"
)

type DeleteImagesJSON struct {
	Image []string `json:"image"`
}

type PullImagesJSON struct {
	Image []string `json:"image"`
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

type pullImageResponse struct {
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

func pullImage(imageName string) error {
	authorizer := docker.NewDockerAuthorizer()
	if strings.HasPrefix(imageName, "192.168.1.199:5000") {
		authorizer = docker.NewDockerAuthorizer(
			docker.WithAuthClient(http.DefaultClient),
			docker.WithAuthCreds(func(host string) (string, string, error) {
				return "admin", "bfgeY8qckIhH2aeh", nil
			}),
		)
	}

	resolver := docker.NewResolver(docker.ResolverOptions{
		Hosts: docker.ConfigureDefaultRegistries(
			docker.WithPlainHTTP(docker.MatchAllHosts),
			docker.WithAuthorizer(authorizer),
		),
	})

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	_, err = client.Pull(ctx, imageName, containerd.WithPullUnpack, containerd.WithResolver(resolver))
	if err != nil {
		return err
	}

	return nil
}

func PullImages(pullImages *PullImagesJSON) *[]pullImageResponse {
	var responses []pullImageResponse
	for _, image := range pullImages.Image {
		err := pullImage(image)
		if err != nil {
			responses = append(responses, pullImageResponse{
				Image:   image,
				Success: err.Error(),
			})
		} else {
			responses = append(responses, pullImageResponse{
				Image:   image,
				Success: "pull successfully",
			})
		}
	}
	return &responses
}

func DeleteImages(deleteImages *DeleteImagesJSON) *[]deleteImageResponse {
	client, err := containerd.New("/run/containerd/containerd.sock")
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	var responses []deleteImageResponse

	for _, image := range deleteImages.Image {
		err = client.ImageService().Delete(ctx, image, images.SynchronousDelete())
		if err != nil {
			responses = append(responses, deleteImageResponse{
				Image:   image,
				Success: err.Error(),
			})
		} else {
			responses = append(responses, deleteImageResponse{
				Image:   image,
				Success: "delete successfully",
			})
		}
	}

	return &responses
}
