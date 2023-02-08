//Generation random, valid HTTP User-Agent header and web navgator JS object.

//Functions:
//* GenerateUserAgent: generates User-Agent HTTP header
//* GenerateNavigator:  generates web navigator's config
//* GenerateNavigatorJS:  generates web navigator's config with keys
//   identical keys used in navigator object
/*
FIXME:
* add Edge, Safari and Opera support
* add random config i.e. windows is more common than linux
*/
//Specs:
//* https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent/Firefox
//* http://msdn.microsoft.com/en-us/library/ms537503(VS.85).aspx
//* https://developer.chrome.com/multidevice/user-agent
//* http://www.javascriptkit.com/javatutors/navigator.shtml

//Release history:
//* https://en.wikipedia.org/wiki/Firefox_release_history
//* https://en.wikipedia.org/wiki/Google_Chrome_release_history
//* https://en.wikipedia.org/wiki/Internet_Explorer_version_history
//* https://en.wikipedia.org/wiki/Android_version_history

//Lists of user agents:
//* http://www.useragentstring.com/
//* http://www.user-agents.org/
//* http://www.webapps-online.com/online-tools/user-agent-strings

package useragent

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var (
	DEVICE_TYPE_OS = map[string][]string{
		"desktop":    {"win", "mac", "linux"},
		"smartphone": {"android"},
		"tablet":     {"android"},
	}
	OS_DEVICE_TYPE = map[string][]string{
		"win":     {"desktop"},
		"linux":   {"desktop"},
		"mac":     {"desktop"},
		"android": {"smartphone", "tablet"},
	}

	DEVICE_TYPE_NAVIGATOR = map[string][]string{
		"desktop":    {"chrome", "firefox", "ie"},
		"smartphone": {"firefox", "chrome"},
		"tablet":     {"firefox", "chrome"},
	}
	NAVIGATOR_DEVICE_TYPE = map[string][]string{
		"ie":      {"desktop"},
		"chrome":  {"desktop", "smartphone", "tablet"},
		"firefox": {"desktop", "smartphone", "tablet"},
	}

	OS_PLATFORM = map[string][]string{
		"win": {
			// Windows XP
			"Windows NT 5.1",
			// Windows 7
			"Windows NT 6.1",
			// Windows 8
			"Windows NT 6.2",
			// Windows 8.1
			"Windows NT 6.3",
			// Windows 10
			"Windows NT 10.0",
		},
		"mac": {
			"Macintosh; Intel Mac OS X 10.8",
			"Macintosh; Intel Mac OS X 10.9",
			"Macintosh; Intel Mac OS X 10.10",
			"Macintosh; Intel Mac OS X 10.11",
			"Macintosh; Intel Mac OS X 10.12",
		},
		"linux": {
			"X11; Linux",
			"X11; Ubuntu; Linux",
		},
		"android": {
			// 2013-10-31
			"Android 4.4",
			// 2013-12-05
			"Android 4.4.1",
			// 2013-12-09
			"Android 4.4.2",
			// 2014-06-02
			"Android 4.4.3",
			// 2014-06-19
			"Android 4.4.4",
			// 2014-11-12
			"Android 5.0",
			// 2014-12-02
			"Android 5.0.1",
			// 2014-12-19
			"Android 5.0.2",
			// 2015-03-09
			"Android 5.1",
			// 2015-04-21
			"Android 5.1.1",
			// 2015-10-05
			"Android 6.0",
			// 2015-12-07
			"Android 6.0.1",
			// 2016-08-22
			"Android 7.0",
			// 2016-10-04
			"Android 7.1",
			// 2016-12-05
			"Android 7.1.1",
		},
	}

	OS_CPU = map[string][]string{
		"win": {
			//32bit
			"",
			// 64bit
			"Win64; x64",
			// 32bit process on 64bit system
			"WOW64",
		},
		"linux": {
			// 32bit
			"i686",
			// 64bit
			"x86_64",
			// 32bit process on 64bit system
			"i686 on x86_64",
		},
		"mac": {""},
		"android": {
			// 32bit
			"armv7l",
			// 64bit
			"armv8l",
		},
	}
	OS_NAVIGATOR = map[string][]string{
		"win":     {"chrome", "firefox", "ie"},
		"mac":     {"firefox", "chrome"},
		"linux":   {"chrome", "firefox"},
		"android": {"firefox", "chrome"},
	}
	NAVIGATOR_OS = map[string][]string{
		"chrome":  {"win", "linux", "mac", "android"},
		"firefox": {"win", "linux", "mac", "android"},
		"ie":      {"win"},
	}
	MACOSX_CHROME_BUILD_RANGE = map[string][]int{
		// https://en.wikipedia.org/wiki/MacOS#Release_history
		"10.8":  {0, 8},
		"10.9":  {0, 5},
		"10.10": {0, 5},
		"10.11": {0, 6},
		"10.12": {0, 2},
	}

	FIREFOX_VERSION = []FirefoxVersion{
		{"45.0", time.Date(2016, 3, 8, 0, 0, 0, 0, time.UTC)},
		{"46.0", time.Date(2016, 4, 26, 0, 0, 0, 0, time.UTC)},
		{"47.0", time.Date(2016, 6, 7, 0, 0, 0, 0, time.UTC)},
		{"48.0", time.Date(2016, 8, 2, 0, 0, 0, 0, time.UTC)},
		{"49.0", time.Date(2016, 9, 20, 0, 0, 0, 0, time.UTC)},
		{"50.0", time.Date(2016, 11, 15, 0, 0, 0, 0, time.UTC)},
		{"51.0", time.Date(2017, 1, 24, 0, 0, 0, 0, time.UTC)},
	}

	// Top chrome builds from website access log
	// for september, october 2020
	CHROME_BUILD = []string{
		"80.0.3987.132",
		"80.0.3987.149",
		"80.0.3987.99",
		"81.0.4044.117",
		"81.0.4044.138",
		"83.0.4103.101",
		"83.0.4103.106",
		"83.0.4103.96",
		"84.0.4147.105",
		"84.0.4147.111",
		"84.0.4147.125",
		"84.0.4147.135",
		"84.0.4147.89",
		"85.0.4183.101",
		"85.0.4183.102",
		"85.0.4183.120",
		"85.0.4183.121",
		"85.0.4183.127",
		"85.0.4183.81",
		"85.0.4183.83",
		"86.0.4240.110",
		"86.0.4240.111",
		"86.0.4240.114",
		"86.0.4240.183",
		"86.0.4240.185",
		"86.0.4240.75",
		"86.0.4240.78",
		"86.0.4240.80",
		"86.0.4240.96",
		"86.0.4240.99",
	}

	// (numeric ver, string ver, trident ver)
	IE_VERSION = []IEVersion{
		//2009
		{8, "MSIE 8.0", "4.0"},
		//2011
		{9, "MSIE 9.0", "5.0"},
		//2012
		{10, "MSIE 10.0", "6.0"},
		//2013
		{11, "MSIE 11.0", "7.0"},
	}

	USERAGENTTEMPLATE = map[string]any{
		"firefox":           `Mozilla/5.0 ({{index .System "ua_platform"}}; rv:{{index .App "build_version"}}) Gecko/{{index .App "geckotrail"}} Firefox/{{index .App "build_version"}}`,
		"chrome":            `Mozilla/5.0 ({{index .System "ua_platform"}}) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/{{index .App "build_version"}} Safari/537.36`,
		"chrome_smartphone": `Mozilla/5.0 ({{index .System "ua_platform"}}) AppleWebKit/537.36	(KHTML, like Gecko) Chrome/{{index .App "build_version"}} Mobile Safari/537.36`,
		"chrome_tablet":     `Mozilla/5.0 ({{index .System "ua_platform"}}) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/{{index .App "build_version"}} Safari/537.36`,
		"ie_less_11":        `Mozilla/5.0 (compatible; {{index .App "build_version"}}; {{index .System "ua_platform"}}; Trident/{{index .App "trident_version"}})`,
		"ie_11":             `Mozilla/5.0 ({{index .System "ua_platform"}}; Trident/{{index .App "trident_version"}}; rv:11.0) like Gecko`,
	}
)

