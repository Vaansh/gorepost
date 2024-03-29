package subscriber

import (
	"fmt"
	"strings"
	"time"

	"github.com/Vaansh/gore"
	"github.com/Vaansh/gore/internal/gcloud"
	"github.com/Vaansh/gore/internal/http"
	"github.com/Vaansh/gore/internal/model"
	"github.com/Vaansh/gore/internal/repository"
	"github.com/Vaansh/gore/internal/util"
)

type InstagramSubscriber struct {
	instagramId string
	hashtags    string
	frequency   time.Duration
	repository  repository.UserRepository
	client      *http.InstagramClient
}

func NewInstagramSubscriber(instagramId string, metadata model.MetaData, repository repository.UserRepository) *InstagramSubscriber {
	return &InstagramSubscriber{
		instagramId: instagramId,
		hashtags:    metadata.IgPostTags,
		frequency:   metadata.Frequency,
		repository:  repository,
		client:      http.NewInstagramClient(metadata.IgUserId, metadata.IgLongLivedAccessToken),
	}
}

func (s *InstagramSubscriber) SubscribeTo(c <-chan model.Post) {
	for post := range c {
		postId, author, sourcePlatform, caption := post.GetParams()
		exists, err := s.repository.CheckIfRecordExists(s.getTableName(), &post)
		if err != nil {
			gcloud.LogWarning(err.Error())
		}

		if !exists {
			err := gcloud.SaveFile(postId, sourcePlatform)
			if err != nil {
				gcloud.LogError(err.Error())
				break
			}

			coverUrl, err := gcloud.GetCoverUrl(postId, sourcePlatform)
			if err != nil {
				gcloud.LogWarning(err.Error())
			}

			fileName := sourcePlatform.String() + "_" + postId + ".mp4"
			err = gcloud.UploadToBucket(fileName)
			if err != nil {
				gcloud.LogError(err.Error())
				break
			}

			fileUrl := gcloud.GetFileUrl(fileName)
			err = s.client.UploadReel(fileUrl, coverUrl, util.GenerateInstagramCaption(caption, author, s.hashtags, strings.ToUpper(sourcePlatform.String())))
			if err != nil {
				gcloud.LogError(err.Error())
				break
			}

			err = gcloud.DeleteFromBucket(fileName)
			if err != nil {
				gcloud.LogWarning(err.Error())
			}

			err = s.repository.AddRecord(s.getTableName(), &post)
			if err != nil {
				gcloud.LogError(err.Error())
				break
			}

			gcloud.DeleteFile(fileName)
			gcloud.LogInfo(fmt.Sprintf("Instagram subscriber (%s) processed post %s from %s", s.GetSubscriberId(), post.PostId, sourcePlatform))
			time.Sleep(s.frequency)
		}
	}
}

func (s *InstagramSubscriber) GetSubscriberId() string {
	return s.instagramId
}

func (s *InstagramSubscriber) getTableName() string {
	return gore.INSTAGRAM.String() + "_" + s.instagramId
}
