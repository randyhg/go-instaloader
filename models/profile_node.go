package models

type ProfileNode struct {
	AiAgentType           interface{} `json:"ai_agent_type"`
	Biography             string      `json:"biography"`
	BioLinks              []BioLink   `json:"bio_links"`
	FbProfileBiolink      interface{} `json:"fb_profile_biolink"`
	BiographyWithEntities struct {
		RawText  string `json:"raw_text"`
		Entities []struct {
			User struct {
				Username string `json:"username"`
			} `json:"user"`
			Hashtag interface{} `json:"hashtag"`
		} `json:"entities"`
	} `json:"biography_with_entities"`
	BlockedByViewer        bool   `json:"blocked_by_viewer"`
	RestrictedByViewer     bool   `json:"restricted_by_viewer"`
	CountryBlock           bool   `json:"country_block"`
	EimuID                 string `json:"eimu_id"`
	ExternalURL            string `json:"external_url"`
	ExternalURLLinkshimmed string `json:"external_url_linkshimmed"`
	EdgeFollowedBy         struct {
		Count int `json:"count"`
	} `json:"edge_followed_by"`
	Fbid             string `json:"fbid"`
	FollowedByViewer bool   `json:"followed_by_viewer"`
	EdgeFollow       struct {
		Count int `json:"count"`
	} `json:"edge_follow"`
	FollowsViewer         bool        `json:"follows_viewer"`
	FullName              string      `json:"full_name"`
	GroupMetadata         interface{} `json:"group_metadata"`
	HasArEffects          bool        `json:"has_ar_effects"`
	HasClips              bool        `json:"has_clips"`
	HasGuides             bool        `json:"has_guides"`
	HasChaining           bool        `json:"has_chaining"`
	HasChannel            bool        `json:"has_channel"`
	HasBlockedViewer      bool        `json:"has_blocked_viewer"`
	HighlightReelCount    int         `json:"highlight_reel_count"`
	HasRequestedViewer    bool        `json:"has_requested_viewer"`
	HideLikeAndViewCounts bool        `json:"hide_like_and_view_counts"`
	ID                    string      `json:"id"`
	IsBusinessAccount     bool        `json:"is_business_account"`
	IsProfessionalAccount bool        `json:"is_professional_account"`
	IsSupervisionEnabled  bool        `json:"is_supervision_enabled"`
	IsGuardianOfViewer    bool        `json:"is_guardian_of_viewer"`
	IsSupervisedByViewer  bool        `json:"is_supervised_by_viewer"`
	IsSupervisedUser      bool        `json:"is_supervised_user"`
	IsEmbedsDisabled      bool        `json:"is_embeds_disabled"`
	IsJoinedRecently      bool        `json:"is_joined_recently"`
	GuardianID            interface{} `json:"guardian_id"`
	BusinessAddressJSON   string      `json:"business_address_json"`
	BusinessContactMethod string      `json:"business_contact_method"`
	BusinessEmail         interface{} `json:"business_email"`
	BusinessPhoneNumber   interface{} `json:"business_phone_number"`
	BusinessCategoryName  interface{} `json:"business_category_name"`
	OverallCategoryName   interface{} `json:"overall_category_name"`
	CategoryEnum          interface{} `json:"category_enum"`
	CategoryName          string      `json:"category_name"`
	IsPrivate             bool        `json:"is_private"`
	IsVerified            bool        `json:"is_verified"`
	IsVerifiedByMv4B      bool        `json:"is_verified_by_mv4b"`
	IsRegulatedC18        bool        `json:"is_regulated_c18"`
	EdgeMutualFollowedBy  struct {
		Count int `json:"count"`
		Edges []struct {
			Node struct {
				Username string `json:"username"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"edge_mutual_followed_by"`
	PinnedChannelsListCount        int           `json:"pinned_channels_list_count"`
	ProfilePicURL                  string        `json:"profile_pic_url"`
	ProfilePicURLHd                string        `json:"profile_pic_url_hd"`
	RequestedByViewer              bool          `json:"requested_by_viewer"`
	ShouldShowCategory             bool          `json:"should_show_category"`
	ShouldShowPublicContacts       bool          `json:"should_show_public_contacts"`
	ShowAccountTransparencyDetails bool          `json:"show_account_transparency_details"`
	TransparencyLabel              interface{}   `json:"transparency_label"`
	TransparencyProduct            interface{}   `json:"transparency_product"`
	Username                       string        `json:"username"`
	ConnectedFbPage                interface{}   `json:"connected_fb_page"`
	Pronouns                       []interface{} `json:"pronouns"`
	EdgeOwnerToTimelineMedia       struct {
		Count    int `json:"count"`
		PageInfo struct {
			HasNextPage bool   `json:"has_next_page"`
			EndCursor   string `json:"end_cursor"`
		} `json:"page_info"`
		Edges []interface{} `json:"edges"`
	} `json:"edge_owner_to_timeline_media"`
}

type BioLink struct {
	Title    string `json:"title"`
	LynxURL  string `json:"lynx_url"`
	URL      string `json:"url"`
	LinkType string `json:"link_type"`
}
