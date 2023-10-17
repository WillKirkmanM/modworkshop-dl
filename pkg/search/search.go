package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type SearchResponse struct {
	Data []struct {
		ID                  int       `json:"id"`
		CategoryID          int       `json:"category_id"`
		GameID              int       `json:"game_id"`
		UserID              int       `json:"user_id"`
		Name                string    `json:"name"`
		Desc                string    `json:"desc"`
		ShortDesc           string    `json:"short_desc"`
		Changelog           string    `json:"changelog"`
		License             string    `json:"license"`
		Instructions        string    `json:"instructions"`
		Visibility          string    `json:"visibility"`
		LegacyBannerURL     string    `json:"legacy_banner_url"`
		Downloads           int       `json:"downloads"`
		Likes               int       `json:"likes"`
		Views               int       `json:"views"`
		Version             string    `json:"version"`
		Donation            string    `json:"donation"`
		Suspended           bool      `json:"suspended"`
		CommentsDisabled    bool      `json:"comments_disabled"`
		Score               string    `json:"score"`
		DailyScore          string    `json:"daily_score"`
		WeeklyScore         string    `json:"weekly_score"`
		BumpedAt            time.Time `json:"bumped_at"`
		PublishedAt         time.Time `json:"published_at"`
		DownloadID          any       `json:"download_id"`
		DownloadType        any       `json:"download_type"`
		LastUserID          int       `json:"last_user_id"`
		HasDownload         bool      `json:"has_download"`
		Approved            bool      `json:"approved"`
		AllowedStorage      any       `json:"allowed_storage"`
		CreatedAt           time.Time `json:"created_at"`
		UpdatedAt           time.Time `json:"updated_at"`
		ThumbnailID         int       `json:"thumbnail_id"`
		BannerID            int       `json:"banner_id"`
		InstructsTemplateID int       `json:"instructs_template_id"`
		User                struct {
			ID               int       `json:"id"`
			Name             string    `json:"name"`
			Ban              any       `json:"ban"`
			GameBan          any       `json:"game_ban"`
			UniqueName       string    `json:"unique_name"`
			CreatedAt        time.Time `json:"created_at"`
			Avatar           string    `json:"avatar"`
			RoleNames        []string  `json:"role_names"`
			Tag              any       `json:"tag"`
			RoleIds          []int     `json:"role_ids"`
			LastOnline       time.Time `json:"last_online"`
			CustomColor      string    `json:"custom_color"`
			HighestRoleOrder any       `json:"highest_role_order"`
			Banner           string    `json:"banner"`
			Bio              string    `json:"bio"`
			Invisible        bool      `json:"invisible"`
			PrivateProfile   bool      `json:"private_profile"`
			CustomTitle      string    `json:"custom_title"`
			DonationURL      string    `json:"donation_url"`
			ShowTag          string    `json:"show_tag"`
			ActiveSupporter  any       `json:"active_supporter"`
			ModsCount        any       `json:"mods_count"`
		} `json:"user"`
		Game struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			ShortName string `json:"short_name"`
		} `json:"game"`
		Category struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"category"`
		Thumbnail struct {
			ID        int       `json:"id"`
			UserID    int       `json:"user_id"`
			ModID     int       `json:"mod_id"`
			HasThumb  bool      `json:"has_thumb"`
			File      string    `json:"file"`
			Type      string    `json:"type"`
			Size      int       `json:"size"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"thumbnail"`
	} `json:"data"`
	Meta struct {
		CurrentPage int `json:"current_page"`
		From        int `json:"from"`
		LastPage    int `json:"last_page"`
		PerPage     int `json:"per_page"`
		To          int `json:"to"`
		Total       int `json:"total"`
	} `json:"meta"`
}

func Search(args []string) {
	query := strings.Join(args, "%20")
	
	searchURL := fmt.Sprintf("https://modworkshop.net/api/mods?limit=10&sort=views&query=%v", query) 
	fmt.Println(searchURL)

	res, err := http.Get(searchURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var searchData SearchResponse

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	json.Unmarshal(responseData, &searchData)

	fmt.Println(string(responseData))

}