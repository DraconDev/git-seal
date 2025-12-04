# Git-Seal: Transparent Encryption for Git

Git-seal provides seamless, transparent encryption for sensitive files in your Git repositories. Perfect for personal projects where you need to keep environment variables, API keys, and configuration files secure without complicated setup.

## 🚀 What is Git-Seal?

Git-seal is a lightweight Git filter that automatically encrypts specified files when committing to Git and decrypts them when checking out. It runs transparently in the background, so you work with plain text locally while Git stores encrypted versions.

### **Key Benefits:**

- **Zero Friction**: No manual unlock/lock commands
- **Git-Native**: Works seamlessly with all Git commands
- **Fast**: Stream-based encryption/decryption
- **Simple**: 3 commands to complete setup
- **Portable**: Single binary, cross-platform

## 🛡️ Security Model

- **Encryption**: AES-256-CFB encryption
- **Key Management**: Single master key stored locally (`~/.git-seal.key`)
- **Deterministic**: Fixed IV ensures proper Git diff behavior
- **Local Only**: Your key never leaves your machine

## 📦 Installation

### Prerequisites

- Go 1.16+ (for building from source)
- Git

### Build & Install

```bash
# Clone or download the source code
# Save main.go from the repository

# Compile
go build -o git-seal main.go

# Install system-wide (optional)
sudo mv git-seal /usr/local/bin/
```

## ⚙️ Setup (One-Time)

### 1. Generate Master Key

```bash
git-seal keygen
```

This creates `~/.git-seal.key` - **BACK IT UP IMMEDIATELY!**

### 2. Configure Git

```bash
git-seal setup
```

## 📋 Usage

### Basic Workflow

#### 1. Create/Edit Files

```bash
echo "API_KEY=secret123" > .env
```

#### 2. Configure Encryption

Create `.gitattributes` in your repository:

```
.env filter=git-seal diff=git-seal
```

#### 3. Commit

```bash
git add .
git commit -m "Add environment configuration"
```

**Result:**

- **Local**: `.env` shows `API_KEY=secret123` (readable)
- **Git**: File is stored as encrypted binary (unreadable)

### The "Clone & Delete" Cycle

The ultimate convenience - no manual unlocking required:

```bash
# Delete local copy
rm -rf my-project

# Clone fresh
git clone <your-repo-url>
cd my-project

# File is automatically decrypted!
cat .env  # Shows: API_KEY=secret123
```

## 🔧 Advanced Usage

### Encrypt Multiple File Types

In `.gitattributes`:

```
.env filter=git-seal diff=git-seal
config/secrets.json filter=git-seal diff=git-seal
*.key filter=git-seal diff=git-seal
```

### Manual Encryption/Decryption

```bash
# Encrypt a file
cat sensitive.txt | git-seal clean > encrypted.bin

# Decrypt a file
cat encrypted.bin | git-seal smudge > sensitive.txt
```

### Check What's Encrypted

```bash
# See original content locally
cat .env

# See encrypted content in Git
git show HEAD:.env
```

## 📊 Git-Seal vs git-crypt

## ⚠️ Security Consideration: Duplicate API Keys

### **The Problem Demonstrated**

When the same API key is used across multiple projects, git-seal produces **identical encrypted content**:

```bash
# Project 1: API_KEY=secret123 → [ENCRYPTED BYTES: 1f87 30e6...]
# Project 2: API_KEY=secret123 → [ENCRYPTED BYTES: 1f87 30e6...]
```

**Result**: Identical bytes reveal you're reusing secrets across projects!

### **Why This Happens**

Git-seal uses **deterministic encryption** (fixed IV derived from key) for Git compatibility:

- **Same input** → **Same encrypted output**
- Enables proper Git diffs and merges
- But **sacrifices some security** for convenience

### **Attack Scenario**

```bash
# Attacker clones your multiple repos
git clone your-project-1 your-project-2

# Compares encrypted .env files
diff project1/.env project2/.env  # IDENTICAL!
# → Attacker knows you reuse the same API key
```

### **🛡️ Solutions to Avoid This**

#### 1. **Use Unique Environment Variable Names**

```bash
# ❌ Bad: Same key everywhere
API_KEY=secret123

# ✅ Good: Unique names
PROJECT1_API_KEY=secret123
PROJECT2_API_KEY=secret123
```

#### 2. **Environment-Specific Files**

```bash
# Project 1
.env.prod filter=git-seal diff=git-seal

# Project 2
.env.staging filter=git-seal diff=git-seal
```

#### 3. **Namespace Your Secrets**

```bash
# Instead of generic names
# API_KEY=secret123

# Use descriptive names
WEB_API_KEY=secret123
MOBILE_API_KEY=secret123
ADMIN_API_KEY=secret123
```

### **When This Matters**

| Use Case                | Risk Level | Recommendation                      |
| ----------------------- | ---------- | ----------------------------------- |
| Personal projects       | 🟢 Low     | Fine as-is, but use unique names    |
| Team projects           | 🟡 Medium  | Always use unique environment names |
| Enterprise environments | 🔴 High    | Consider HashiCorp Vault instead    |

