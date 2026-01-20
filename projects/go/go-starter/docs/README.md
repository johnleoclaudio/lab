# Documentation Guide

Welcome to the AI-assisted development documentation for this Go backend project. This folder contains comprehensive guides for working with AI coding agents to maximize productivity.

## 📚 Documentation Structure

```
docs/
├── README.md                 # This file - start here
├── AGENTS.md                 # Agent roles and responsibilities
├── SKILLS.md                 # Reusable capabilities catalog
├── QUICK_START.md            # Ready-to-use prompts
├── ARCHITECTURE.md           # System design and structure
├── CODING_STANDARDS.md       # Go conventions (from your original)
├── DATABASE.md               # sqlc and migration guidelines (from your original)
├── API_STANDARDS.md          # REST and JSON:API spec (from your original)
├── TESTING.md                # Testing strategies (from your original)
└── SECURITY.md               # Security best practices (from your original)
```

## 🚀 Quick Start

### For New Developers
1. Read **QUICK_START.md** for common prompts
2. Skim **AGENTS.md** to understand agent roles
3. Reference **ARCHITECTURE.md** for system design
4. Check **CODING_STANDARDS.md** for Go conventions

### For AI Agent Interaction
1. Identify your task type in **QUICK_START.md**
2. Copy the relevant prompt template
3. Customize with your specific requirements
4. Reference **SKILLS.md** for advanced capabilities

### For Code Review
1. Use **Code Reviewer Agent** from **AGENTS.md**
2. Check against **CODING_STANDARDS.md**
3. Verify **DATABASE.md** compliance (sqlc usage)
4. Review **SECURITY.md** checklist

## 📖 Document Purposes

### AGENTS.md
**Purpose**: Define AI agent roles and workflows  
**Use When**: 
- Starting a new feature
- Need to understand which agent to use
- Planning agent collaboration

**Key Sections**:
- Agent roles and responsibilities
- Pre-generation checklists
- Collaboration workflows
- Agent selection guide

### SKILLS.md
**Purpose**: Catalog of reusable capabilities  
**Use When**:
- Need specific code generation
- Want to understand available tools
- Building complex workflows

**Key Sections**:
- Code generation skills
- Database skills
- Testing skills
- Review skills
- Infrastructure skills

### QUICK_START.md
**Purpose**: Ready-to-use prompt templates  
**Use When**:
- Want to get started immediately
- Need prompt examples
- Don't want to read full documentation

**Key Sections**:
- Common task prompts
- Agent selection matrix
- Best practices for prompts
- Troubleshooting guide

### ARCHITECTURE.md
**Purpose**: System design and structure  
**Use When**:
- Understanding project organization
- Designing new features
- Reviewing architectural decisions

**Key Sections**:
- Layered architecture
- Dependency flow
- Layer responsibilities
- Context propagation
- Error handling strategy

### CODING_STANDARDS.md
**Purpose**: Go coding conventions  
**Use When**:
- Writing new code
- Reviewing code
- Onboarding team members

**Key Sections**:
- Naming conventions
- Error handling
- Logging standards
- Code organization

### DATABASE.md
**Purpose**: Database and sqlc guidelines  
**Use When**:
- Creating migrations
- Writing queries
- Setting up database

**Key Sections**:
- sqlc usage (CRITICAL: no manual SQL)
- Migration patterns
- Query definitions
- Transaction management

### API_STANDARDS.md
**Purpose**: REST and JSON:API specification  
**Use When**:
- Designing API endpoints
- Formatting responses
- Handling errors

**Key Sections**:
- RESTful principles
- JSON:API format
- Status codes
- Versioning

### TESTING.md
**Purpose**: Testing strategies and patterns  
**Use When**:
- Writing tests
- Reviewing test coverage
- Setting up test infrastructure

**Key Sections**:
- Table-driven tests
- Mocking strategies
- Integration testing
- Coverage requirements

### SECURITY.md
**Purpose**: Security best practices  
**Use When**:
- Implementing authentication
- Handling sensitive data
- Reviewing for vulnerabilities

**Key Sections**:
- Input validation
- Password handling
- SQL injection prevention
- Authentication/authorization

## 🎯 Common Workflows

### Workflow 1: Add New Feature
```
1. Read QUICK_START.md → Find "Add Complete CRUD Resource" template
2. Customize prompt with your resource details
3. Reference AGENTS.md → Go Backend Developer + Database Engineer
4. Check ARCHITECTURE.md → Understand layer responsibilities
5. Verify CODING_STANDARDS.md → Follow conventions
6. Review DATABASE.md → Ensure sqlc compliance
```

### Workflow 2: Code Review
```
1. Use AGENTS.md → Code Reviewer Agent
2. Apply CODING_STANDARDS.md checklist
3. Verify DATABASE.md compliance (sqlc)
4. Check SECURITY.md for vulnerabilities
5. Confirm TESTING.md coverage requirements
```

### Workflow 3: Database Changes
```
1. Use AGENTS.md → Database Engineer Agent
2. Follow DATABASE.md migration patterns
3. Reference ARCHITECTURE.md → Repository layer
4. Update queries using sqlc
5. Test with TESTING.md guidelines
```

## 💡 How to Use These Docs with AI

### Basic Pattern
```
1. Identify task type
2. Find relevant agent in AGENTS.md
3. Copy prompt template from QUICK_START.md
4. Customize for your needs
5. Reference supporting docs as needed
```

