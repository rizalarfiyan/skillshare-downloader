package models

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rizalarfiyan/skillshare-downloader/utils"
)

type ClassData struct {
	ID                         int    `json:"id"`
	Gid                        string `json:"gid"`
	Sku                        int    `json:"sku"`
	Title                      string `json:"title"`
	ProjectTitle               string `json:"project_title"`
	ImageHuge                  string `json:"image_huge"`
	ImageSmall                 string `json:"image_small"`
	ImageThumbnail             string `json:"image_thumbnail"`
	WebURL                     string `json:"web_url"`
	EnrollmentType             int    `json:"enrollment_type"`
	Category                   string `json:"category"`
	Price                      any    `json:"price"`
	NumVideos                  int    `json:"num_videos"`
	TotalVideosDuration        string `json:"total_videos_duration"`
	TotalVideosDurationSeconds int    `json:"total_videos_duration_seconds"`
	NumReviews                 int    `json:"num_reviews"`
	NumPositiveReviews         int    `json:"num_positive_reviews"`
	NumStudents                int    `json:"num_students"`
	NumProjects                int    `json:"num_projects"`
	NumDiscussions             int    `json:"num_discussions"`
	IsStaffPick                bool   `json:"is_staff_pick"`
	IsSkillshareProduced       bool   `json:"is_skillshare_produced"`
	RelativePublishTime        string `json:"relative_publish_time"`
	Actions                    []any  `json:"actions"`
	Links                      struct {
		Self     ClassDataLink `json:"self"`
		Teacher  ClassDataLink `json:"teacher"`
		Units    ClassDataLink `json:"units"`
		Sessions ClassDataLink `json:"sessions"`
		Category ClassDataLink `json:"category"`
		Students ClassDataLink `json:"students"`
		Projects ClassDataLink `json:"projects"`
		Reviews  ClassDataLink `json:"reviews"`
	} `json:"_links"`
	Embedded struct {
		Teacher struct {
			ID               int    `json:"id"`
			Gid              string `json:"gid"`
			Username         int    `json:"username"`
			FirstName        string `json:"first_name"`
			LastName         string `json:"last_name"`
			FullName         string `json:"full_name"`
			Headline         string `json:"headline"`
			URL              string `json:"url"`
			Pic              string `json:"pic"`
			PicSm            string `json:"pic_sm"`
			PicLg            string `json:"pic_lg"`
			IsTeacher        bool   `json:"is_teacher"`
			IsTopTeacher     bool   `json:"is_top_teacher"`
			NumFollowers     int    `json:"numFollowers"`
			NumFollowing     int    `json:"numFollowing"`
			IsProfilePrivate bool   `json:"is_profile_private"`
			VanityUsername   string `json:"vanity_username"`
			Links            struct {
				Self          ClassDataLink `json:"self"`
				MyClasses     ClassDataLink `json:"my_classes"`
				Completions   ClassDataLink `json:"completions"`
				Discussions   ClassDataLink `json:"discussions"`
				Followers     ClassDataLink `json:"followers"`
				Following     ClassDataLink `json:"following"`
				Notifications ClassDataLink `json:"notifications"`
				Projects      ClassDataLink `json:"projects"`
				Rosters       ClassDataLink `json:"rosters"`
				UserTags      ClassDataLink `json:"userTags"`
				Votes         ClassDataLink `json:"votes"`
				Wishlist      ClassDataLink `json:"wishlist"`
			} `json:"_links"`
		} `json:"teacher"`
		Units    []any `json:"units"`
		Sessions struct {
			Links struct {
				Self ClassDataLink `json:"self"`
			} `json:"_links"`
			Embedded struct {
				Sessions []struct {
					ID                   int    `json:"id"`
					ParentClassSku       int    `json:"parent_class_sku"`
					UnitID               int    `json:"unit_id"`
					Index                int    `json:"index"`
					Title                string `json:"title"`
					Rank                 int    `json:"rank"`
					LastPlayedTime       int    `json:"last_played_time"`
					VideoHashedID        string `json:"video_hashed_id"`
					VideoHashedIDAlt     any    `json:"video_hashed_id_alt"`
					VideoDuration        string `json:"video_duration"`
					VideoDurationSeconds int    `json:"video_duration_seconds"`
					VideoThumbnailURL    string `json:"video_thumbnail_url"`
					VideoMidThumbnailURL string `json:"video_mid_thumbnail_url"`
					ImageThumbnail       string `json:"image_thumbnail"`
					CreateTime           string `json:"create_time"`
					UpdateTime           string `json:"update_time"`
					IsCloudflareReady    bool   `json:"is_cloudflare_ready"`
					Links                struct {
						Self        ClassDataLink `json:"self"`
						Download    ClassDataLink `json:"download"`
						ParentClass ClassDataLink `json:"parentClass"`
						Stream      ClassDataLink `json:"stream"`
						Unit        any           `json:"unit"`
					} `json:"_links"`
				} `json:"sessions"`
			} `json:"_embedded"`
		} `json:"sessions"`
	} `json:"_embedded"`
}

