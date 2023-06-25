package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

/*********************************************
 * Workspace Root
 *********************************************/

type workspaceRoot struct {
	path string
};

func systemWorkspaceRoot() workspaceRoot {
	const WORKSPACE_DIR = "ws"

	home, err := os.UserHomeDir()

	if err != nil {
		log.Fatalf("The user home directory is not defined: %v", err)
	}

	workspaceRootPath := filepath.Join(home, WORKSPACE_DIR);

	_, err = os.Stat(workspaceRootPath)

	if os.IsNotExist(err) {
		if err := os.Mkdir(workspaceRootPath, 0755); err != nil {
			log.Fatalf("Failed to create directory: %s", workspaceRootPath);
		}
	}

	return workspaceRoot { workspaceRootPath }
}

type workspaceDirs []*workspaceDir

func (wsRoot *workspaceRoot) readWorkspaces() workspaceDirs {
	dirs, err := os.ReadDir(wsRoot.path)

	if err != nil {
		log.Fatalf("Failed to read workspaces in %s: %v", wsRoot.path, err);
	}

	workspaces := make([]*workspaceDir, 0)

	for _, dir := range dirs {
		if !dir.Type().IsDir() {
			continue
		}

		wsDir, err := newWorkspaceFromName(wsRoot, dir.Name())

		if err != nil {
			log.Printf("Error reading directory in workspace %s: %s: %v", wsRoot.path, dir.Name(), err);
			continue
		}

		workspaces = append(workspaces, wsDir);
	}

	return workspaces
}

func (dirs workspaceDirs) findDirByName(name string) *workspaceDir {
	birthTime, nonce, err := parseWorkspaceName(name)

	if err != nil {
		log.Fatalf("Cannot find directory %s, error parsing name", name);
	}

	var dir *workspaceDir = nil

	for _, dir = range dirs {
		if dir.birthTime.Equal(*birthTime) && dir.nonce == *nonce {
			break
		}
	}

	return dir
}

/*********************************************
 * Workspace Directory Object
 *********************************************/

type workspaceDir struct {
	birthTime time.Time
	nonce string
	accessTime time.Time
	parent *workspaceRoot
}

var MisformattedWorkspaceError = errors.New("Misformatted workspace name error")

func formatWorkspaceName(birthTime time.Time, nonce string) string {
	dateTime := birthTime.UTC().Format(time.DateTime)
	dateTime = strings.ReplaceAll(dateTime, " ", "_")
	return dateTime + "." + nonce
}

func parseWorkspaceName(name string) (birthTime *time.Time, nonce *string, err error) {
	parts := strings.Split(name, ".")

	if len(parts) != 2 {
		return nil, nil, MisformattedWorkspaceError
	}

	dateTime := strings.ReplaceAll(parts[0], "_", " ")

	timePart, err := time.ParseInLocation(time.DateTime, dateTime, time.UTC)
	noncePart := parts[1]

	if err != nil {
		return nil, nil, MisformattedWorkspaceError
	}

	return &timePart, &noncePart, nil
}

func newWorkspaceFromName(parent *workspaceRoot, name string) (dir *workspaceDir, err error) {
	dir = nil

	birthTime, nonce, err := parseWorkspaceName(name)
	if err != nil {
		return
	}

	st, err := os.Stat(path.Join(parent.path, name))
	if err != nil {
		return
	}

	dir = &workspaceDir {
		*birthTime,
		*nonce,
		st.ModTime(),
		parent,
	}

	return
}

func createWorkspaceDir(ws *workspaceRoot) workspaceDir {
	birthTime := time.Now()
	namePattern := formatWorkspaceName(birthTime, "*")

	dir, err := os.MkdirTemp(ws.path, namePattern)
	if err != nil {
		log.Fatalf("Failed to create new workspace directory: %v", err);
	}

	dir = path.Base(dir)

	_, nonce, err := parseWorkspaceName(dir)

	if err != nil {
		log.Fatalf("Failed to parse the workspace name: %v\n", err);
	}

	return workspaceDir {
		birthTime,
		*nonce,
		birthTime,
		ws,
	};
}

func (dir *workspaceDir) name() string {
	return formatWorkspaceName(dir.birthTime, dir.nonce)
}

func (dir *workspaceDir) path() string {
	return path.Join(dir.parent.path, dir.name())
}

func (dir *workspaceDir) timeSinceLastAccess() time.Duration {
	return time.Since(dir.accessTime)
}

/*********************************************
 * Commands Helpers
 *********************************************/

