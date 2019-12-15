package models

import (
	"log"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
)

// DB is a connection to your database to be used
// throughout your application.
var DB *pop.Connection

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"

	//--------------------------------------------------------------------------
	for _, repository := range repositories {
		err := DB.Eager().Create(&repository)
		if err != nil {
			return
		}
	}
	DB.Eager().Create(&account)
}

	var account = Admin{
		Email: "admin",
		Password: "admin123456789",
	}

	var repositories = Repositories{
		{
			Project: Project{
				DisplayName: "Common Voice",
				FirstColor:  "#ABDEF5",
				SecondColor: nulls.String{String: "#CDCFEE", Valid: true},
				Description: "The Common Voice project is Mozilla’s initiative to help teach machines how real people speak.\nVoice is natural, voice is human. That’s why we’re fascinated with creating usable voice technology for our machines. But to create voice systems, an extremely large amount of voice data is required.\n\nMost of the data used by large companies isn’t available to the majority of people. We think that stifles innovation. So we’ve launched Project Common Voice, a project to help make voice recognition open to everyone.",
				Logo:        "https://voice.mozilla.org/img/cv-logo-bw.svg",
				Link:        "https://voice.mozilla.org/en",
				Tags:        []string{"nodejs", "npm", "ffmpeg", "docker", "yarn"},
				IssuesCount: 33,
				IsGitHub:		 true,
			},
			RepositoryUrl: "https://github.com/mozilla/voice-web",
			IsGitHub:		 true,
		},
		{
			Project: Project{
				DisplayName:   "A-Frame",
				FirstColor:    "#24CAFF",
				SecondColor:   nulls.String{String: "#24CAFF", Valid: true},
				Description:   "A web framework for building virtual reality experiences. Make WebVR with HTML and Entity-Component. Works on Vive, Rift, Daydream, GearVR, desktop\n\nVirtual Reality Made Simple: A-Frame handles the 3D and WebVR boilerplate required to get running across platforms including mobile, desktop, Vive, and Rift just by dropping in <a-scene>.",
				Logo:          "https://www.drupal.org/files/project-images/download_6.png",
				Link:          "https://github.com/aframevr/aframe/",
				SetupDuration: nulls.String{String: "1'", Valid: true},
				Tags:          []string{"javascript", "vr", "webvr", "threejs", "html", "oculus"},
				IssuesCount:   238,
				IsGitHub:		 true,
			},
			RepositoryUrl: "https://github.com/aframevr/aframe/",
			IsGitHub:		 true,
		},
		{
			Project: Project{
				DisplayName:   "Firefox Focus for Android",
				FirstColor:    "#A3007F",
				SecondColor:   nulls.String{String: "#A3007F", Valid: true},
				Description:   "Browse like no one’s watching. The new Firefox Focus automatically blocks a wide range of online trackers — from the moment you launch it to the second you leave it. Easily erase your history, passwords and cookies, so you won’t get followed by things like unwanted ads.\n\nFirefox Focus provides automatic ad blocking and tracking protection on an easy-to-use private browser.\n\nGet it on Google Play",
				Logo:          "https://www.underconsideration.com/brandnew/archives/firefox_2017_focus.jpg",
				Link:          "https://github.com/mozilla-mobile/focus-android",
				SetupDuration: nulls.String{String: "17'", Valid: true},
				Tags:          []string{"javascript", "html", "mobile", "browser"},
				IssuesCount:   238,
				IsGitHub:		 true,
			},
			RepositoryUrl: "https://github.com/mozilla-mobile/focus-android",
			IsGitHub:		 true,
		},
//-------------------------------------------------------------------
		{
			Project: Project{
				DisplayName:   "Vesta",
				FirstColor:    "#A3097F",
				SecondColor:   nulls.String{String: "#A3097F", Valid: true},
				Description:   "A Ruby on Rails app to facilitate on-campus housing procedures. Developed for Yale's undergraduate housing process.",
				Logo:          "https://cwrc.ca/sites/default/files/road-sign-361514_960_720.png",
				Link:          "https://gitlab.com/yale-sdmp/vesta",
				SetupDuration: nulls.String{String: "107'", Valid: true},
				Tags:          []string{"test", "html2121", "mo32bile", "brows99er"},
				IssuesCount:   238,
				IsGitHub:		 	 false,
			},
			RepositoryUrl: "https://gitlab.com/yale-sdmp/vesta",
			IsGitHub:		 false,
		},
		{
			Project: Project{
				DisplayName:   "redream",
				FirstColor:    "#B3007F",
				SecondColor:   nulls.String{String: "#B3007F", Valid: true},
				Description:   "Work In Progress SEGA Dreamcast emulator",
				Logo:          "https://lh3.googleusercontent.com/dOcWhgQtj6fiPDQCjbzkKpq4jI0xV82Z-0UeJuUI4yqZICiQ8Avf5h1lFRIUgHgdpHI=w412-h220-rw",
				Link:          "https://gitlab.com/inolen/redream",
				SetupDuration: nulls.String{String: "1777'", Valid: true},
				Tags:          []string{"html", "sega", "emulator"},
				IssuesCount:   28,
				IsGitHub:		 	 false,
			},
			RepositoryUrl: "https://gitlab.com/inolen/redream",
			IsGitHub:		 false,
		},
//-------------------------------------------------------------------
		{
			Project: Project{
				DisplayName:   "debugger.html",
				FirstColor:    "#15638D",
				SecondColor:   nulls.String{String: "#158D63", Valid: true},
				Description:   "debugger.html is a hackable debugger for modern times, built from the ground up using React and Redux. It is designed to be approachable, yet powerful.\n\nMozilla created this debugger for use in the Firefox Developer Tools. And we've purposely created this project in GitHub, using modern toolchains.\n\nWe hope to not only to create a great debugger that works with the Firefox and Chrome debugging protocols but develop a broader community that wants to create great tools for the web.",
				Logo:          "https://projects.ossn.club/assets/images/firefox-logo.png",
				Link:          "https://firefox-dev.tools/debugger.html/#getting-involved",
				SetupDuration: nulls.String{String: "5'", Valid: true},
				Tags:          []string{"nodejs", "npm", "react", "redux", "javascript"},
				IssuesCount:   7,
				IsGitHub:		 true,
			},
			RepositoryUrl: "https://github.com/firefox-devtools/debugger",
			IsGitHub:		 true,
		}}
