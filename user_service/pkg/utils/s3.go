package utils

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nfnt/resize"
	"github.com/pavel/user_service/config"
	"github.com/pavel/user_service/pkg/logger"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type Aws interface {
	UploadFileToS3(reader io.Reader, folder string, width, height uint) (string, error)
	DeleteFileFromS3(filename string)
}

type AwsService struct {
	cfg  *config.Config
	log  *logger.Logger
	sess *session.Session
}

func InitAwsService(cfg *config.Config, log *logger.Logger, sess *session.Session) Aws {
	return &AwsService{
		cfg:  cfg,
		log:  log,
		sess: sess,
	}
}

func (a AwsService) DeleteFileFromS3(filename string) {
	request := &s3.DeleteObjectInput{
		Bucket: aws.String(a.cfg.Aws.BucketName),
		Key:    aws.String(filename),
	}
	client := s3.New(a.sess)
	_, err := client.DeleteObject(request)
	if err != nil {
		a.log.Errorf("File wasn't delete from: %v", err)
	}
}

func (a AwsService) UploadFileToS3(reader io.Reader, folder string, width, height uint) (string, error) {

	img, _, err := image.Decode(reader)
	if err != nil {
		a.log.Errorf("Image decode error: %v", err)
		return "", err
	}
	resizedImg := resize.Resize(width, height, img, resize.Lanczos3)

	uploadedFile, err := a.saveTemplate(resizedImg)
	if err != nil {
		return "", err
	}
	s3URL, err := a.uploadToS3(folder, uploadedFile)
	if err != nil {
		return "", err
	}
	a.deleteFromTemplate(uploadedFile)
	return s3URL, nil
}

func (a AwsService) uploadToS3(folder, filename string) (string, error) {
	templateFolder := a.getTemplateFolder()
	file, err := os.Open(fmt.Sprintf("%s%s", templateFolder, filename))
	if err != nil {
		a.log.Errorf("Cannot open file. Err %v", err)
		return "", err
	}
	defer file.Close()
	s3Filename := fmt.Sprintf("%s/%s", folder, filename)
	uploader := s3manager.NewUploader(a.sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.cfg.Aws.BucketName),
		ACL:    aws.String("public-read"),
		Key:    aws.String(s3Filename),
		Body:   file,
	})
	if err != nil {
		a.log.Errorf("Cannot upload file to Amazon. Err %v", err)
		return "", err
	}

	return "https://" + a.cfg.Aws.BucketName + "." + "s3-" + a.cfg.Aws.Region + ".amazonaws.com/" + s3Filename, nil

}

func (a AwsService) deleteFromTemplate(filename string) {
	err := os.Remove(fmt.Sprintf("%s/%s", a.getTemplateFolder(), filename))
	if err != nil {
		a.log.Errorf("Cannot delete file `%s` . Err: %v", fmt.Sprintf("%s/%s", a.getTemplateFolder(), filename), err)
	}
}

func (a AwsService) getTemplateFolder() string {
	_, b, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(b), "../..") + "/uploads/"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		a.log.Errorf("Can't create template uploads folder. Err %v", err)
	}

	return path
}

func (a AwsService) saveTemplate(resizedImg image.Image) (string, error) {
	path := a.getTemplateFolder()
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, resizedImg, nil)
	if err != nil {
		a.log.Error("encode jpeg. Err %v", err)
		return "", err
	}
	filename := fmt.Sprintf("%s%s", GetRandomString(32), ".jpg")
	err = os.WriteFile(path+filename, buf.Bytes(), 0755)
	if err != nil {
		a.log.Errorf("error write to file. Err %v", err)
		return "", err
	}

	return filename, nil
}