const (
	secondsInSec = 1
	minsInSec = 60 * secondsInSec
	hourInSec = 60 * minsInSec
	dayInSec = 24 * hourInSec
)

func formatTimestamp(duration time.Duration) string {
	secs := int64(duration.Seconds())
	days := secs / dayInSec
	secs %= dayInSec
	hours := secs / hourInSec
	secs %= hourInSec
	mins := secs / minsInSec
	secs %= minsInSec
	
	return fmt.Sprintf("%d days %d hours %d min %d secs", days, hours, mins, secs)
}

func readRecencySortedWorkspaces() (*workspaceRoot, workspaceDirs) {
	wsRoot := systemWorkspaceRoot()
	workspaces := wsRoot.readWorkspaces()

	sort.Slice(workspaces, func(i, j int) bool {
		return workspaces[i].accessTime.Compare(workspaces[j].accessTime) < 0
	})

	return &wsRoot, workspaces
}

func currentWorkspace(wsRoot *workspaceRoot, workspaces workspaceDirs) int {

	cwd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to obtain the working directory: %s\n", err);
	}

	if !strings.HasPrefix(cwd, wsRoot.path) {
		log.Fatalf("The current working directory is not in a workspace\n")
	}

	dirName := path.Base(cwd)

	cur_idx := 0
	for ; cur_idx < len(workspaces); cur_idx++ {
		dir := workspaces[cur_idx]
		if dir.name() == dirName {
			break
		}
	}

	if cur_idx == len(workspaces) {
		log.Fatalf("Failed to find %s\n", dirName)
	}

	return cur_idx
}

/*********************************************
 * Commands
 *********************************************/

func createNewWorkspace() {
	ws := systemWorkspaceRoot()
	workspace := createWorkspaceDir(&ws)
	path := workspace.path()

	if err := os.Chdir(path); err != nil {
		log.Fatalf("Failed to change into directory: %v\n", err)
	}

	fmt.Println(path)
}

func listWorkspaces() {
	_, workspaces := readRecencySortedWorkspaces()

	for i, ws := range workspaces {
		delta := ws.timeSinceLastAccess()

		fmt.Println(i + 1, ws.name(), formatTimestamp(delta))
	}
}

func printRecentWorkspace() {
	_, workspaces := readRecencySortedWorkspaces()

	if len(workspaces) == 0 {
		log.Fatalln("There are currently no workspaces. Create a workspace with `ws new`.")
	}

	fmt.Println(workspaces[len(workspaces) - 1].path())
}

type printWorkspaceSubcmd int

const (
	PrevWs printWorkspaceSubcmd = iota
	NextWs
	CurrWs
)

func printWorkspace(subcmd printWorkspaceSubcmd) {
	wsRoot, workspaces := readRecencySortedWorkspaces()

	idx := currentWorkspace(wsRoot, workspaces)

	switch subcmd {
	case PrevWs:
		if idx <= 0 {
			log.Fatalln("The first workspace has been reached")
		}

		idx--
	case NextWs:
		if idx + 1 >= len(workspaces) {
			log.Fatalln("The current workspace is the newest workspace")
		}

		idx++
	}

	fmt.Println(workspaces[idx].path())
}

func printWorkspaceRoot() {
	wsRoot := systemWorkspaceRoot()
	fmt.Println(wsRoot.path)
}

func fatalPrintUsage() {
fmt.Printf(`
	ws internal

	ws_internal create_new_workspace
	ws_internal list_workspaces
	ws_internal print_recent_workspace
	ws_internal print_current_workspace
	ws_internal print_next_workspace
	ws_internal print_prev_workspace
	ws_internal print_workspace_root
	ws_internal activate
`)
	os.Exit(1)
}

func main() {

	if len(os.Args) < 1 {
		os.Exit(127)
	}

	prog := os.Args[0]

	log.SetPrefix(prog + ": ")
	log.SetFlags(0)

	if len(os.Args) < 2 {
		fatalPrintUsage()
	}

	cmd := os.Args[1]

	switch cmd {
	case "create_new_workspace":
		createNewWorkspace()
	case "list_workspaces":
		listWorkspaces()
	case "print_recent_workspace":
		printRecentWorkspace()
	case "print_current_workspace":
		printWorkspace(CurrWs)
	case "print_next_workspace":
		printWorkspace(NextWs)
	case "print_prev_workspace":
		printWorkspace(PrevWs)
	case "print_workspace_root":
		printWorkspaceRoot()
	case "activate":
		fmt.Print(shellWrapper)
	default:
		fatalPrintUsage()
	}
	
}