### **Key Takeaway**

Git-seal prioritizes **workflow convenience over perfect security**. For maximum security, use unique environment variable names or consider enterprise key management solutions.
| Feature | Git-Seal | git-crypt |
| ----------------- | ------------------------- | ----------------------------- |
| Setup Complexity | ⭐ Simple (3 commands) | ⭐⭐⭐ Complex (GPG setup) |
| Performance | ⭐⭐⭐ Fast (streaming) | ⭐⭐ Moderate (file-based) |
| Workflow Friction | ⭐⭐⭐ Zero unlock needed | ⭐⭐ Manual unlock/lock |
| Dependencies | ⭐ None (single binary) | ⭐ GPG required |
| Cross-Platform | ⭐⭐⭐ Perfect | ⭐⭐ GPG compatibility issues |
| Git Integration | ⭐⭐⭐ Seamless | ⭐⭐ Good |

## 🔐 Security Considerations

### Strengths

- **AES-256** encryption is cryptographically strong
- **Local key storage** - no server-side key exposure
- **Deterministic encryption** enables proper Git diffs

### Trade-offs (By Design)

- **Fixed IV**: Less secure than random IV, but required for Git compatibility
- **Local key**: If someone accesses your machine, they can decrypt files
- **No passphrase**: Prioritizes convenience over additional security layer

### Best Practices

1. **Backup your key** to a secure location (password manager, encrypted storage)
2. **Use on personal projects** only - not for team environments
3. **Keep the binary secure** - anyone with the binary and key can decrypt
4. **Regular key rotation** if security requirements demand it

## 🛠️ Troubleshooting

### Key Not Found

```
Error: Key not found at ~/.git-seal.key. Run 'git-seal keygen' first.
```

**Solution**: Run `git-seal keygen`

### Permission Denied

```
Error: Failed to run git config
```

**Solution**: Ensure you have Git configured and proper permissions

### Files Not Encrypting

1. Check `.gitattributes` is in repository root
2. Verify Git filter is configured: `git config --get-regexp filter`
3. Ensure file pattern matches exactly

### Git Diff Shows Garbled Text

This is expected! Git diff shows encrypted content. Use `git show` to see the actual diff:

```bash
git show HEAD:.env | git-seal smudge | diff - .env
```

## 🏗️ Technical Architecture

### Git Filter Flow

```
1. git add .env
   → Git calls: git-seal clean < .env > encrypted_version
   → Stores encrypted_version in Git index

2. git checkout
   → Git calls: git-seal smudge < encrypted_version > .env
   → Creates readable .env file locally
```

### Encryption Process

```
Input Text → AES-256-CFB → Encrypted Binary → Git Storage
     ↑                                                    ↓
Local File ← AES-256-CFB ← Decrypted Binary ← Git Storage
```

## 📝 Example Project Structure

```
my-app/
├── .env                    # API_KEY=secret123 (readable locally)
├── .gitattributes         # .env filter=git-seal diff=git-seal
├── .git-seal.key          # Your master key (don't commit!)
├── src/
└── README.md
```

## ⚡ Quick Start Summary

```bash
# 1. Build & Install
go build -o git-seal main.go

# 2. Setup (one-time)
git-seal keygen    # Backup ~/.git-seal.key!
git-seal setup

# 3. Use
echo "API_KEY=secret" > .env
echo ".env filter=git-seal diff=git-seal" > .gitattributes
git add . && git commit -m "Secure config"
```

**That's it!** Your files are now encrypted in Git, decrypted locally, with zero ongoing effort.

---

_Built with ❤️ for developers who value both security and simplicity._

# Git-Seal: Transparent Encryption for Git

Git-seal provides seamless, transparent encryption for sensitive files in your Git repositories. Perfect for personal projects where you need to keep environment variables, API keys, and configuration files secure without complicated setup.

## 🚀 What is Git-Seal?

Git-seal is a lightweight Git filter that automatically encrypts specified files when committing to Git and decrypts them when checking out. It runs transparently in the background, so you work with plain text locally while Git stores encrypted versions.

### **Key Benefits:**

- **Zero Friction**: No manual unlock/lock commands
- **Git-Native**: Works seamlessly with all Git commands
- **Fast**: Stream-based encryption/decryption
- **Simple**: 3 commands to complete setup
- **Portable**: Single binary, cross-platform

## 🛡️ Security Model

- **Encryption**: AES-256-CFB encryption
- **Key Management**: Single master key stored locally (`~/.git-seal.key`)
- **Deterministic**: Fixed IV ensures proper Git diff behavior
- **Local Only**: Your key never leaves your machine

## 📦 Installation

### Prerequisites

- Go 1.16+ (for building from source)
- Git

### Build & Install

```bash
# Clone or download the source code
# Save main.go from the repository

# Compile
go build -o git-seal main.go

# Install system-wide (optional)
sudo mv git-seal /usr/local/bin/
```

## ⚙️ Setup (One-Time)

### 1. Generate Master Key

```bash
git-seal keygen
```

