package backend

import (
    "context"
    "fmt"
    "io"

    "around/constants"

    "cloud.google.com/go/storage"
)

var (
    GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
    client *storage.Client // client共享 每次请求都是一样的 
    bucket string
}

func InitGCSBackend() { //怎么限定他只初始化一次？init在main main只调用一次 所以init只call一次
	// router 会多线程调用 会接受到多个请求 会并发进行 会触发多个handler
    client, err := storage.NewClient(context.Background())
	//context.background 可以设置timeout context.WithTimeout可以设置 这里就没有设置
    if err != nil {
        panic(err)
    }

    GCSBackend = &GoogleCloudStorageBackend{
		// & 返回指向对象的指针
        client: client,
        bucket: constants.GCS_BUCKET,
    }
}

func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
    // 返回的string 是url
	ctx := context.Background()
    object := backend.client.Bucket(backend.bucket).Object(objectName)
    wc := object.NewWriter(ctx)
    if _, err := io.Copy(wc, r); err != nil {
        return "", err
    }

    if err := wc.Close(); err != nil {
        return "", err
    }

	// 所有人都有reader权限 可以使得前端拿到url可以访问到具体文件 因为前后端分开所以如果是private前端没法访问
    if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
        return "", err
    }

    attrs, err := object.Attrs(ctx)
    if err != nil {
        return "", err
    }

    fmt.Printf("File is saved to GCS: %s\n", attrs.MediaLink)
    return attrs.MediaLink, nil
}