type ClassDataLink struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

func (cd *ClassData) IsValidVideoId() bool {
	for _, session := range cd.Embedded.Sessions.Embedded.Sessions {
		if session.VideoHashedID == "" {
			return false
		}
	}

	return true
}

type SkillshareClass struct {
	ID                         int               `json:"id"`
	Title                      string            `json:"title"`
	ProjectTitle               string            `json:"project_title"`
	Category                   string            `json:"category"`
	TotalVideosDuration        string            `json:"total_videos_duration"`
	TotalVideosDurationSeconds int               `json:"total_videos_duration_seconds"`
	ImageHuge                  string            `json:"image_huge"`
	ImageSmall                 string            `json:"image_small"`
	ImageThumbnail             string            `json:"image_thumbnail"`
	Videos                     []SkillshareVideo `json:"videos"`
}

type SkillshareVideo struct {
	ID                   int                       `json:"id"`
	Title                string                    `json:"title"`
	VideoID              string                    `json:"video_id"`
	VideoDuration        string                    `json:"video_duration"`
	VideoDurationSeconds int                       `json:"video_duration_seconds"`
	VideoThumbnailURL    string                    `json:"video_thumbnail_url"`
	VideoMidThumbnailURL string                    `json:"video_mid_thumbnail_url"`
	ImageThumbnail       string                    `json:"image_thumbnail"`
	Sources              []SkillshareVideoSource   `json:"sources"`
	Subtitles            []SkillshareVideoSubtitle `json:"subtitles"`
}

type SkillshareVideoSource struct {
	Src        string `json:"src"`
	AvgBitrate int    `json:"avg_bitrate"`
	Codec      string `json:"codec"`
	Container  string `json:"container"`
	Duration   int    `json:"duration"`
	Height     int    `json:"height"`
	Size       int    `json:"size"`
	Width      int    `json:"width"`
}

type SkillshareVideoSubtitle struct {
	Src   string `json:"src"`
	Lang  string `json:"lang"`
	Label string `json:"label"`
}

func (cd *ClassData) Mapper() SkillshareClass {
	ssData := SkillshareClass{
		ID:                         cd.ID,
		Title:                      utils.DecodeAscii(cd.Title),
		ProjectTitle:               cd.ProjectTitle,
		Category:                   cd.Category,
		TotalVideosDuration:        cd.TotalVideosDuration,
		TotalVideosDurationSeconds: cd.TotalVideosDurationSeconds,
		ImageHuge:                  cd.ImageHuge,
		ImageSmall:                 cd.ImageSmall,
		ImageThumbnail:             cd.ImageThumbnail,
	}

	for _, session := range cd.Embedded.Sessions.Embedded.Sessions {
		var videoId int
		videoArr := strings.Split(session.VideoHashedID, ":")
		if len(videoArr) > 1 {
			videoId, _ = strconv.Atoi(videoArr[1])
		}

		ssData.Videos = append(ssData.Videos, SkillshareVideo{
			ID:                   videoId,
			Title:                utils.DecodeAscii(session.Title),
			VideoID:              session.VideoHashedID,
			VideoDuration:        session.VideoDuration,
			VideoDurationSeconds: session.VideoDurationSeconds,
			VideoThumbnailURL:    session.VideoThumbnailURL,
			VideoMidThumbnailURL: session.VideoMidThumbnailURL,
			ImageThumbnail:       session.ImageThumbnail,
		})
	}

	return ssData
}

