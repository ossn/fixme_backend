package models

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
)

// DB is a connection to your database to be used
// throughout your application.
var DB * pop.Connection

func init() {
	var err error
	env := envy.Get("GO_ENV", "development")
	DB, err = pop.Connect(env)
	if err != nil {
		log.Fatal(err)
	}
	pop.Debug = env == "development"

	//--------------------------------------------------------------------------
	db_projects := Projects {}
	err = DB.All( & db_projects)
	if err != nil || len(db_projects) == 0 {
		for _, project := range projects {
			err := DB.Eager().Create( & project)
			if err != nil {
				return
			}
		}
	}
	return
	//DB.Eager().Create(&account)
}

var account = Admin {
	Email: "admin",
	Password: "admin123456789",
}

var projects = Projects {
	{
		Name: "Common Voice",
		FirstColor: "#ABDEF5",
		SecondColor: "#CDCFEE",
		Description: "The Common Voice project is Mozilla’s initiative to help teach machines how real people speak.\nVoice is natural, voice is human. That’s why we’re fascinated with creating usable voice technology for our machines. But to create voice systems, an extremely large amount of voice data is required.\n\nMost of the data used by large companies isn’t available to the majority of people. We think that stifles innovation. So we’ve launched Project Common Voice, a project to help make voice recognition open to everyone.",
		Logo: "https://voice.mozilla.org/img/cv-logo-bw.svg",
		Link: "https://github.com/mozilla/voice-web",
		IsGitHub: true,
	}, {
		Name: "Firefox Focus for Android",
		FirstColor: "#A3007F",
		SecondColor: "#A3007F",
		Description: "Browse like no one’s watching. The new Firefox Focus automatically blocks a wide range of online trackers — from the moment you launch it to the second you leave it. Easily erase your history, passwords and cookies, so you won’t get followed by things like unwanted ads.\n\nFirefox Focus provides automatic ad blocking and tracking protection on an easy-to-use private browser.\n\nGet it on Google Play",
		Logo: "https://www.underconsideration.com/brandnew/archives/firefox_2017_focus.jpg",
		Link: "https://github.com/mozilla-mobile/focus-android",
		IsGitHub: true,
	}, {
		ProjectID: 4983582,
		Name: "Vesta",
		FirstColor: "#A3097F",
		SecondColor: "#A3097F",
		Description: "A Ruby on Rails app to facilitate on-campus housing procedures. Developed for Yale's undergraduate housing process.",
		Logo: "https://www.putnamhousing.com/wp-content/uploads/2017/02/Fair-Housing.png",
		Link: "https://gitlab.com/yale-sdmp/vesta",
		IsGitHub: false,
	}, {
		ProjectID: 4856282,
		Name: "redream",
		FirstColor: "#B3007F",
		SecondColor: "#B3007F",
		Description: "Work In Progress SEGA Dreamcast emulator",
		Logo: "http://www.fasebonus.net/wp-content/uploads/2014/02/wdreamcast-400x300.jpg",
		Link: "https://gitlab.com/inolen/redream",
		IsGitHub: false,
	}, {
		ProjectID: 4339844,
		Name: "freedesktop-sdk",
		FirstColor: "#FF300FF",
		SecondColor: "#B3BB7F",
		Description: "A minimal Linux runtime",
		Logo: "https://gitlab.com/uploads/-/system/project/avatar/4339844/MM_freedesktop2-01__copy_.png",
		Link: "https://gitlab.com/freedesktop-sdk/freedesktop-sdk",
		IsGitHub: false,
	}, {
		Name: "debugger.html",
		FirstColor: "#15638D",
		SecondColor: "#158D63",
		Description: "debugger.html is a hackable debugger for modern times, built from the ground up using React and Redux. It is designed to be approachable, yet powerful.\n\nMozilla created this debugger for use in the Firefox Developer Tools. And we've purposely created this project in GitHub, using modern toolchains.\n\nWe hope to not only to create a great debugger that works with the Firefox and Chrome debugging protocols but develop a broader community that wants to create great tools for the web.",
		Logo: "https://projects.ossn.club/assets/images/firefox-logo.png",
		Link: "https://github.com/firefox-devtools/debugger",
		IsGitHub: true,
	}, {
		Name: "mozillians",
		FirstColor: "#15778F",
		SecondColor: "#25778F",
		Description: "Mozilla community directory -- A centralized directory of all Mozilla contributors!",
		Logo: "https://i1.wp.com/cdn.mozillians.org/media/img/default_avatar.png?ssl=1",
		Link: "https://github.com/mozilla/mozillians",
		IsGitHub: true,
	},
}
