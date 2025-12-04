# Git-Seal: Transparent Encryption for Git

**Simple, secure encryption for your environment files and sensitive configuration data.**

## 🚀 Quick Start

### 1. Install the Binary

```bash
# Make it executable and move to system path
chmod +x git-seal
sudo mv git-seal /usr/local/bin/
```

### 2. Generate Your Master Key (One-Time)

```bash
git-seal keygen
```

**⚠️ IMPORTANT**: Backup your `~/.git-seal.key` file immediately!

### 3. Configure Git Integration (One-Time)

```bash
git-seal setup
```

### 4. Protect Your Files

Create `.gitattributes` in your repository:

```
.env filter=git-seal diff=git-seal
```

Add to `.gitignore` for safety:

```
.env
```

### 5. Use It

```bash
# Local development - plain text .env
echo "API_KEY=secret123" > .env

# When ready to commit - encrypted automatically
git add -f .env
git commit -m "Add encrypted config"

# Clone anywhere - automatically decrypted!
git clone <your-repo>
# .env is instantly available as plain text
```

## 🔧 Advanced Usage

### Multiple File Types

```
# .gitattributes
.env filter=git-seal diff=git-seal
*.key filter=git-seal diff=git-seal
config/secrets.json filter=git-seal diff=git-seal
```

### Manual Encryption/Decryption

```bash
# Encrypt
cat sensitive.txt | ./git-seal clean > encrypted.bin

# Decrypt
cat encrypted.bin | ./git-seal smudge > sensitive.txt
```

## 🛡️ Security Features

- **AES-256-CFB** encryption
- **Transparent operation** - no manual unlock commands
- **Git-native integration** - works with all Git operations
- **Local key storage** - your key never leaves your machine
- **Cross-platform** - single binary for all operating systems

## ⚠️ Important Notes

1. **Backup your key**: `~/.git-seal.key` contains your master encryption key
2. **Safe by default**: Add `.env` to `.gitignore` to prevent accidental commits
3. **Flexible override**: Use `git add -f .env` when you want to commit encrypted version
4. **Key rotation**: Run `git-seal keygen` to generate a new master key

## 🔧 System Requirements

- Git (for integration)
- Linux/macOS/Windows (binary included)
- No additional dependencies

## 📚 Documentation

For detailed explanation of how Git's filter system works automatically, see `GIT_FILTER_AUTOMATION.md`.

For advanced workflow patterns and best practices, see `GIT_SEAL_WORKFLOW.md`.

## 🤝 License

MIT License - Free to use for personal and commercial projects.

---

**Built with ❤️ for developers who value both security and simplicity.**