func (sc *SkillshareVideo) AddSourceSubtitle(video VideoData) {
	sources := []SkillshareVideoSource{}
	tempIdx := make(map[string]bool)
	sort.Slice(video.Sources, func(i, j int) bool {
		return video.Sources[i].Src > video.Sources[j].Src
	})
	for _, source := range video.Sources {
		if source.Codecs == "avc1,mp4a" || source.Codec == "" {
			continue
		}
		key := strings.TrimPrefix(source.Src, "https://")
		key = strings.TrimPrefix(key, "http://")
		if _, isExist := tempIdx[key]; !isExist {
			tempIdx[key] = true
			sources = append(sources, SkillshareVideoSource{
				Src:        source.Src,
				AvgBitrate: source.AvgBitrate,
				Codec:      source.Codec,
				Container:  source.Container,
				Duration:   source.Duration,
				Height:     source.Height,
				Size:       source.Size,
				Width:      source.Width,
			})
		}
	}

	sc.Sources = sources

	subtitles := []SkillshareVideoSubtitle{}
	for _, subtitle := range video.TextTracks {
		if subtitle.Kind != "subtitles" {
			continue
		}
		subtitles = append(subtitles, SkillshareVideoSubtitle{
			Src:   subtitle.Src,
			Lang:  subtitle.Srclang,
			Label: subtitle.Label,
		})
	}
	sc.Subtitles = subtitles
}

type VideoData struct {
	Poster           string            `json:"poster"`
	Thumbnail        string            `json:"thumbnail"`
	PosterSources    []VideoDataSource `json:"poster_sources"`
	ThumbnailSources []VideoDataSource `json:"thumbnail_sources"`
	Description      any               `json:"description"`
	Tags             []any             `json:"tags"`
	CuePoints        []any             `json:"cue_points"`
	CustomFields     struct {
	} `json:"custom_fields"`
	AccountID string `json:"account_id"`
	Sources   []struct {
		Codecs      string `json:"codecs,omitempty"`
		ExtXVersion string `json:"ext_x_version,omitempty"`
		Src         string `json:"src"`
		Type        string `json:"type,omitempty"`
		Profiles    string `json:"profiles,omitempty"`
		AvgBitrate  int    `json:"avg_bitrate,omitempty"`
		Codec       string `json:"codec,omitempty"`
		Container   string `json:"container,omitempty"`
		Duration    int    `json:"duration,omitempty"`
		Height      int    `json:"height,omitempty"`
		Size        int    `json:"size,omitempty"`
		Width       int    `json:"width,omitempty"`
	} `json:"sources"`
	Name            string `json:"name"`
	ReferenceID     any    `json:"reference_id"`
	LongDescription any    `json:"long_description"`
	Duration        int    `json:"duration"`
	Economics       string `json:"economics"`
	TextTracks      []struct {
		ID        string            `json:"id"`
		AccountID string            `json:"account_id"`
		Src       string            `json:"src"`
		Srclang   string            `json:"srclang"`
		Label     string            `json:"label"`
		Kind      string            `json:"kind"`
		MimeType  string            `json:"mime_type"`
		AssetID   any               `json:"asset_id"`
		Sources   []VideoDataSource `json:"sources"`
		Default   bool              `json:"default"`
	} `json:"text_tracks"`
	PublishedAt    time.Time `json:"published_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	OfflineEnabled bool      `json:"offline_enabled"`
	Link           any       `json:"link"`
	ID             string    `json:"id"`
	AdKeys         any       `json:"ad_keys"`
}

type VideoDataSource struct {
	Src string `json:"src"`
}
