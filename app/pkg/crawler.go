package pkg

import (
	"agit-crawler/app/lib"
	"fmt"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type agitInfo struct {
	name     string
	userId   string
	password string
}

type Crawler struct {
	browser  *rod.Browser
	page     *rod.Page
	agitInfo agitInfo
}

// 생성자
func NewCrawler(env lib.Env) *Crawler {
	url := launcher.New().Headless(false).MustLaunch()
	browser := rod.New().ControlURL(url).MustConnect()
	page := browser.MustPage("")

	return &Crawler{
		browser:  browser,
		page:     page,
		agitInfo: agitInfo{name: env.Agit.Name, userId: env.Agit.UserId, password: env.Agit.Password},
	}
}

// 브라우저 닫기
func (c *Crawler) Close() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Close():", r)
		}
	}()
	c.browser.MustClose()
}

// 로그인 처리
func (c *Crawler) LoginAgit() {
	loginURL := "https://appm.agit.io/login"
	c.page.Navigate(loginURL)
	c.page.MustWaitLoad()
	c.page.MustElement("#email").MustInput(c.agitInfo.userId)
	c.page.MustElement("#password").MustInput(c.agitInfo.password)
	c.page.MustElement("input[type=submit]").MustClick()
	c.page.MustWaitLoad()
	c.page.MustElement(".navbar__myprofile-agit-id").MustWaitVisible()
	fmt.Println("Login success")
}

// 게시글 가져오기
func (c *Crawler) GetPosts() {
	if c.page == nil {
		fmt.Println("Error: Page is not initialized")
		return
	}

	searchURL := fmt.Sprintf("https://%s.agit.io/search?parent=parent&q=&scope=wall&sort=recent", c.agitInfo.name)
	c.page.Navigate(searchURL)
	c.page.MustWaitLoad()

	// 로그인 페이지로 리디렉션됐는지 확인
	currentURL := c.page.MustInfo().URL
	if currentURL == "https://appm.agit.io/login" || currentURL == fmt.Sprintf("https://%s.agit.io/login", c.agitInfo.name) {
		fmt.Println("Detected logged-out state. Logging in again...")

		c.LoginAgit()

		// 로그인 후 다시 크롤링 페이지로 이동
		c.page.Navigate(searchURL)
		c.page.MustWaitLoad()
	}

	c.page.MustElement(".search-list-page__item").MustWaitVisible()

	// 게시글 정보 파싱
	items := c.page.MustElements(".search-list-page__item")
	for _, item := range items {
		groupTitle := item.MustElement(".wall-message__group-title").MustText()
		fmt.Println("Group Title:", groupTitle)

		authorName := item.MustElement(".wall-message__author-name").MustText()
		fmt.Println("Author Name:", authorName)

		markedContent := item.MustElement(".react-afm").MustText()
		fmt.Println("Marked React AFM:", markedContent)

		btn := item.MustElement(".msg-button.msg-button--clickable")
		href := btn.MustAttribute("href")
		if href != nil {
			fmt.Println("Button href:", *href)
		} else {
			fmt.Println("Button href not found")
		}
		return //수정 필요
	}
}
