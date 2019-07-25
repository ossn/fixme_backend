package grifts

import (
	"github.com/gobuffalo/nulls"
	"github.com/markbates/grift/grift"
	"github.com/ossn/fixme_backend/models"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		for _, repository := range repositories {
			err := models.DB.Eager().Create(&repository)
			if err != nil {
				return err
			}
		}
		return nil
	})

})

var repositories = models.Repositories{
	{
		Project: models.Project{
			DisplayName:   "Common Voice",
			FirstColor:    "#ABDEF5",
			SecondColor:   nulls.String{String: "#CDCFEE", Valid: true},
			Description:   "The Common Voice project is Mozilla’s initiative to help teach machines how real people speak.\nVoice is natural, voice is human. That’s why we’re fascinated with creating usable voice technology for our machines. But to create voice systems, an extremely large amount of voice data is required.\n\nMost of the data used by large companies isn’t available to the majority of people. We think that stifles innovation. So we’ve launched Project Common Voice, a project to help make voice recognition open to everyone.",
			SetupDuration: nulls.String{String: "1'", Valid: false},
			Logo:          "https://voice.mozilla.org/img/cv-logo-bw.svg",
			Link:          "https://voice.mozilla.org/en",
			Tags:          []string{"nodejs", "npm", "ffmpeg", "docker", "yarn"},
			IssuesCount:   33,
		},
		RepositoryUrl: "https://github.com/mozilla/voice-web",
	},
	{
		Project: models.Project{
			DisplayName:   "A-Frame",
			FirstColor:    "#24CAFF",
			SecondColor:   nulls.String{String: "#24CAFF", Valid: true},
			Description:   "A web framework for building virtual reality experiences. Make WebVR with HTML and Entity-Component. Works on Vive, Rift, Daydream, GearVR, desktop\n\nVirtual Reality Made Simple: A-Frame handles the 3D and WebVR boilerplate required to get running across platforms including mobile, desktop, Vive, and Rift just by dropping in <a-scene>.",
			Logo:          "https://www.drupal.org/files/project-images/download_6.png",
			Link:          "https://github.com/aframevr/aframe/",
			SetupDuration: nulls.String{String: "1'", Valid: true},
			Tags:          []string{"javascript", "vr", "webvr", "threejs", "html", "oculus"},
			IssuesCount:   238,
		},
		RepositoryUrl: "https://github.com/aframevr/aframe/",
	},
	{
		Project: models.Project{
			DisplayName:   "Firefox Focus for Android",
			FirstColor:    "#A3007F",
			SecondColor:   nulls.String{String: "#A3007F", Valid: true},
			Description:   "Browse like no one’s watching. The new Firefox Focus automatically blocks a wide range of online trackers — from the moment you launch it to the second you leave it. Easily erase your history, passwords and cookies, so you won’t get followed by things like unwanted ads.\n\nFirefox Focus provides automatic ad blocking and tracking protection on an easy-to-use private browser.\n\nGet it on Google Play",
			Logo:          "https://www.underconsideration.com/brandnew/archives/firefox_2017_focus.jpg",
			Link:          "https://github.com/mozilla-mobile/focus-android",
			SetupDuration: nulls.String{String: "17'", Valid: true},
			Tags:          []string{"javascript", "html", "mobile", "browser"},
			IssuesCount:   238,
		},
		RepositoryUrl: "https://github.com/mozilla-mobile/focus-android",
	},
	{
		Project: models.Project{
			DisplayName:   "debugger.html",
			FirstColor:    "#15638D",
			SecondColor:   nulls.String{String: "#158D63", Valid: true},
			Description:   "debugger.html is a hackable debugger for modern times, built from the ground up using React and Redux. It is designed to be approachable, yet powerful.\n\nMozilla created this debugger for use in the Firefox Developer Tools. And we've purposely created this project in GitHub, using modern toolchains.\n\nWe hope to not only to create a great debugger that works with the Firefox and Chrome debugging protocols but develop a broader community that wants to create great tools for the web.",
			Logo:          "https://projects.ossn.club/assets/images/firefox-logo.png",
			Link:          "https://firefox-dev.tools/debugger.html/#getting-involved",
			SetupDuration: nulls.String{String: "5'", Valid: true},
			Tags:          []string{"nodejs", "npm", "react", "redux", "javascript"},
			IssuesCount:   7,
		},
		RepositoryUrl: "https://github.com/firefox-devtools/debugger",
	}}
