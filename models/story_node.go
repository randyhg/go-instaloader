package models

type StoryNode struct {
	Audience   string `json:"audience"`
	Typename   string `json:"__typename"`
	ID         string `json:"id"`
	Dimensions struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"dimensions"`
	DisplayResources []struct {
		Src          string `json:"src"`
		ConfigWidth  int    `json:"config_width"`
		ConfigHeight int    `json:"config_height"`
	} `json:"display_resources"`
	DisplayURL             string      `json:"display_url"`
	MediaPreview           string      `json:"media_preview"`
	GatingInfo             interface{} `json:"gating_info"`
	FactCheckOverallRating interface{} `json:"fact_check_overall_rating"`
	FactCheckInformation   interface{} `json:"fact_check_information"`
	SharingFrictionInfo    struct {
		ShouldHaveSharingFriction bool        `json:"should_have_sharing_friction"`
		BloksAppURL               interface{} `json:"bloks_app_url"`
	} `json:"sharing_friction_info"`
	MediaOverlayInfo        interface{} `json:"media_overlay_info"`
	SensitivityFrictionInfo interface{} `json:"sensitivity_friction_info"`
	TakenAtTimestamp        int         `json:"taken_at_timestamp"`
	ExpiringAtTimestamp     int         `json:"expiring_at_timestamp"`
	StoryCtaURL             interface{} `json:"story_cta_url"`
	StoryViewCount          int         `json:"story_view_count"`
	IsVideo                 bool        `json:"is_video"`
	Owner                   struct {
		ID                string `json:"id"`
		ProfilePicURL     string `json:"profile_pic_url"`
		Username          string `json:"username"`
		FollowedByViewer  bool   `json:"followed_by_viewer"`
		RequestedByViewer bool   `json:"requested_by_viewer"`
	} `json:"owner"`
	TrackingToken          string        `json:"tracking_token"`
	TappableObjects        []interface{} `json:"tappable_objects"`
	StoryAppAttribution    interface{}   `json:"story_app_attribution"`
	EdgeMediaToSponsorUser struct {
		Edges []interface{} `json:"edges"`
	} `json:"edge_media_to_sponsor_user"`
	MutingInfo   interface{} `json:"muting_info"`
	IphoneStruct struct {
		TakenAt                             int     `json:"taken_at"`
		Pk                                  int64   `json:"pk"`
		ID                                  string  `json:"id"`
		HideViewAllCommentEntrypoint        bool    `json:"hide_view_all_comment_entrypoint"`
		IsVisualReplyCommenterNoticeEnabled bool    `json:"is_visual_reply_commenter_notice_enabled"`
		LikeAndViewCountsDisabled           bool    `json:"like_and_view_counts_disabled"`
		StickerTranslationsEnabled          bool    `json:"sticker_translations_enabled"`
		IsPostLiveClipsMedia                bool    `json:"is_post_live_clips_media"`
		IsReshareOfTextPostAppMediaInIg     bool    `json:"is_reshare_of_text_post_app_media_in_ig"`
		IsReelMedia                         bool    `json:"is_reel_media"`
		Fbid                                int64   `json:"fbid"`
		DeviceTimestamp                     int64   `json:"device_timestamp"`
		CaptionIsEdited                     bool    `json:"caption_is_edited"`
		StrongID                            string  `json:"strong_id__"`
		DeletedReason                       int     `json:"deleted_reason"`
		HasSharedToFb                       int     `json:"has_shared_to_fb"`
		ExpiringAt                          int     `json:"expiring_at"`
		ShouldRequestAds                    bool    `json:"should_request_ads"`
		HasDelayedMetadata                  bool    `json:"has_delayed_metadata"`
		MezqlToken                          string  `json:"mezql_token"`
		ShowStoryDeletedErrorLabel          bool    `json:"show_story_deleted_error_label"`
		CommentThreadingEnabled             bool    `json:"comment_threading_enabled"`
		IsTerminalVideoSegment              bool    `json:"is_terminal_video_segment"`
		IsUnifiedVideo                      bool    `json:"is_unified_video"`
		HasPrivatelyLiked                   bool    `json:"has_privately_liked"`
		CommercialityStatus                 string  `json:"commerciality_status"`
		ImportedTakenAt                     int     `json:"imported_taken_at"`
		FilterType                          int     `json:"filter_type"`
		ClientCacheKey                      string  `json:"client_cache_key"`
		CaptionPosition                     float64 `json:"caption_position"`
		TimezoneOffset                      int     `json:"timezone_offset"`
		IntegrityReviewDecision             string  `json:"integrity_review_decision"`
		CommentingDisabledForViewer         bool    `json:"commenting_disabled_for_viewer"`
		IgIabPostClickData                  struct {
			EligibleExperienceTypes               []string `json:"eligibleExperienceTypes"`
			BuyWithPrimeIABPostClickDataExtension struct {
				BuyWithPrimeExperienceType string      `json:"buyWithPrimeExperienceType"`
				APIKey                     string      `json:"apiKey"`
				ClientID                   string      `json:"clientID"`
				AccessToken                string      `json:"accessToken"`
				AccessTokenTTL             int         `json:"accessTokenTTL"`
				AccessTokenCreationTime    interface{} `json:"accessTokenCreationTime"`
				PageName                   string      `json:"pageName"`
				BauProductURL              string      `json:"bauProductUrl"`
			} `json:"buyWithPrimeIABPostClickDataExtension"`
		} `json:"ig_iab_post_click_data"`
		IsCommentsGifComposerEnabled    bool          `json:"is_comments_gif_composer_enabled"`
		CommentInformTreatment          interface{}   `json:"comment_inform_treatment"`
		ClipsTabPinnedUserIds           []interface{} `json:"clips_tab_pinned_user_ids"`
		CanViewerSave                   bool          `json:"can_viewer_save"`
		ShopRoutingUserID               interface{}   `json:"shop_routing_user_id"`
		IsOrganicProductTaggingEligible bool          `json:"is_organic_product_tagging_eligible"`
		ProductSuggestions              []interface{} `json:"product_suggestions"`
		CanSeeInsightsAsBrand           bool          `json:"can_see_insights_as_brand"`
		MediaType                       int           `json:"media_type"`
		Code                            string        `json:"code"`
		Caption                         interface{}   `json:"caption"`
		SharingFrictionInfo             struct {
			ShouldHaveSharingFriction bool        `json:"should_have_sharing_friction"`
			BloksAppURL               interface{} `json:"bloks_app_url"`
			SharingFrictionPayload    interface{} `json:"sharing_friction_payload"`
		} `json:"sharing_friction_info"`
		HasTranslation                   bool          `json:"has_translation"`
		OriginalMediaHasVisualReplyMedia bool          `json:"original_media_has_visual_reply_media"`
		CoauthorProducers                []interface{} `json:"coauthor_producers"`
		InvitedCoauthorProducers         []interface{} `json:"invited_coauthor_producers"`
		IsInProfileGrid                  bool          `json:"is_in_profile_grid"`
		ProfileGridControlEnabled        bool          `json:"profile_grid_control_enabled"`
		ImageVersions2                   struct {
			Candidates []struct {
				Width        int    `json:"width"`
				Height       int    `json:"height"`
				URL          string `json:"url"`
				ScansProfile string `json:"scans_profile"`
			} `json:"candidates"`
		} `json:"image_versions2"`
		OriginalWidth              int         `json:"original_width"`
		OriginalHeight             int         `json:"original_height"`
		EnableMediaNotesProduction bool        `json:"enable_media_notes_production"`
		ProductType                string      `json:"product_type"`
		IsPaidPartnership          bool        `json:"is_paid_partnership"`
		MusicMetadata              interface{} `json:"music_metadata"`
		OrganicTrackingToken       string      `json:"organic_tracking_token"`
		IgMediaSharingDisabled     bool        `json:"ig_media_sharing_disabled"`
		BoostUnavailableIdentifier interface{} `json:"boost_unavailable_identifier"`
		BoostUnavailableReason     interface{} `json:"boost_unavailable_reason"`
		IsAutoCreated              bool        `json:"is_auto_created"`
		IsCutoutStickerAllowed     bool        `json:"is_cutout_sticker_allowed"`
		Owner                      struct {
			IsPrivate bool   `json:"is_private"`
			Pk        int64  `json:"pk"`
			StrongID  string `json:"strong_id__"`
		} `json:"owner"`
		FbAggregatedLikeCount                                int           `json:"fb_aggregated_like_count"`
		FbAggregatedCommentCount                             int           `json:"fb_aggregated_comment_count"`
		IsTaggedMediaSharedToViewerProfileGrid               bool          `json:"is_tagged_media_shared_to_viewer_profile_grid"`
		ShouldShowAuthorPogForTaggedMediaSharedToProfileGrid bool          `json:"should_show_author_pog_for_tagged_media_shared_to_profile_grid"`
		CollapseComments                                     bool          `json:"collapse_comments"`
		Likers                                               []interface{} `json:"likers"`
		IsOpenToPublicSubmission                             bool          `json:"is_open_to_public_submission"`
		ArchiveStoryDeletionTs                               int           `json:"archive_story_deletion_ts"`
		CanSendPrompt                                        bool          `json:"can_send_prompt"`
		HasSharedToFbDating                                  int           `json:"has_shared_to_fb_dating"`
		IsFirstTake                                          bool          `json:"is_first_take"`
		IsRollcallV2                                         bool          `json:"is_rollcall_v2"`
		SourceType                                           int           `json:"source_type"`
		StoryIsSavedToArchive                                bool          `json:"story_is_saved_to_archive"`
		SupportsReelReactions                                bool          `json:"supports_reel_reactions"`
		CanPlaySpotifyAudio                                  bool          `json:"can_play_spotify_audio"`
		IsFromDiscoverySurface                               bool          `json:"is_from_discovery_surface"`
		IsSuperlative                                        bool          `json:"is_superlative"`
		ShowOneTapFbShareTooltip                             bool          `json:"show_one_tap_fb_share_tooltip"`
		StoryLinkStickers                                    []struct {
			X           float64 `json:"x"`
			Y           float64 `json:"y"`
			Z           int     `json:"z"`
			Width       float64 `json:"width"`
			Height      float64 `json:"height"`
			Rotation    float64 `json:"rotation"`
			IsPinned    int     `json:"is_pinned"`
			IsHidden    int     `json:"is_hidden"`
			IsSticker   int     `json:"is_sticker"`
			IsFbSticker int     `json:"is_fb_sticker"`
			StartTimeMs float64 `json:"start_time_ms"`
			EndTimeMs   float64 `json:"end_time_ms"`
			StoryLink   struct {
				LinkType   string `json:"link_type"`
				ClickID    string `json:"click_id"`
				URL        string `json:"url"`
				LinkTitle  string `json:"link_title"`
				DisplayURL string `json:"display_url"`
			} `json:"story_link"`
		} `json:"story_link_stickers"`
		User struct {
			Pk        int64  `json:"pk"`
			IsPrivate bool   `json:"is_private"`
			StrongID  string `json:"strong_id__"`
		} `json:"user"`
		CanReshare           bool          `json:"can_reshare"`
		CanReply             bool          `json:"can_reply"`
		Viewers              []interface{} `json:"viewers"`
		ViewerCount          int           `json:"viewer_count"`
		FbViewerCount        interface{}   `json:"fb_viewer_count"`
		ViewerCursor         interface{}   `json:"viewer_cursor"`
		TotalViewerCount     int           `json:"total_viewer_count"`
		MultiAuthorReelNames []interface{} `json:"multi_author_reel_names"`
	} `json:"iphone_struct"`
}
