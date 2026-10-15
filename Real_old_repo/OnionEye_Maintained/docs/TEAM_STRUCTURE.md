# OnionEye Team Structure

> **Last Updated:** September 9, 2026

This document defines the team structure for the OnionEye project. All other documentation references teams by name rather than individual members to maintain flexibility.

---

## 📋 Team Roster

### 🤖 AI Team
**Responsibilities:** Machine learning models, AI integration, intelligent analysis
- Jinesh
- Rudra

### 💻 Backend & Frontend Team
**Responsibilities:** Service development, API design, UI/UX implementation
- Shivansh
- Chanchal

### 🔍 Live Testing & Debugging Team
**Responsibilities:** Quality assurance, bug fixes, production testing, system monitoring
- Shubham
- Reshit

---

## 🔄 Team Responsibilities Matrix

| Area | Primary Team | Support Team |
|------|--------------|--------------|
| AI/ML Models | AI Team | - |
| Correlation Engine | AI Team | Backend & Frontend Team |
| Content Classification | AI Team | Backend & Frontend Team |
| Backend Services | Backend & Frontend Team | - |
| Frontend Dashboard | Backend & Frontend Team | - |
| API Development | Backend & Frontend Team | - |
| Testing & QA | Live Testing & Debugging Team | All Teams |
| Production Issues | Live Testing & Debugging Team | Relevant Team |
| CI/CD Pipeline | Live Testing & Debugging Team | Backend & Frontend Team |
| Documentation | All Teams | - |

---

## 📞 Contact & Escalation

### For Questions About:
- **AI/ML Features** → Contact AI Team
- **Backend/Frontend Issues** → Contact Backend & Frontend Team
- **Bugs/Testing** → Contact Live Testing & Debugging Team

### Escalation Path:
1. Check documentation first
2. Contact relevant team
3. If cross-team issue, coordinate through team leads

---

## 📝 How to Update This File

When team membership changes:
1. Update this file (`TEAM_STRUCTURE.md`)
2. No need to update other documentation (they reference teams, not individuals)
3. Commit with message: `docs: update team structure`

---

## 🔗 Referenced In:
- `README.md`
- `CONTRIBUTING.md`
- `CODEOWNERS`
- All testing documentation
- All architecture documentation

---

**Note:** Documentation should always reference "AI Team", "Backend & Frontend Team", or "Live Testing & Debugging Team" instead of individual names to maintain flexibility as the team evolves.