type (
	FirefoxVersion struct {
		Version string
		Date    time.Time
	}

	IEVersion struct {
		NumericVersion int
		StringVersion  string
		TridentVersion string
	}

	DevIDs []string

	UserAgentConfig struct {
		//OS limit list of os for generation
		//OS can be string or array
		//optional:default
		OS any
		//Navigator is the limit list of browser engines for generation
		//Navigator can be a string or a list
		//Default:Desktop
		//Optional
		Navigator any
		//DeviceType limits possible oses by device type
		//DeviceType is a list, possible values:"desktop", "smartphone", "tablet", "all"
		DeviceType []string
		//Platform limits possible platforms by platform
		//Default:""
		//Optional
		Platform []string
	}

	uatmpl struct {
		System map[string]string
		App    map[string]string
	}
)

func getFirefoxBuild() (string, string, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(FIREFOX_VERSION))))
	if err != nil {
		return "", "", err
	}
	fxvs := FIREFOX_VERSION[nBig.Int64()]
	//var index int
	//var dateto time.Time
	build_ver, date_from := fxvs.Version, fxvs.Version
	return build_ver, date_from, nil //build_rnd_time.strftime("%Y%m%d%H%M%S")
}

func getChromeBuild() (string, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(CHROME_BUILD))))
	if err != nil {
		return "", err
	}
	return CHROME_BUILD[nBig.Int64()], nil
}

