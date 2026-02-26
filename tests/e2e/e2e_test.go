//go:build e2e

package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestE2E_SopsvLifecycle(t *testing.T) {
	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "sopsv")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, "../../main.go")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	err := buildCmd.Run()
	require.NoError(t, err, "failed to build sopsv binary")

	// Set up temporary directories to avoid messing with real user config
	tempHome := t.TempDir()
	env := os.Environ()
	env = append(env, "XDG_CONFIG_HOME="+filepath.Join(tempHome, "config"))
	env = append(env, "XDG_DATA_HOME="+filepath.Join(tempHome, "data"))
	
	// Add local bin to path for installed tools
	localBin := filepath.Join(tempHome, ".local", "bin")
	err = os.MkdirAll(localBin, 0755)
	require.NoError(t, err)

	pathIdx := -1
	for i, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			pathIdx = i
			break
		}
	}
	if pathIdx != -1 {
		env[pathIdx] = "PATH=" + localBin + string(os.PathListSeparator) + strings.TrimPrefix(env[pathIdx], "PATH=")
	} else {
		env = append(env, "PATH="+localBin)
	}
	
	// Helper function to run the binary
	runCmd := func(args ...string) (string, error) {
		cmd := exec.Command(binaryPath, args...)
		cmd.Env = env
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		return out.String(), err
	}

	// 0. Install tools
	out, err := runCmd("tool", "install", "-d", localBin)
	require.NoError(t, err, "failed to install tools: %s", out)

	// 1. Vault New
	vaultName := "e2evault"
	out, err = runCmd("vault", "new", vaultName)
	require.NoError(t, err, "failed to create new vault: %s", out)
	require.Contains(t, out, "created successfully")

	// 2. Set Secret
	out, err = runCmd("set", "-k", "MY_SECRET", "-v", "supersecret123", "--vault", vaultName)
	require.NoError(t, err, "failed to set secret: %s", out)
	require.Contains(t, out, "set successfully")

	// 3. Get Secret
	out, err = runCmd("get", "-k", "MY_SECRET", "--vault", vaultName)
	require.NoError(t, err, "failed to get secret: %s", out)
	require.Equal(t, "supersecret123", strings.TrimSpace(out))

	// 4. Remove Secret
	out, err = runCmd("rm", "-k", "MY_SECRET", "--vault", vaultName)
	require.NoError(t, err, "failed to rm secret: %s", out)
	require.Contains(t, out, "removed from vault")

	// 5. Delete Vault
	out, err = runCmd("vault", "rm", vaultName)
	require.NoError(t, err, "failed to delete vault: %s", out)
	require.Contains(t, out, "deleted successfully")
}
