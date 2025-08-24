package protobuf

import (
	"github.com/rishirishhh/vought/src/cmd/api/models"
	contracts "github.com/rishirishhh/vought/src/pkg/contracts/v1"
	log "github.com/sirupsen/logrus"
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

var modelToProtoStatus = []contracts.Video_VideoStatus{
	models.UNSPECIFIED: contracts.Video_VIDEO_STATUS_UNSPECIFIED,
	models.UPLOADING:   contracts.Video_VIDEO_STATUS_UPLOADING,
	models.UPLOADED:    contracts.Video_VIDEO_STATUS_UPLOADED,
	models.ENCODING:    contracts.Video_VIDEO_STATUS_ENCODING,
	models.COMPLETE:    contracts.Video_VIDEO_STATUS_COMPLETE,
	models.UNKNOWN:     contracts.Video_VIDEO_STATUS_UNKNOWN,
	models.FAIL_UPLOAD: contracts.Video_VIDEO_STATUS_FAIL_UPLOAD,
	models.FAIL_ENCODE: contracts.Video_VIDEO_STATUS_FAIL_ENCODE,
}

func VideoProtobufToVideo(videoProto *contracts.Video) *models.Video {
	if videoProto == nil {
		log.Error("Cannot convert protobuf video to video, video nil")
		return nil
	}

	video := models.Video{
		ID:         videoProto.Id,
		Status:     protoToModelStatus[videoProto.Status],
		SourcePath: videoProto.Source,
		CoverPath:  videoProto.CoverPath,
	}

	return &video
}

func VideoToVideoProtobuf(video *models.Video) *contracts.Video {
	if video == nil {
		log.Error("Cannot convert protobuf video to video, video nil")
		return nil
	}

	videoData := &contracts.Video{
		Id:        video.ID,
		Status:    modelToProtoStatus[video.Status],
		Source:    video.SourcePath,
		CoverPath: video.CoverPath,
	}

	return videoData
}
