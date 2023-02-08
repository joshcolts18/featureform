package main

import (
	"fmt"
	"github.com/featureform/logging"
	"github.com/featureform/provider"
	"github.com/joho/godotenv"
	"os"
	"time"

	"github.com/featureform/backup"
	help "github.com/featureform/helpers"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	godotenv.Load(".env")
	etcdHost := help.GetEnv("ETCD_HOSTNAME", "localhost")
	etcdPort := help.GetEnv("ETCD_PORT", "2379")
	etcdUsername := help.GetEnv("ETCD_USERNAME", "")
	etcdPassword := help.GetEnv("ETCD_PASSWORD", "")

	address := fmt.Sprintf("%s:%s", etcdHost, etcdPort)

	etcdConfig := clientv3.Config{
		Endpoints:   []string{address},
		DialTimeout: time.Second * 10,
		Username:    etcdUsername,
		Password:    etcdPassword,
	}

	client, err := clientv3.New(etcdConfig)
	if err != nil {
		panic(err)
	}

	p := provider.FileStoreType(help.GetEnv("CLOUD_PROVIDER", provider.FileSystem))

	var backupProvider backup.Provider
	switch p {
	case provider.Azure:
		backupProvider = &backup.Azure{
			AccountName:   os.Getenv("AZURE_STORAGE_ACCOUNT"),
			AccountKey:    os.Getenv("AZURE_STORAGE_KEY"),
			ContainerName: os.Getenv("AZURE_CONTAINER_NAME"),
			Path:          os.Getenv("AZURE_STORAGE_PATH"),
		}
	case provider.GCS:
		backupProvider = &backup.GCS{
			BucketName:  os.Getenv("GCS_BUCKET_NAME"),
			BucketPath:  os.Getenv("GCS_BUCKET_PATH"),
			Credentials: []byte(help.GetEnv("GCS_CREDENTIALS", "")), // Uses local creds if empty
		}
	case provider.S3:
		backupProvider = &backup.S3{
			AWSAccessKeyId: os.Getenv("AWS_ACCESS_KEY"),
			AWSSecretKey:   os.Getenv("AWS_SECRET_KEY"),
			BucketRegion:   os.Getenv("AWS_BUCKET_REGION"),
			BucketName:     os.Getenv("AWS_BUCKET_NAME"),
			BucketPath:     os.Getenv("AWS_BUCKET_PATH"),
		}
	case provider.FileSystem:
		backupProvider = &backup.Local{
			Path: help.GetEnv("LOCAL_FILESTORE_PATH", "file://./"),
		}
	default:
		panic(fmt.Errorf("the cloud provider '%s' is not supported", p))
	}

	logger := logging.NewLogger("Restore")

	backupExecutor := backup.BackupManager{
		ETCDClient: client,
		Provider:   backupProvider,
		Logger:     logger,
	}

	err = backupExecutor.Save()
	if err != nil {
		panic(err)
	}

}