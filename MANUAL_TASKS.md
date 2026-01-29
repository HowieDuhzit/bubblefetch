# Manual Tasks for Documentation Updates

## 📸 Screenshots & Visual Assets Needed

### P0 - Critical (Hero Demo)

1. **Hero Section Screenshots**
   - [ ] Take screenshot of `bf` (default theme output)
   - [ ] Take screenshot of `bf --theme dracula`
   - [ ] Save to `site/images/hero-default.png`
   - [ ] Save to `site/images/hero-dracula.png`
   - [ ] Optional: Create animated GIF showing theme switching
   - [ ] Update README-NEW.md line 28-29 to reference actual images

2. **What's New Section Proof**
   - [ ] Take screenshot of config wizard (`bf --config-wizard`)
   - [ ] Save to `site/images/config-wizard.png`
   - [ ] Take screenshot of Solana token output
   - [ ] Save output to `docs/examples/solana-bonk.txt`
   - [ ] Update README-NEW.md line 108 with actual screenshot link

### P2 - Medium Priority (Theme Gallery)

3. **Theme Previews**
   For each theme, take a screenshot and save to `site/images/themes/`:
   - [ ] `bf --theme default` → `themes-default.png`
   - [ ] `bf --theme dracula` → `themes-dracula.png`
   - [ ] `bf --theme nord` → `themes-nord.png`
   - [ ] `bf --theme gruvbox` → `themes-gruvbox.png`
   - [ ] `bf --theme tokyo-night` → `themes-tokyo-night.png`
   - [ ] `bf --theme solarized-dark` → `themes-solarized-dark.png`
   - [ ] `bf --theme monokai` → `themes-monokai.png`
   - [ ] `bf --theme minimal` → `themes-minimal.png`
   - [ ] Update README-NEW.md theme section with actual image references

---

## 📝 Documentation Updates

### Files to Replace

1. **README.md**
   - [ ] Back up current README.md: `cp README.md README-OLD.md`
   - [ ] Replace with new version: `mv README-NEW.md README.md`
   - [ ] Review and adjust any project-specific details

2. **Landing Page** (`index.html` or equivalent)
   - [ ] Apply same P0-P3 changes to landing page
   - [ ] Add hero demo section with screenshots
   - [ ] Add "Why Bubblefetch" comparison table
   - [ ] Add common setups section
   - [ ] Add privacy section
   - [ ] Add troubleshooting FAQ
   - [ ] Unify performance number to ~1.3ms

3. **CHANGELOG.md**
   - [ ] Update any references to collection time to use "~1.3ms"
   - [ ] Ensure consistency across all versions

---

## 🔧 Configuration Examples

4. **Create Example Configs**
   Save these to `docs/examples/`:
   - [ ] `config-minimal-laptop.yaml` (from Common Setups)
   - [ ] `config-streamer-flex.yaml` (from Common Setups)
   - [ ] `config-theme-author.yaml` (from Common Setups)

---

## 🌐 Website Updates

5. **Landing Page Priority**
   - [ ] Add hero demo section (1 command → 1 result)
   - [ ] Add "Why Bubblefetch" table
   - [ ] Add "What's New" with visual proof
   - [ ] Add privacy section
   - [ ] Add common setups
   - [ ] Add themes gallery with copy-paste
   - [ ] Add outputs showcase
   - [ ] Add troubleshooting FAQ
   - [ ] Add project status badge

---

## 📊 Performance Verification

6. **Verify Performance Claims**
   - [ ] Run `bf --benchmark --format json > benchmark-results.json`
   - [ ] Verify median is ~1.3ms
   - [ ] If different, update all references in:
     - README.md (line 63, 281)
     - Landing page
     - docs/PERFORMANCE.md
   - [ ] Document benchmark methodology in docs/PERFORMANCE.md

---

## 🎨 Asset Generation (Already Done)

✅ SVG export example: `docs/examples/bubblefetch-default.svg`
✅ PNG export example: `docs/examples/bubblefetch-default.png`
✅ HTML export example: `docs/examples/bubblefetch-default.html`
✅ Default output text: `docs/examples/output-default.txt`

---

## 🔗 Links to Verify

7. **Internal Links**
   - [ ] Check all `docs/` links work
   - [ ] Check all theme file links work
   - [ ] Check all example file links work
   - [ ] Verify GitHub release link
   - [ ] Verify AUR package link
   - [ ] Verify landing page link

---

## 📢 Communication

8. **After Updates**
   - [ ] Tweet/announce the improved documentation
   - [ ] Update any external documentation (wiki, forum posts)
   - [ ] Submit updated package descriptions (AUR, etc.)
   - [ ] Consider adding a "What's New" banner to landing page

---

## ✅ Completion Checklist

When all above tasks are done:
- [ ] All screenshots added and linked
- [ ] README.md replaced with new version
- [ ] Landing page updated with P0-P3 changes
- [ ] Performance numbers unified across all docs
- [ ] Example configs created
- [ ] All links verified
- [ ] Documentation reviewed for consistency

---

## 🎯 Quick Priority Summary

**Do First (P0 - Highest Impact):**
1. Hero demo screenshots (2 images)
2. Replace README.md with new version
3. Config wizard screenshot
4. Update landing page hero section

**Do Next (P1 - High Priority):**
5. Verify and unify performance number everywhere
6. Add privacy section to landing page
7. Add common setups to landing page

**Do When Possible (P2-P3):**
8. Theme gallery screenshots (8 images)
9. Create example config files
10. Add troubleshooting to landing page

---

## 🤖 Automated vs Manual

**Already Automated:**
- ✅ Export examples generated
- ✅ New README content written
- ✅ Documentation structure created

**Requires Manual Work:**
- 📸 Screenshots (cannot be automated in headless environment)
- 🌐 Landing page updates (external to this repo)
- 🔗 Link verification (requires human judgment)
- 📝 Final review and project-specific adjustments
