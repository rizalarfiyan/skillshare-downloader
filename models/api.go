package models

type ClassApiResponse struct {
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
		Self     ClassApiLink `json:"self"`
		Teacher  ClassApiLink `json:"teacher"`
		Units    ClassApiLink `json:"units"`
		Sessions ClassApiLink `json:"sessions"`
		Category ClassApiLink `json:"category"`
		Students ClassApiLink `json:"students"`
		Projects ClassApiLink `json:"projects"`
		Reviews  ClassApiLink `json:"reviews"`
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
				Self          ClassApiLink `json:"self"`
				MyClasses     ClassApiLink `json:"my_classes"`
				Completions   ClassApiLink `json:"completions"`
				Discussions   ClassApiLink `json:"discussions"`
				Followers     ClassApiLink `json:"followers"`
				Following     ClassApiLink `json:"following"`
				Notifications ClassApiLink `json:"notifications"`
				Projects      ClassApiLink `json:"projects"`
				Rosters       ClassApiLink `json:"rosters"`
				UserTags      ClassApiLink `json:"userTags"`
				Votes         ClassApiLink `json:"votes"`
				Wishlist      ClassApiLink `json:"wishlist"`
			} `json:"_links"`
		} `json:"teacher"`
		Units    []any `json:"units"`
		Sessions struct {
			Links struct {
				Self ClassApiLink `json:"self"`
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
						Self        ClassApiLink `json:"self"`
						Download    ClassApiLink `json:"download"`
						ParentClass ClassApiLink `json:"parentClass"`
						Stream      ClassApiLink `json:"stream"`
						Unit        any          `json:"unit"`
					} `json:"_links"`
				} `json:"sessions"`
			} `json:"_embedded"`
		} `json:"sessions"`
	} `json:"_embedded"`
}

type ClassApiLink struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}
