package distro

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dv-net/dv-updater/pkg/logger"
	"io"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
)

const (
	unknownName    = "Unknown"
	unknownVersion = "unknown"
)

var FileSystemRoot = string(os.PathSeparator)

// equalsSplitter is a regex to split apart key value pairs delimited with an equals sign
var equalsSplitter = regexp.MustCompile(`^\s*(\S+)\s*=\s*([\S ]+)\s*`)

// releaseSplitter is a regex to split apart the contents of /etc/*-release files in the Red Hat Format
var releaseSplitter = regexp.MustCompile(`^(.+) (release|version)? (\S+)\s*(\S+)?`)

type IDistro interface {
	DiscoverDistro() (LinuxDistro, error)
}

type Service struct {
	logger logger.Logger
}

func New(l logger.Logger) *Service {
	return &Service{
		logger: l,
	}
}

func (s *Service) DiscoverDistro() (LinuxDistro, error) {
	if runtime.GOOS != "linux" {
		return LinuxDistro{}, fmt.Errorf("OS %s is not supported", runtime.GOOS)
	}

	lsbProperties, _ := s.readReleaseFile("/etc/lsb-release")
	osReleaseProperties, _ := s.readReleaseFile("/etc/os-release")

	return discoverDistroFromProperties(lsbProperties, osReleaseProperties), nil
}

func discoverDistroFromProperties(lsbProperties ReleaseDetails, osReleaseProperties ReleaseDetails) LinuxDistro {
	var detectedDistro LinuxDistro
	wasDetected := false

	for _, distroTest := range List {
		wasDetected, detectedDistro = distroTest(lsbProperties, osReleaseProperties)

		if wasDetected {
			break
		}
	}

	if !wasDetected {
		detectedDistro = BestGuess(lsbProperties, osReleaseProperties)
	}

	return detectedDistro
}

func (s *Service) readReleaseFile(filePath string) (ReleaseDetails, error) {
	reader, pathRead, openErr := readBinaryFileFunc([]string{filePath})
	if openErr != nil {
		if pathRead != "" {
			s.logger.Info("unable to read release file at the path: %s", pathRead)
		}

		return ReleaseDetails{}, openErr
	}
	defer func() { _ = reader.Close() }()

	properties, parseErr := parseOSRelease(reader)
	return properties, parseErr
}

func readFileFunc(filePaths ...string) (bool, string) {
	reader, _, err := readBinaryFileFunc(filePaths)
	if err != nil {
		return false, ""
	}

	defer func() { _ = reader.Close() }()

	contents, err := io.ReadAll(reader)
	if err != nil {
		return false, ""
	}

	return true, string(contents)
}

func readBinaryFileFunc(filePaths []string) (io.ReadCloser, string, error) {
	for _, filePath := range filePaths {
		if FileSystemRoot != string(os.PathSeparator) {
			filePath = path.Clean(FileSystemRoot + string(os.PathSeparator) + filePath)
		}

		fileInfo, statErr := os.Stat(filePath)
		if statErr != nil || fileInfo.IsDir() {
			continue
		}

		reader, readErr := os.Open(filePath)
		if readErr != nil {
			return nil, filePath, readErr
		}

		return reader, filePath, nil
	}

	errMsg := fmt.Sprintf("unable to create a reader for any of the specified paths: %v", filePaths)
	return nil, "", errors.New(errMsg)
}

func parseOSRelease(reader io.Reader) (ReleaseDetails, error) {
	properties := ReleaseDetails{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		key, val, splitErr := splitEqualsKeyVal(line)
		if splitErr != nil {
			continue
		}

		properties[key] = val
	}

	return properties, scanner.Err()
}

func splitEqualsKeyVal(line string) (string, string, error) {
	if line == "" {
		return "", "", errors.New("can't split a blank line")
	}

	if line[0] == '#' {
		return "", "", fmt.Errorf("ignoring commented line: %s", line)
	}

	match := equalsSplitter.FindStringSubmatch(line)
	if len(match) == 0 {
		return "", "", fmt.Errorf("no splittable character for line: %s", line)
	}
	if len(match) != 3 {
		return "", "", fmt.Errorf("unexpected number of matches (%d) for line: %s", len(match), line)
	}

	withoutTrailingWhitespace := strings.TrimSpace(match[2])
	withoutEnclosingQuotes := strings.Trim(withoutTrailingWhitespace, "\"")

	return match[1], withoutEnclosingQuotes, nil
}

func BestGuess(lsbProperties ReleaseDetails, osReleaseProperties ReleaseDetails) LinuxDistro {
	var id string
	switch {
	case osReleaseProperties["ID"] != "":
		id = osReleaseProperties["ID"]
	case lsbProperties["DISTRIB_ID"] != "":
		id = strings.ToLower(lsbProperties["DISTRIB_ID"])
	default:
		id = unknownVersion
	}

	var name string
	switch {
	case osReleaseProperties["NAME"] != "":
		name = osReleaseProperties["NAME"]
	case osReleaseProperties["PRETTY_NAME"] != "":
		segments := strings.SplitN(osReleaseProperties["PRETTY_NAME"], " ", 2)
		name = segments[0]
	case lsbProperties["DISTRIB_ID"] != "":
		name = lsbProperties["DISTRIB_ID"]
	case osReleaseProperties["ID"] != "":
		name = osReleaseProperties["ID"]
	default:
		name = unknownName
	}

	var version string
	switch {
	case osReleaseProperties["VERSION_ID"] != "":
		version = osReleaseProperties["VERSION_ID"]
	case lsbProperties["DISTRIB_RELEASE"] != "":
		version = lsbProperties["DISTRIB_RELEASE"]
	case osReleaseProperties["VERSION"] != "":
		segments := strings.SplitN(osReleaseProperties["VERSION"], " ", 2)
		version = segments[0]
	default:
		version = unknownVersion
	}

	return LinuxDistro{
		Name:       name,
		ID:         id,
		Version:    version,
		LsbRelease: lsbProperties,
		OsRelease:  osReleaseProperties,
	}
}

var List = []func(ReleaseDetails, ReleaseDetails) (bool, LinuxDistro){
	IsCentOS,
	IsRHEL,
	IsUbuntu,
	IsDebian,
	IsAmazonLinux,
	IsFedora,
	IsOpenSuSE,
	IsSLES,
	IsOracleLinux,
	IsPhoton,
	IsAlpine,
	IsArchLinux,
	IsGentoo,
	IsKali,
	IsScientificLinux,
	IsSlackware,
	IsMageia,
	IsClearLinux,
	IsMint,
	IsMXLinux,
	IsNovellOES,
	IsPuppy,
	IsRancherOS,
	IsNixOS,
	IsAlt,
	IsCrux,
	IsSourceMage,
	IsAndroid,
	IsYellowDog,
}

func parseRedhatReleaseContents(contents string, expectedDistro string) (bool, string) {
	matches := releaseSplitter.FindStringSubmatch(contents)

	if !strings.HasPrefix(matches[0], expectedDistro) {
		return false, ""
	}

	var version string

	if len(matches) > 3 {
		version = strings.TrimSpace(matches[3])
	} else {
		version = unknownVersion
	}

	return true, version
}