func getIEBuild() (IEVersion, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(IE_VERSION))))
	if err != nil {
		return IEVersion{}, err
	}
	return IE_VERSION[nBig.Int64()], nil

}

//Fix chrome version on mac OS.
//Chrome on Mac OS adds minor version number and uses underscores instead
//of dots. E.g. platform for Firefox will be: 'Intel Mac OS X 10.11'
//but for Chrome it will be 'Intel Mac OS X 10_11_6'.
//param platform: - string like "Macintosh; Intel Mac OS X 10.8"
//return: platform with version number including minor number and formatted
//    with underscores, e.g. "Macintosh; Intel Mac OS X 10_8_2"

func fixChromeMacPlatform(platform string) (string, error) {
	ver := strings.Split(platform, "OS X ")[1]
	build_range := (MACOSX_CHROME_BUILD_RANGE[ver])
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(build_range[1])))
	if err != nil {
		return "", err
	}
	build := int(nBig.Int64())
	mac_ver := strings.Replace(ver, ".", "_", -1) + "_" + strconv.Itoa(build)
	return fmt.Sprintf("Macintosh; Intel Mac OS X %s", mac_ver), nil
}

//Build random platform and oscpu components for given parameters.
//Returns dict {platform_version, platform, ua_platform, oscpu}
//platform_version is OS name used in different places
//ua_platform goes to navigator.platform
//platform is used in building navigator.userAgent
//oscpu goes to navigator.oscpu

