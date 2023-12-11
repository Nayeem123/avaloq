package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/pelletier/go-toml"
)

type Path struct {
	CurrentPath string `json:"currentPath"`
}

func WhoAmI(w http.ResponseWriter, r *http.Request) {
	config, err := toml.LoadFile("./configs/config.toml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	command := config.Get("whoami.execute").(string)

	//executableCommand := command + " " + path.CurrentPath
	fmt.Println("Command that to run: ", command)
	var out []byte
	switch runtime.GOOS {
	case "windows":
		// For Windows, you might need to adjust the command and use the appropriate shell
		cmd := exec.Command("cmd", "/C", command)
		out, err = cmd.Output()
	case "linux", "darwin":
		// For Linux and macOS, you might not need to adjust the command, using sh should be fine
		cmd := exec.Command("sh", "-c", command)
		out, err = cmd.Output()
	default:
		http.Error(w, "Unsupported Operating System", http.StatusInternalServerError)
		return
	}

	if err != nil {
		errorMessage := fmt.Sprintf("Error executing command: %s\n%s", command, err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// Replace newlines with <br> for better readability in HTML responses
	//output := strings.ReplaceAll(string(out), "\n", "<br>")
	//fmt.Fprint(w, output)

	fmt.Fprint(w, string(out))
}

func ExecuteCommandHandler(w http.ResponseWriter, r *http.Request) {
	res := Authentication(w, r)
	if !bool(res) {
		http.Error(w, "API Authentication Failed , Pass valid Token", http.StatusUnauthorized)
		return
	}
	isAdmin := ExecuteUserRolePolicy
	//fmt.Println("policy = ", bool(isAdmin()))
	if bool(isAdmin()) != true {
		http.Error(w, "Access Denied", http.StatusUnauthorized)
		return
	}
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	// Parse JSON data
	var path Path
	err = json.Unmarshal(body, &path)
	if err != nil {
		http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
		return
	}

	fmt.Println("body data: ", path.CurrentPath)

	config, err := toml.LoadFile("./configs/config.toml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	command := config.Get("command.execute").(string)

	executableCommand := command + " " + path.CurrentPath
	fmt.Println("Command that to run: ", executableCommand)
	var out []byte
	switch runtime.GOOS {
	case "windows":
		// For Windows, you might need to adjust the command and use the appropriate shell
		cmd := exec.Command("cmd", "/C", executableCommand)
		out, err = cmd.Output()
	case "linux", "darwin":
		// For Linux and macOS, you might not need to adjust the command, using sh should be fine
		cmd := exec.Command("sh", "-c", executableCommand)
		out, err = cmd.Output()
	default:
		http.Error(w, "Unsupported Operating System", http.StatusInternalServerError)
		return
	}

	if err != nil {
		errorMessage := fmt.Sprintf("Error executing command: %s\n%s", executableCommand, err.Error())
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}

	// Replace newlines with <br> for better readability in HTML responses
	//output := strings.ReplaceAll(string(out), "\n", "<br>")
	//fmt.Fprint(w, output)

	fmt.Fprint(w, string(out))

}
