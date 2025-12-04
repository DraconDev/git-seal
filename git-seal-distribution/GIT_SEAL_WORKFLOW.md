# Git-Seal: Perfect Safe + Flexible Workflow

## 🎯 **The Ideal Setup**

### **1. Safety First (.env in .gitignore)**

```bash
# In .gitignore
.env
.env.*
config.json
```

### **2. Git-Seal Configuration (.gitattributes)**

```bash
# In .gitattributes
.env filter=git-seal diff=git-seal
*.env filter=git-seal diff=git-seal
config.json filter=git-seal diff=git-seal
```

### **3. The Perfect Workflow**

#### **Default State (Safe):**

- `.env` ignored by Git ✅
- Plain text locally ✅
- No risk of accidental commit ✅

#### **When You Want to Commit (Flexible):**

```bash
# Override gitignore (explicit action)
git add -f .env

# Git automatically:
# 1. Calls: git-seal clean (encrypts)
# 2. Stores: Encrypted binary in Git
# 3. Result: Secure commit!
```

#### **On Clone/Fresh Checkout:**

```bash
# Git automatically:
# 1. Calls: git-seal smudge (decrypts)
# 2. Creates: Plain text .env locally
# 3. Result: Instant access!
```

## ✅ **Benefits of This Approach**

1. **🔒 Safe by Default**: Can't accidentally commit plain text
2. **🔓 Flexible**: Override with `git add -f` when needed
3. **🤖 Automatic**: Encryption/decryption is transparent
4. **🎯 Simple**: Uses familiar Git patterns
5. **🛡️ Defense in Depth**: .gitignore + encryption

## 📋 **Complete Setup Commands**

```bash
# 1. Add to .gitignore (safety)
echo ".env" >> .gitignore

# 2. Configure git-seal (encryption)
git-seal setup
echo ".env filter=git-seal diff=git-seal" >> .gitattributes

# 3. Development (plain text locally)
echo "DATABASE_URL=..." > .env

# 4. When ready to commit (encrypted)
git add -f .env
git commit -m "Add encrypted config"

# 5. Clone anywhere (auto-decrypts)
git clone <repo>
# .env is already plain text!
```

This workflow gives you the **best of both worlds**: maximum safety with maximum flexibility!
