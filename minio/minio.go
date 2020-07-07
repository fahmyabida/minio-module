package minio

import (
	"errors"
	"github.com/minio/minio-go/v6"
	"io"
	"net/url"
	"strings"
	"time"
)

type MinioEngine struct {
	minioClient *minio.Client
}

type IMinioEngine interface {
	AddBucket(bucketName string) error
	GetFile(bucketName, objectNameInMinio string, durationShare time.Duration) (*ResultShareFile, error)
	DownloadFile(bucketName, objectNameInMinio string, durationShare time.Duration) (*ResultShareFile, error)
	UploadFileWithPathFile(bucketName, objectNameWillBeInMinioWithExtension, pathOfOriginFile string) (*ResultUploadFile,
																										error)
	UploadFileWithFile(bucketName, objectNameWillBeInMinioWithExtension string, uploadFile io.Reader) (*ResultUploadFile,
																										error)
}

func NewMinioHelper(minioClient *minio.Client) IMinioEngine{
	return &MinioEngine{minioClient}
}

// AddNewBucket
func (m *MinioEngine) AddBucket(bucketName string) error{
	location := ""
	err := m.minioClient.MakeBucket(bucketName, location)
	if err != nil {
		exists, errBucketExists := m.minioClient.BucketExists(bucketName)
		if errBucketExists == nil && exists {
			return nil
		} else {
			return err
		}
	}
	return nil
}

type ResultShareFile struct {
	PathUrl string
	Duration string
	ValidUntil string
}

// Get File in temporary link with time duration
func (m *MinioEngine) GetFile(bucketName, objectNameInMinio string, durationShare time.Duration) (*ResultShareFile,
																									error){
	return m.getShareableFile("inline", bucketName, objectNameInMinio, durationShare)
}

// Download File in temporary link with time duration
func (m *MinioEngine) DownloadFile(bucketName, objectNameInMinio string, durationShare time.Duration) (*ResultShareFile,
																										error){
	return m.getShareableFile("attachment", bucketName, objectNameInMinio, durationShare)
}

func (m *MinioEngine) getShareableFile(contentDisposition, bucketName, objectNameInMinio string,
										durationShare time.Duration) (*ResultShareFile, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", contentDisposition)
	contentType, _ := m.getMIMEcontentType(objectNameInMinio)
	if contentType != "" {
		reqParams.Set("content-type", contentType)
	}
	url, err := m.minioClient.PresignedGetObject(bucketName, objectNameInMinio, durationShare, reqParams)
	if err != nil {
		return nil, err
	}
	return &ResultShareFile{
		PathUrl: url.Host+url.Path+"?"+url.RawQuery,
		Duration: durationShare.String(),
		ValidUntil: time.Now().Add(durationShare).Format(time.RFC3339),
	}, nil
}

// UploadFile use contentType which the content is already define on const ^^
func (m *MinioEngine) UploadFileWithPathFile(bucketName, objectNameWillBeInMinioWithExtension,
										pathOfOriginFile string) (*ResultUploadFile, error){
	MIMEcontentType, err := m.getMIMEcontentType(objectNameWillBeInMinioWithExtension)
	if err != nil {
		return nil, err
	}
	_, err = m.minioClient.FPutObject(bucketName, objectNameWillBeInMinioWithExtension, pathOfOriginFile,
									minio.PutObjectOptions{ContentDisposition:MIMEcontentType})
	if err != nil {
		return nil, err
	}
	return &ResultUploadFile{
		Bucket:     bucketName,
		MimeType:   MIMEcontentType,
		ObjectFile: objectNameWillBeInMinioWithExtension,
	}, nil
}

type ResultUploadFile struct {
	Bucket string
	MimeType string
	ObjectFile string
}

// UploadFile with type of io.reader
func (m *MinioEngine) UploadFileWithFile(bucketName, objectNameWillBeInMinioWithExtension string,
									uploadFile io.Reader) (*ResultUploadFile, error){
	MIMEcontentType, err := m.getMIMEcontentType(objectNameWillBeInMinioWithExtension)
	if err != nil {
		return nil, err
	}
	_, err = m.minioClient.PutObject(bucketName, objectNameWillBeInMinioWithExtension, uploadFile, -1,
									minio.PutObjectOptions{ContentType:MIMEcontentType})
	if err != nil {
		return nil, err
	}
	return &ResultUploadFile{
		Bucket:     bucketName,
		MimeType:   MIMEcontentType,
		ObjectFile: objectNameWillBeInMinioWithExtension,
	}, nil
}

func (m *MinioEngine) getMIMEcontentType(fileName string) (mimeContentType string, err error){
	arrFileName := strings.Split(fileName, ".")
	if len(arrFileName) == 0 {
		return "", errors.New("file not found")
	} else if len(arrFileName) != 2 {
		return "", errors.New("file extension not valid")
	}
	extension := arrFileName[1]
	switch extension {
	case "jpg", "jpeg":
		return "image/jpeg", nil
	case "zip":
		return "application/zip", nil
	case "xslx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
	case "xls":
		return "application/vnd.ms-excel", nil
	case "csv" :
		return "text/csv", nil
	case "doc" :
		return "application/msword", nil
	case "docx" :
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", nil
	case "pdf" :
		return "application/pdf", nil
	default:
		return "", errors.New("extension file not registered yet")
	}
}



