SHELL = /bin/bash
DESTDIR = ./dist

RELEASE_BASENAME = Spotlight.Search
BUNDLE_NAME = Spotlight\ Search
BUNDLE_VERSION = $(shell cat VERSION)
BUNDLE_IDENTIFIER = nbjahan.launchbar.spotlight
BUNDLE_ICON = Spotlight.icns
AUTHOR = nbjahan
TWITTER = @nbjahan
SLUG = launchbar-spotlight
WEBSITE = http://github.com/nbjahan/$(SLUG)
SCRIPT_NAME = spotlight

LBACTION_PATH = $(DESTDIR)/$(RELEASE_BASENAME).lbaction
RELEASE_FILENAME = $(RELEASE_BASENAME)-$(BUNDLE_VERSION).lbaction
LDFLAGS=

UPDATE_LINK = https://raw.githubusercontent.com/nbjahan/$(SLUG)/master/src/Info.plist
DOWNLOAD_LINK = https://github.com/nbjahan/$(SLUG)/releases/download/v$(BUNDLE_VERSION)/$(RELEASE_FILENAME)

all:
	@$(RM) -rf $(DESTDIR)

	go install github.com/nbjahan/go-launchbar

	@install -d ${LBACTION_PATH}/Contents/{Resources,Scripts}
	@plutil -replace CFBundleName -string $(BUNDLE_NAME) $(PWD)/src/Info.plist
	@plutil -replace CFBundleVersion -string $(BUNDLE_VERSION) $(PWD)/src/Info.plist
	@plutil -replace CFBundleIdentifier -string $(BUNDLE_IDENTIFIER) $(PWD)/src/Info.plist
	@plutil -replace CFBundleIconFile -string $(BUNDLE_ICON) $(PWD)/src/Info.plist
	@plutil -replace LBDescription.LBAuthor -string $(AUTHOR) $(PWD)/src/Info.plist
	@plutil -replace LBDescription.LBTwitter -string $(TWITTER) $(PWD)/src/Info.plist
	@plutil -replace LBDescription.LBWebsite -string $(WEBSITE) $(PWD)/src/Info.plist
	@plutil -replace LBDescription.LBUpdate -string $(UPDATE_LINK) $(PWD)/src/Info.plist
	@plutil -replace LBDescription.LBDownload -string $(DOWNLOAD_LINK) $(PWD)/src/Info.plist
	@plutil -replace LBScripts.LBDefaultScript.LBScriptName -string $(SCRIPT_NAME) $(PWD)/src/Info.plist
	@install -pm 0644 ./src/Info.plist $(LBACTION_PATH)/Contents/
	go build $(LDFLAGS) -o $(LBACTION_PATH)/Contents/Scripts/$(SCRIPT_NAME) ./src
	-@cp -r ./resources/* $(LBACTION_PATH)/Contents/Resources/

	@echo "Refreshing the LaunchBar"
	@osascript -e 'run script "tell application \"LaunchBar\" \n repeat with rule in indexing rules \n if name of rule is \"Actions\" then \n update rule \n exit repeat \n end if \n end repeat \n activate \n end tell"'

	@echo "Making a release"
	@install -d $(DESTDIR)/release
	@ditto -ck --keepParent $(LBACTION_PATH)/ $(DESTDIR)/release/$(RELEASE_FILENAME)

dev: LDFLAGS := -ldflags "-X main.InDev true"
dev: all

release: all

.PHONY: all dev