func buildSystemComponents(deviceType, OSID, navigatorID string) (map[string]string, error) {
	if !contains([]string{"win", "linux", "mac", "android"}, OSID) {
		return nil, errors.New("Invalid platform")
	}
	var platform string
	if OSID == "win" {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(OS_PLATFORM["win"]))))
		if err != nil {
			return nil, err
		}
		platform_version := OS_PLATFORM["win"][nBig.Int64()]
		nBig, err = rand.Int(rand.Reader, big.NewInt(int64(len(OS_CPU["win"]))))
		if err != nil {
			return nil, err
		}
		cpu := OS_CPU["win"][nBig.Int64()]
		if cpu != "" {
			platform = fmt.Sprintf("%s %s", platform_version, cpu)
		} else {
			platform = platform_version
		}
		return map[string]string{
			"platform_version": platform_version,
			"platform":         platform,
			"ua_platform":      platform,
			"oscpu":            platform,
		}, nil
	}
	if OSID == "linux" {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(OS_CPU["linux"]))))
		if err != nil {
			return nil, err
		}
		cpu := OS_CPU["linux"][nBig.Int64()]
		nBig, err = rand.Int(rand.Reader, big.NewInt(int64(len(OS_PLATFORM["linux"]))))
		if err != nil {
			return nil, err
		}
		platform_version := OS_PLATFORM["linux"][nBig.Int64()]
		platform := fmt.Sprintf("%s %s", platform_version, cpu)
		return map[string]string{
			"platform_version": platform_version,
			"platform":         platform,
			"ua_platform":      platform,
			"oscpu":            fmt.Sprintf("Linux %s", cpu),
		}, nil
	}
	if OSID == "mac" {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(OS_PLATFORM["mac"]))))
		if err != nil {
			return nil, err
		}
		platform_version := OS_PLATFORM["mac"][nBig.Int64()]
		platform := platform_version
		if navigatorID == "chrome" {
			platform, err = fixChromeMacPlatform(platform)
			if err != nil {
				return nil, err
			}
		}
		return map[string]string{
			"platform_version": platform_version,
			"platform":         "MacIntel",
			"ua_platform":      platform,
			"oscpu": fmt.Sprintf("Intel Mac OS X %s",
				strings.Split(platform, " ")[len(strings.Split(platform, " "))-1]),
		}, nil
	}

	// OSID could be only "android" here
	if !contains([]string{"firefox", "chrome"}, navigatorID) {
		return nil, errors.New("assertion error")
	}

	if !contains([]string{"smartphone", "tablet"}, deviceType) {
		return nil, errors.New("assertion error")
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(OS_PLATFORM["android"]))))
	if err != nil {
		return nil, err
	}
	platform_version := OS_PLATFORM["android"][nBig.Int64()]
	var ua_platform, oscpu string
	if navigatorID == "firefox" {
		if deviceType == "smartphone" {
			ua_platform = fmt.Sprintf("%s; Mobile", platform_version)
		} else if deviceType == "tablet" {
			ua_platform = fmt.Sprintf("%s; Tablet", platform_version)
		}
	} else if navigatorID == "chrome" {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(OS_PLATFORM["android"]))))
		if err != nil {
			return nil, err
		}
		platform_version = OS_PLATFORM["android"][nBig.Int64()]
		nBig, err = rand.Int(rand.Reader, big.NewInt(int64(len(SMARTPHONE_DEV_IDS))))
		if err != nil {
			return nil, err
		}
		device_id := SMARTPHONE_DEV_IDS[nBig.Int64()]
		ua_platform = fmt.Sprintf("Linux; %s; %s", platform_version, device_id)
		nBig, err = rand.Int(rand.Reader, big.NewInt(int64(len(OS_CPU["android"]))))
		if err != nil {
			return nil, err
		}
		oscpu = fmt.Sprintf("Linux %s", OS_CPU["android"][nBig.Int64()])
	}
	return map[string]string{
		"platform_version": platform_version,
		"ua_platform":      ua_platform,
		"platform":         oscpu,
		"oscpu":            oscpu,
	}, nil
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

//Build app features for given os and navigator.
//Returns dict {name, product_sub, vendor, build_version, build_id}

func buildAppComponents(OSID, navigatorID string) (map[string]string, error) {
	if !contains([]string{"firefox", "chrome", "ie"}, navigatorID) {
		return nil, errors.New("invalid browser")
	}
	if navigatorID == "firefox" {
		//fxbuild, err :=
		build_version, build_id, err := getFirefoxBuild()
		if err != nil {
			return nil, err
		}
		var geckotrail string
		if contains([]string{"win", "linux", "mac"}, OSID) {
			geckotrail = "20100101"
		} else {
			geckotrail = build_version
		}
		return map[string]string{
			"name":          "Netscape",
			"product_sub":   "20100101",
			"vendor":        "",
			"build_version": build_version,
			"build_id":      build_id,
			"geckotrail":    geckotrail,
		}, nil
	}
	if navigatorID == "chrome" {
		chromebuild, err := getChromeBuild()
		if err != nil {
			return nil, err
		}
		return map[string]string{
			"name":          "Netscape",
			"product_sub":   "20030107",
			"vendor":        "Google Inc.",
			"build_version": chromebuild,
			"build_id":      "",
		}, nil
	}
	// navigator_id could be only "ie" here
	iebuild, err := getIEBuild()
	if err != nil {
		return nil, err
	}
	num_ver, build_version, trident_version := iebuild.NumericVersion, iebuild.StringVersion, iebuild.TridentVersion
	var app_name string
	if num_ver >= 11 {
		app_name = "Netscape"
	} else {
		app_name = "Microsoft Internet Explorer"
	}
	return map[string]string{
		"name":            app_name,
		"product_sub":     "",
		"vendor":          "",
		"build_version":   build_version,
		"build_id":        "",
		"trident_version": trident_version,
	}, nil
}

//Generate something.
//Long uninformative description: Generate possible choices for the
//option `opt_name` limited to `opt_value` value with default value
//as `default_value`

