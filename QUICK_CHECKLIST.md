# Quick Checklist - Documentation Updates

## 🎯 Critical Path (30 minutes)

### Step 1: Screenshots (15 min)
```bash
# Hero screenshots
bf                              # Screenshot this → site/images/hero-default.png
bf --theme dracula             # Screenshot this → site/images/hero-dracula.png

# Config wizard
bf --config-wizard             # Screenshot this → site/images/config-wizard.png

# Solana example
bf -s DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263  # Screenshot → site/images/solana-example.png
```

- [ ] hero-default.png
- [ ] hero-dracula.png
- [ ] config-wizard.png
- [ ] solana-example.png

---

### Step 2: Replace README (2 min)
```bash
cd ~/Documents/bubblefetch
cp README.md README-OLD.md       # Backup
mv README-NEW.md README.md       # Use new version
```

- [ ] Backed up old README
- [ ] New README in place
- [ ] Reviewed for project-specific changes

---

### Step 3: Update README Image Links (5 min)
Edit `README.md` and replace placeholders:
- Line ~28: Replace `*(Manual: Add screenshot)*` with actual image link
- Line ~108: Replace `*(Manual: Add screenshot)*` with wizard screenshot

- [ ] Hero image links updated
- [ ] Wizard screenshot link updated

---

### Step 4: Verify Performance (3 min)
```bash
bf --benchmark --format json > benchmark.json
cat benchmark.json | grep average_ms
```

If not ~1.3ms, update these lines:
- README.md line 63, 281
- Landing page performance claim

- [ ] Performance verified
- [ ] Numbers updated if needed

---

### Step 5: Landing Page Priority (5 min)
Update your landing page with:
- [ ] Hero demo section (top of page)
- [ ] Comparison table
- [ ] Link to new README sections

---

## ✅ Done!

Critical updates complete. Website now has:
- ✨ Instant visual proof
- 📊 Clear value proposition
- 🚀 Copy-paste examples
- 🔒 Privacy transparency

---

## 🎨 Optional: Full Theme Gallery (45 min)

If you want the complete theme showcase:

```bash
# Take 8 theme screenshots
bf --theme default         # → site/images/themes/default.png
bf --theme dracula         # → site/images/themes/dracula.png
bf --theme nord            # → site/images/themes/nord.png
bf --theme gruvbox         # → site/images/themes/gruvbox.png
bf --theme tokyo-night     # → site/images/themes/tokyo-night.png
bf --theme solarized-dark  # → site/images/themes/solarized-dark.png
bf --theme monokai         # → site/images/themes/monokai.png
bf --theme minimal         # → site/images/themes/minimal.png
```

Then update README.md theme section (~line 400-520) with image links.

---

## 📱 Social Media Assets

Quick wins for announcements:

```bash
# Comparison table screenshot
# (Open README.md, screenshot the table around line 76)

# Solana feature showcase
bf -s DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263 -t dracula

# Image export showcase
bf -o showcase.png
```

- [ ] Comparison table screenshot
- [ ] Solana feature screenshot
- [ ] Export showcase ready

---

## 🔍 Final Check

Before committing:
- [ ] All critical screenshots added
- [ ] README.md replaced
- [ ] Image links work
- [ ] Performance number consistent
- [ ] Landing page updated
- [ ] Commit message ready

```bash
git add .
git commit -m "docs: major documentation overhaul with visual proof

- Add hero demo section with screenshots
- Add comparison table vs fastfetch/neofetch
- Add common setups with ready-to-run configs
- Add privacy & safety section
- Add troubleshooting FAQ
- Restructure themes as gallery
- Showcase outputs section
- Unify performance claims to ~1.3ms
- Add Solana token feature documentation
- Add project status badge

All P0-P3 items from change sheet implemented."
```

---

## 📦 Files Changed

```
Modified:
  README.md (major rewrite)
  docs/CHANGELOG.md

New:
  docs/examples/config-minimal-laptop.yaml
  docs/examples/config-streamer-flex.yaml
  docs/examples/config-theme-author.yaml
  docs/examples/bubblefetch-default.svg
  docs/examples/bubblefetch-default.png
  docs/examples/bubblefetch-default.html

Screenshots Needed:
  site/images/hero-default.png
  site/images/hero-dracula.png
  site/images/config-wizard.png
  site/images/solana-example.png
```

---

**Ready to ship! 🚀**
