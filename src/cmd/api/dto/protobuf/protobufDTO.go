package protobuf

import (
	"github.com/rishirishhh/vought/src/cmd/api/models"
	contracts "github.com/rishirishhh/vought/src/pkg/contracts/v1"
)

var protoToModelStatus = []models.VideoStatus{
	contracts.Video_VIDEO_STATUS_UNSPECIFIED: models.UNSPECIFIED,
	contracts.Video_VIDEO_STATUS_UPLOADING:   models.UPLOADING,
	contracts.Video_VIDEO_STATUS_UPLOADED:    models.UPLOADED,
	contracts.Video_VIDEO_STATUS_ENCODING:    models.ENCODING,
	contracts.Video_VIDEO_STATUS_COMPLETE:    models.COMPLETE,
	contracts.Video_VIDEO_STATUS_UNKNOWN:     models.UNKNOWN,
	contracts.Video_VIDEO_STATUS_FAIL_UPLOAD: models.FAIL_UPLOAD,
	contracts.Video_VIDEO_STATUS_FAIL_ENCODE: models.FAIL_ENCODE,
}