func getOptionChoices(opt_name string, opt_value []string, default_value, all_choices []string) []string {
	var choices []string
	choices = opt_value
	if len(opt_value) == 0 {
		choices = default_value
	}
	if contains(opt_value, "all") {
		choices = all_choices
	}
	for _, item := range choices {
		if !contains(all_choices, item) {
			panic(fmt.Sprintf("Choices of option %s contains invalid item: %s", opt_name, item))
		}
	}
	return choices
}

// Select one item from all possible combinations of (device, os, navigator) items.
func pickConfigIDs(cfg *UserAgentConfig) (string, string, string, error) {
	var (
		defaultDevTypes  []string
		deviceTypeOSKeys []string
		osNavigatorKeys  []string
		navigatorOSkeys  []string
		variants         [][]string
		deviceType       []string
		OS               []string
		navigator        []string
		osPlatformKeys   []string
	)
	for k := range DEVICE_TYPE_OS {
		deviceTypeOSKeys = append(deviceTypeOSKeys, k)
	}

	for j := range OS_NAVIGATOR {
		osNavigatorKeys = append(osNavigatorKeys, j)
	}
	for l := range NAVIGATOR_OS {
		navigatorOSkeys = append(navigatorOSkeys, l)
	}
	for k := range OS_PLATFORM {
		osPlatformKeys = append(osPlatformKeys, k)
	}
	if cfg.OS == nil {
		defaultDevTypes = []string{"desktop"}
	} else {
		defaultDevTypes = deviceTypeOSKeys

	}
	if cfg.DeviceType == nil {
		deviceType = nil
	}

	if cfg.OS == nil {
		OS = nil
	}
	if cfg.Navigator == nil {
		navigator = nil
	}
	devTypeChoices := getOptionChoices("device_type", deviceType, defaultDevTypes, deviceTypeOSKeys)
	osChoices := getOptionChoices("os", OS, osNavigatorKeys, osNavigatorKeys)
	navChoices := getOptionChoices("navigator", navigator, navigatorOSkeys, navigatorOSkeys)
	n := product(devTypeChoices, osChoices, navChoices)
	i := len(devTypeChoices) * len(osChoices) * len(navChoices)
	fmt.Println(i)
	for i >= 0 {
		i--
		prod := n()
		if len(prod) != 0 {
			iter_dev, iter_os, iter_nav := prod[0], prod[1], prod[2]
			fmt.Println(iter_dev, iter_os, iter_nav)
			if contains(DEVICE_TYPE_OS[iter_dev], iter_os) && contains(DEVICE_TYPE_NAVIGATOR[iter_dev], iter_nav) && contains(OS_NAVIGATOR[iter_os], iter_nav) {
				variants = append(variants, []string{iter_dev, iter_os, iter_nav})
			}
		}
	}
	if len(variants) == 0 {
		panic("Options device_type, os and navigator conflicts with each other")
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(variants))))
	if err != nil {
		panic(err)
	}
	//fmt.Prinln
	device_type, os_id, navigator_id := variants[nBig.Int64()][0], variants[nBig.Int64()][1], variants[nBig.Int64()][2]

	if !contains(osPlatformKeys, os_id) {
		panic("error")
	}
	if !contains(navigatorOSkeys, navigator_id) {
		panic("error")
	}
	if !contains(deviceTypeOSKeys, device_type) {
		panic("error")
	}

	return device_type, os_id, navigator_id, nil

}
func chooseUATemplate(device_type, navigator_id string, app map[string]string) any {
	tpl_name := navigator_id
	if navigator_id == "ie" {
		if app["build_version"] == "MSIE 11.0" {
			tpl_name = "ie_11"
		} else {
			tpl_name = "ie_less_11"
		}
	}
	if navigator_id == "chrome" {
		if device_type == "smartphone" {
			tpl_name = "chrome_smartphone"
		}
		if device_type == "tablet" {
			tpl_name = "chrome_tablet"
		}
	}
	return USERAGENTTEMPLATE[tpl_name]
}

func build_navigator_app_version(OSID, navigatorID, platformVersion, userAgent string) string {
	if navigatorID == "firefox" {
		if OSID == "android" {
			return fmt.Sprintf("5.0 (%s)", platformVersion)
		}
		osToken := map[string]string{
			"win":   "Windows",
			"mac":   "Macintosh",
			"linux": "X11",
		}[OSID]
		return fmt.Sprintf("5.0 (%s)", osToken)
	}
	// here navigator_id could be only "chrome" and "ie"
	if !strings.HasPrefix(userAgent, "Mozilla/") {
		panic("assertion error")
	}
	return strings.Split(userAgent, "Mozilla/")[1]
}