This creates `~/.git-seal.key` - **BACK IT UP IMMEDIATELY!**

### 2. Configure Git

```bash
git-seal setup
```

## 📋 Usage

### Basic Workflow

#### 1. Create/Edit Files

```bash
echo "API_KEY=secret123" > .env
```

#### 2. Configure Encryption

Create `.gitattributes` in your repository:

```
.env filter=git-seal diff=git-seal
```

#### 3. Commit

```bash
git add .
git commit -m "Add environment configuration"
```

**Result:**

- **Local**: `.env` shows `API_KEY=secret123` (readable)
- **Git**: File is stored as encrypted binary (unreadable)

### The "Clone & Delete" Cycle

The ultimate convenience - no manual unlocking required:

```bash
# Delete local copy
rm -rf my-project

# Clone fresh
git clone <your-repo-url>
cd my-project

# File is automatically decrypted!
cat .env  # Shows: API_KEY=secret123
```

## 🔧 Advanced Usage

### Encrypt Multiple File Types

In `.gitattributes`:

```
.env filter=git-seal diff=git-seal
config/secrets.json filter=git-seal diff=git-seal
*.key filter=git-seal diff=git-seal
```

### Manual Encryption/Decryption

```bash
# Encrypt a file
cat sensitive.txt | git-seal clean > encrypted.bin

# Decrypt a file
cat encrypted.bin | git-seal smudge > sensitive.txt
```

### Check What's Encrypted

```bash
# See original content locally
cat .env

# See encrypted content in Git
git show HEAD:.env
```

## 📊 Git-Seal vs git-crypt

| Feature           | Git-Seal                  | git-crypt                     |
| ----------------- | ------------------------- | ----------------------------- |
| Setup Complexity  | ⭐ Simple (3 commands)    | ⭐⭐⭐ Complex (GPG setup)    |
| Performance       | ⭐⭐⭐ Fast (streaming)   | ⭐⭐ Moderate (file-based)    |
| Workflow Friction | ⭐⭐⭐ Zero unlock needed | ⭐⭐ Manual unlock/lock       |
| Dependencies      | ⭐ None (single binary)   | ⭐ GPG required               |
| Cross-Platform    | ⭐⭐⭐ Perfect            | ⭐⭐ GPG compatibility issues |
| Git Integration   | ⭐⭐⭐ Seamless           | ⭐⭐ Good                     |

## 🔐 Security Considerations

### Strengths

- **AES-256** encryption is cryptographically strong
- **Local key storage** - no server-side key exposure
- **Deterministic encryption** enables proper Git diffs

### Trade-offs (By Design)

- **Fixed IV**: Less secure than random IV, but required for Git compatibility
- **Local key**: If someone accesses your machine, they can decrypt files
- **No passphrase**: Prioritizes convenience over additional security layer

### Best Practices

1. **Backup your key** to a secure location (password manager, encrypted storage)
2. **Use on personal projects** only - not for team environments
3. **Keep the binary secure** - anyone with the binary and key can decrypt
4. **Regular key rotation** if security requirements demand it

## 🛠️ Troubleshooting

### Key Not Found

```
Error: Key not found at ~/.git-seal.key. Run 'git-seal keygen' first.
```

**Solution**: Run `git-seal keygen`

### Permission Denied

```
Error: Failed to run git config
```

**Solution**: Ensure you have Git configured and proper permissions

### Files Not Encrypting

1. Check `.gitattributes` is in repository root
2. Verify Git filter is configured: `git config --get-regexp filter`
3. Ensure file pattern matches exactly

### Git Diff Shows Garbled Text

This is expected! Git diff shows encrypted content. Use `git show` to see the actual diff:

```bash
git show HEAD:.env | git-seal smudge | diff - .env
```

## 🏗️ Technical Architecture

### Git Filter Flow

```
1. git add .env
   → Git calls: git-seal clean < .env > encrypted_version
   → Stores encrypted_version in Git index

2. git checkout
   → Git calls: git-seal smudge < encrypted_version > .env
   → Creates readable .env file locally
```

### Encryption Process

```
Input Text → AES-256-CFB → Encrypted Binary → Git Storage
     ↑                                                    ↓
Local File ← AES-256-CFB ← Decrypted Binary ← Git Storage
```

## 📝 Example Project Structure

```
my-app/
├── .env                    # API_KEY=secret123 (readable locally)
├── .gitattributes         # .env filter=git-seal diff=git-seal
├── .git-seal.key          # Your master key (don't commit!)
├── src/
└── README.md
```

## ⚡ Quick Start Summary

```bash
# 1. Build & Install
go build -o git-seal main.go

# 2. Setup (one-time)
git-seal keygen    # Backup ~/.git-seal.key!
git-seal setup

# 3. Use
echo "API_KEY=secret" > .env
echo ".env filter=git-seal diff=git-seal" > .gitattributes
git add . && git commit -m "Secure config"
```

**That's it!** Your files are now encrypted in Git, decrypted locally, with zero ongoing effort.

---

_Built with ❤️ for developers who value both security and simplicity._