### Advanced Pattern
```
1. Review AGENTS.md → Understand agent collaboration
2. Check SKILLS.md → Identify needed capabilities
3. Build custom prompt combining skills
4. Reference architecture/standards docs
5. Iterate based on output
```

### Example Interaction
```
You: "I need to add a comments feature to my blog"

Step 1: Check QUICK_START.md
→ Find "Add Complete CRUD Resource" template

Step 2: Customize prompt
→ Replace [resource_name] with "comments"
→ Add specific fields

Step 3: Submit to AI with agent context
→ Agent: Go Backend Developer + Database Engineer
→ Skills: scaffold_crud_resource

Step 4: Review output against standards
→ CODING_STANDARDS.md: Check naming, error handling
→ DATABASE.md: Verify sqlc usage
→ TESTING.md: Confirm test coverage
```

## 🔧 Maintenance

### When to Update These Docs

**AGENTS.md**: 
- Add new agent types
- Update agent responsibilities
- Modify workflows

**SKILLS.md**:
- Add new capabilities
- Update skill templates
- Deprecate outdated patterns

**QUICK_START.md**:
- Add new common use cases
- Update prompt templates
- Include troubleshooting tips

**Other Docs**:
- Update when project standards change
- Add new patterns and practices
- Include lessons learned

### Keeping Docs in Sync
- Review docs when merging major features
- Update examples to match current codebase
- Deprecate outdated patterns clearly
- Version docs with major releases

## 📊 Document Priority by Role

### New Developer
**Priority Order**:
1. QUICK_START.md (get productive fast)
2. ARCHITECTURE.md (understand system)
3. CODING_STANDARDS.md (write correct code)
4. DATABASE.md (avoid common mistakes)

### Experienced Developer
**Priority Order**:
1. AGENTS.md (leverage AI effectively)
2. SKILLS.md (advanced capabilities)
3. QUICK_START.md (common tasks)
4. Reference docs as needed

### Code Reviewer
**Priority Order**:
1. CODING_STANDARDS.md (consistency)
2. DATABASE.md (sqlc compliance)
3. SECURITY.md (vulnerabilities)
4. TESTING.md (coverage)

### Tech Lead / Architect
**Priority Order**:
1. ARCHITECTURE.md (system design)
2. AGENTS.md (team workflows)
3. SECURITY.md (compliance)
4. All others for reference

## 🎓 Learning Path

### Week 1: Foundations
- [ ] Read QUICK_START.md completely
- [ ] Understand ARCHITECTURE.md layers
- [ ] Practice with simple prompts
- [ ] Review CODING_STANDARDS.md

### Week 2: Database & Testing
- [ ] Master DATABASE.md (sqlc is critical)
- [ ] Practice migration creation
- [ ] Study TESTING.md patterns
- [ ] Write tests for existing code

### Week 3: Advanced Features
- [ ] Explore AGENTS.md workflows
- [ ] Use SKILLS.md for complex tasks
- [ ] Implement complete feature
- [ ] Security review with SECURITY.md

### Week 4: Mastery
- [ ] Create custom agent workflows
- [ ] Combine multiple skills
- [ ] Contribute improvements to docs
- [ ] Mentor others using these docs

## 🤝 Contributing to Documentation

### Adding Examples
1. Test the prompt/code thoroughly
2. Add to relevant document
3. Include expected output
4. Update QUICK_START.md if generally useful

### Reporting Issues
1. Note which document is unclear
2. Explain what's confusing
3. Suggest improvements
4. Submit PR with changes

### Best Practices
- Keep examples realistic and tested
- Use consistent formatting
- Link related sections
- Update timestamps on major changes

## 📞 Getting Help

### If Documentation Unclear
1. Check QUICK_START.md for examples
2. Review ARCHITECTURE.md for context
3. Ask specific questions referencing docs
4. Suggest documentation improvements

### If AI Output Incorrect
1. Review prompt against templates
2. Check agent selection (AGENTS.md)
3. Verify against CODING_STANDARDS.md
4. Use Code Reviewer Agent to identify issues

### If Standards Conflict
1. CODING_STANDARDS.md takes precedence for Go code
2. DATABASE.md is absolute for sqlc usage
3. ARCHITECTURE.md guides design decisions
4. Escalate unresolved conflicts to team

## 🔗 Quick Links

- [AGENTS.md](./AGENTS.md) - Who does what
- [SKILLS.md](./SKILLS.md) - What can be done
- [QUICK_START.md](./QUICK_START.md) - How to start fast
- [ARCHITECTURE.md](./ARCHITECTURE.md) - How it's organized
- [CODING_STANDARDS.md](./CODING_STANDARDS.md) - How to write it
- [DATABASE.md](./DATABASE.md) - How to store it
- [API_STANDARDS.md](./API_STANDARDS.md) - How to expose it
- [TESTING.md](./TESTING.md) - How to verify it
- [SECURITY.md](./SECURITY.md) - How to protect it

---

## Summary

This documentation system is designed to maximize productivity when working with AI coding agents while maintaining high code quality and consistency. Start with **QUICK_START.md** for immediate productivity, then explore other documents as needed for deeper understanding.

Remember: 
- ✅ Always use sqlc (never manual SQL)
- ✅ Follow TDD approach
- ✅ Write comprehensive tests
- ✅ Propagate context through layers
- ✅ Handle errors with context wrapping

Happy coding! 🚀