// Generate web navigator's config
func generateNavigator(config *UserAgentConfig) map[string]string {
	device_type, os_id, navigator_id, _ := pickConfigIDs(config)
	fmt.Println("device type", device_type, os_id, navigator_id)
	system, _ := buildSystemComponents(device_type, os_id, navigator_id)
	fmt.Println(system)
	app, _ := buildAppComponents(os_id, navigator_id)
	fmt.Println(app)
	ua_template := chooseUATemplate(device_type, navigator_id, app)
	fmt.Println(ua_template)
	t := template.Must(template.New("letter").Parse(ua_template.(string)))
	var tpl bytes.Buffer
	_ = t.Execute(&tpl, uatmpl{
		system,
		app,
	})
	user_agent := tpl.String()
	app_version := build_navigator_app_version(os_id, navigator_id, system["platform_version"], user_agent)
	return map[string]string{
		// ids
		"os_id":        os_id,
		"navigator_id": navigator_id,
		// system components
		"platform": system["platform"],
		"oscpu":    system["oscpu"],
		// app components
		"build_version": app["build_version"],
		"build_id":      app["build_id"],
		"app_version":   app_version,
		"app_name":      app["name"],
		"app_code_name": "Mozilla",
		"product":       "Gecko",
		"product_sub":   app["product_sub"],
		"vendor":        app["vendor"],
		"vendor_sub":    "",
		// compiled user agent
		"user_agent": user_agent,
	}
}

// Generate HTTP User-Agent header.
// Returns a string of HTTP header
func GenerateUserAgent(uaconfig ...UserAgentConfig) string {
	var cfg UserAgentConfig
	if len(uaconfig) == 0 {
		cfg = UserAgentConfig{}
	} else {
		cfg = uaconfig[0]
	}
	config := generateNavigator(&cfg)
	if config["user_agent"] == "" {
		panic("unable to generate user-agent")
	}
	return config["user_agent"]
}

/*
func GenerateNavigatorJS(*, os: None | str = None,navigator: None | str = None, platform: None | str = None,
    device_type: None | str = None,
) -> dict[str, None | str]:
    """Generate config for `windows.navigator` JavaScript object.

    :param os: limit list of oses for generation
    :type os: string or list/tuple or None
    :param navigator: limit list of browser engines for generation
    :type navigator: string or list/tuple or None
    :param device_type: limit possible oses by device type
    :type device_type: list/tuple or None, possible values:
        "desktop", "smartphone", "tablet", "all"
    :return: User-Agent config
    :rtype: dict with keys (TODO)
    :raises InvalidOption: if could not generate user-agent for
        any combination of allowed oses and navigators
    :raise InvalidOption: if any of passed options is invalid
    """
    config = generate_navigator(
        os=os, navigator=navigator, platform=platform, device_type=device_type
    )
    return {
        "appCodeName": config["app_code_name"],
        "appName": config["app_name"],
        "appVersion": config["app_version"],
        "platform": config["platform"],
        "userAgent": config["user_agent"],
        "oscpu": config["oscpu"],
        "product": config["product"],
        "productSub": config["product_sub"],
        "vendor": config["vendor"],
        "vendorSub": config["vendor_sub"],
        "buildID": config["build_id"],
    }*/
//This the cartesian product of input iterables. Its python equivalent is
//itertools.product() function(https://docs.python.org/3/library/itertools.html)*/
func product(a ...[]string) func() []string {

	if len(a) == 0 {
		panic("empty input")
	}
	p := make([]string, len(a))
	x := make([]int, len(p))
	return func() []string {
		p := p[:len(x)]
		for i, xi := range x {

			if len(a[i]) > xi {
				p[i] = a[i][xi]
			} else {
				p = nil
				break
			}

		}
		for i := len(x) - 1; i >= 0; i-- {
			x[i]++
			if x[i] < len(a) {
				break
			}
			x[i] = 0
			if i <= 0 {
				x = x[0:0]
				break
			}
		}
		return p
	}
}
