package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// CONFIG
const KEY_FILE_NAME = ".git-seal.key"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "setup":
		setupGit()
	case "keygen":
		generateKey()
	case "clean":
		processStream(true) // Encrypt
	case "smudge":
		processStream(false) // Decrypt
	default:
		printHelp()
	}
}

// 1. CRYPTO ENGINE (Deterministic AES)
func processStream(encrypt bool) {
	key := loadKey()

	// We use sha256 of the key to get exactly 32 bytes for AES-256
	keyHash := sha256.Sum256(key)
	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		die("Cipher error: " + err.Error())
	}

	// We use a fixed IV (Initialization Vector)
	// This is slightly less secure than random IV, but REQUIRED for Git.
	// If IV is random, every 'git diff' will show the whole file changed.
	iv := keyHash[:aes.BlockSize]

	stream := cipher.NewCFBEncrypter(block, iv)
	if !encrypt {
		stream = cipher.NewCFBDecrypter(block, iv)
	}

	// Stream Stdin -> Crypto -> Stdout
	writer := &cipher.StreamWriter{S: stream, W: os.Stdout}
	if _, err := io.Copy(writer, os.Stdin); err != nil {
		die("Stream error: " + err.Error())
	}
}

// 2. KEY MANAGEMENT
func loadKey() []byte {
	usr, _ := user.Current()
	path := filepath.Join(usr.HomeDir, KEY_FILE_NAME)

	key, err := os.ReadFile(path)
	if err != nil {
		// Fallback: If decrypting (smudge) and no key found,
		// we MUST output the raw input so the user sees the encrypted garbage
		// instead of crashing.
		if len(os.Args) > 1 && os.Args[1] == "smudge" {
			io.Copy(os.Stdout, os.Stdin)
			os.Exit(0)
		}
		die("Key not found at " + path + ". Run 'git-seal keygen' first.")
	}
	return key
}

func generateKey() {
	usr, _ := user.Current()
	path := filepath.Join(usr.HomeDir, KEY_FILE_NAME)

	// Generate 32 bytes of random data
	key := make([]byte, 32)
	// simple random generation
	f, _ := os.Open("/dev/urandom")
	f.Read(key)
	f.Close()

	if err := os.WriteFile(path, key, 0600); err != nil {
		die("Could not write key: " + err.Error())
	}
	fmt.Printf("✅ Master Key created at: %s\n", path)
	fmt.Println("⚠️  BACK UP THIS FILE IMMEDIATELY!")
}

// 3. GIT CONFIGURATION
func setupGit() {
	// Get path to this executable
	exe, _ := os.Executable()

	// Configure Git Global
	runCmd("git", "config", "--global", "filter.git-seal.clean", fmt.Sprintf("'%s' clean", exe))
	runCmd("git", "config", "--global", "filter.git-seal.smudge", fmt.Sprintf("'%s' smudge", exe))
	runCmd("git", "config", "--global", "filter.git-seal.required", "true")

	fmt.Println("✅ Git configured globally to use git-seal.")
	fmt.Println("👉 To use it in a repo, create .gitattributes and add:")
	fmt.Println("   .env filter=git-seal diff=git-seal")
}

func runCmd(name string, args ...string) {
	cmd := exec.Command(name, args...)
	if err := cmd.Run(); err != nil {
		die("Failed to run git config: " + err.Error())
	}
}

func die(msg string) {
	fmt.Fprintln(os.Stderr, "❌ "+msg)
	os.Exit(1)
}

func printHelp() {
	fmt.Println("git-seal: Transparent Encryption for Git")
	fmt.Println("Commands:")
	fmt.Println("  keygen   -> Generate your master key (~/.git-seal.key)")
	fmt.Println("  setup    -> Configure Git globally to use this tool")
	fmt.Println("  clean    -> (Internal) Encrypt data")
	fmt.Println("  smudge   -> (Internal) Decrypt data")
}
