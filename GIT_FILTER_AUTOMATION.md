# Git Filter Automation: How Git Decides

## 🎯 **The Magic: Git Knows Automatically**

### **Git's Built-in Decision Logic**

Git **automatically** decides to use `clean` or `smudge` based on the **direction of data flow**:

#### **Direction-Based Decision:**

```
Working Directory ← Repository = "Smudge" (decrypt)
Working Directory → Repository = "Clean" (encrypt)
```

### **Automatic Decision Table**

| Git Operation       | Data Flow      | Filter Used | Result                  |
| ------------------- | -------------- | ----------- | ----------------------- |
| `git add .env`      | Working → Repo | `clean`     | Encrypt to repository   |
| `git checkout .env` | Repo → Working | `smudge`    | Decrypt to working dir  |
| `git clone`         | Repo → Working | `smudge`    | Auto-decrypt everything |
| `git diff .env`     | Compare both   | `smudge`    | Decrypt for comparison  |
| `git pull`          | Repo → Working | `smudge`    | Auto-decrypt updates    |

### **The Automation Process**

#### **1. During `git clone`:**

```bash
# What happens internally:
1. Git reads encrypted .env from repository
2. Git sees .gitattributes: ".env filter=git-seal"
3. Git automatically calls: git-seal smudge (decrypt)
4. Creates plain text .env locally
```

#### **2. During `git add`:**

```bash
# What happens internally:
1. Git reads plain .env from working directory
2. Git sees .gitattributes: ".env filter=git-seal"
3. Git automatically calls: git-seal clean (encrypt)
4. Stores encrypted .env in repository
```

### **Why Git Knows Automatically**

**Git was designed with filters in mind:**

- **Clean filter**: Strips/transforms data for storage
- **Smudge filter**: Restores/transforms data for use
- **Direction determines**: Which filter to call
- **No manual selection**: Git handles everything

### **The Configuration Chain**

1. **`.gitattributes`** → "This file needs filtering"
2. **Git config** → "Use git-seal for filtering"
3. **Git operation** → "Which direction is data flowing?"
4. **Automatic decision** → "Use clean or smudge accordingly"

## 🎉 **Result: Zero Manual Commands**

You never need to tell Git "now encrypt" or "now decrypt" - it's all automatic based on the operation you're performing!
