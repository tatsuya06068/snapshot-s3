package main

import (
    "context"
    "fmt"
    "log"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/rds"
)

type Event struct {
    ExportTaskIdentifier string   `json:"export_task_id"`   // 一意な名前
    SnapshotArn          string   `json:"snapshot_arn"`     // スナップショットARN
    S3Bucket             string   `json:"s3_bucket"`        // S3バケット名
    IamRoleArn           string   `json:"iam_role_arn"`     // IAMロールARN
    KmsKeyArn            string   `json:"kms_key_arn"`      // KMSキーARN
    TableNames           []string `json:"table_names"`      // エクスポートするテーブル名（任意）
}

func handler(ctx context.Context, event Event) error {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return fmt.Errorf("failed to load AWS config: %v", err)
    }

    rdsClient := rds.NewFromConfig(cfg)

    input := &rds.StartExportTaskInput{
        ExportTaskIdentifier: aws.String(event.ExportTaskIdentifier),
        SourceArn:            aws.String(event.SnapshotArn),
        S3BucketName:         aws.String(event.S3Bucket),
        IamRoleArn:           aws.String(event.IamRoleArn),
        KmsKeyId:             aws.String(event.KmsKeyArn),
    }

    // テーブルが指定されている場合だけ export-only に追加
    if len(event.TableNames) > 0 {
        input.ExportOnly = event.TableNames
    }

    out, err := rdsClient.StartExportTask(ctx, input)
    if err != nil {
        return fmt.Errorf("failed to start export task: %v", err)
    }

    log.Printf("Started export task: %s (Status: %s)\n", *out.ExportTaskIdentifier, out.Status)
    return nil
}

func main() {
    lambda.Start(handler)
}
