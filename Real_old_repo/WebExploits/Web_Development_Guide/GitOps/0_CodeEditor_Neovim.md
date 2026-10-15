# NEO VIM GUIDE : For dummies like me 

`writting this guide as it's my choice of code editor being one of the best and fastest for terminal lovers like me :)`

Neovim (`nvim`) is a modern, actively-developed fork of Vim. Same core
philosophy — fast, keyboard-only, no mouse needed — but with better plugin
support, built-in Lua scripting, and async features Vim doesn't have. If
you've never touched Vim/Neovim before, this guide gets you from "how do I
even close this thing" to comfortable daily use.
 
---
 
## 1. Installing Neovim
 
**macOS**
```bash
brew install neovim
```
 
**Ubuntu / Debian**
```bash
sudo apt update
sudo apt install neovim
```
(For a newer version than your distro ships, use the official PPA or
download a prebuilt binary from the Neovim GitHub releases page.)
 
**Windows**
```powershell
winget install Neovim.Neovim
```
 
**Verify it installed:**
```bash
nvim --version
```
 
---
 
## 2. The Big Idea You Must Understand First: Modes
 
This is the single biggest mental shift coming from editors like VS Code or
Notepad. Neovim is a **modal editor** — the same key does different things
depending on what mode you're in.
 
| Mode | What it's for | How to enter it |
|---|---|---|
| **Normal** | Moving around, deleting, copying — your "home base" | `Esc` (from anywhere) |
| **Insert** | Actually typing text, like a normal editor | `i`, `a`, `o` (from Normal) |
| **Visual** | Selecting text | `v`, `V`, `Ctrl-v` (from Normal) |
| **Command-line** | Running commands like save/quit/search-replace | `:` (from Normal) |
 
**You spend most of your time in Normal mode**, not Insert mode. This feels
backwards at first — in every other editor you're always "typing." In
Neovim, you drop into Insert mode just long enough to type something, then
hit `Esc` to go back to Normal mode to move around, delete, or fix things.
 
> **Rule #1 for beginners:** if you're ever confused about what's happening
> or stuck, hit `Esc`. That always gets you back to Normal mode, which is
> the safe default state.
 
---
 
## 3. Opening and Closing Neovim
 
```bash
nvim               # open with no file (or a file browser if you type nvim .)
nvim myfile.txt    # open a specific file
nvim .             # open the built-in file explorer in the current folder
```
 
**Saving and quitting** (type these in Normal mode, starting with `:`):
 
| Command | What it does |
|---|---|
| `:w` | Save (write) |
| `:q` | Quit (fails if you have unsaved changes) |
| `:wq` or `:x` | Save and quit |
| `:q!` | Quit **without** saving — discard changes |
| `:wq!` | Force save and quit |
 
> **Rule #2 for beginners:** if you're panicking and just want out without
> saving, `Esc` then `:q!` then `Enter` always works.
 

`Dev note: you can run any command in the vim too use :vsp | term <command>`
---
 
## 4. Your First Five Minutes — Minimum Viable Workflow
 
1. `nvim test.txt` — open a file
2. `i` — enter Insert mode
3. Type some text normally
4. `Esc` — back to Normal mode
5. `:wq` — save and quit
That's genuinely enough to survive. Everything below just makes you faster.
 
---
 
## 5. Movement (Normal Mode)
 
Arrow keys work, but the "home row" keys are faster once they're muscle
memory — your fingers never leave the keyboard's center.
 
| Key | Moves |
|---|---|
| `h` `j` `k` `l` | left, down, up, right |
| `w` | forward one word |
| `b` | back one word |
| `e` | end of current/next word |
| `0` | start of the line |
| `^` | first non-blank character of the line |
| `$` | end of the line |
| `gg` | top of the file |
| `G` | bottom of the file |
| `:42` or `42G` | jump to line 42 |
| `{` / `}` | jump a paragraph back / forward |
| `Ctrl-d` / `Ctrl-u` | half-page down / up |
| `f x` | jump forward to next `x` on the line |
| `%` | jump to matching bracket `(){}[]` |
 
A number before any movement repeats it — `5j` moves down 5 lines, `3w`
jumps 3 words forward.
 
---
 
## 6. Editing (Normal Mode)
 
| Key | Action |
|---|---|
| `x` | delete the character under the cursor |
| `dd` | delete (cut) the whole line |
| `dw` | delete to the end of the word |
| `d$` | delete to the end of the line |
| `cc` | change (delete + enter Insert mode on) the whole line |
| `cw` | change to end of word |
| `yy` | yank (copy) the whole line |
| `p` | paste after the cursor |
| `P` | paste before the cursor |
| `u` | undo |
| `Ctrl-r` | redo |
| `.` | repeat the last change — extremely useful |
| `>>` / `<<` | indent / un-indent the line |
| `J` | join the next line onto this one |
 
**The pattern to notice:** most edit commands are `verb + motion`. `d` (delete)
+ `w` (word) = delete a word. `d` + `$` (end of line) = delete to end of
line. `c` + `w` = change a word. Once this clicks, you can compose commands
instead of memorizing hundreds of them.
---
 
## 7. Visual Mode (Selecting Text)
 
| Key | Selects |
|---|---|
| `v` | character-by-character, starting here |
| `V` | whole lines |
| `Ctrl-v` | a rectangular block (great for editing columns) |
 
While selected, press `d` to delete it, `y` to copy it, or `c` to delete
and start typing a replacement.
 
---
 
## 8. Search and Replace
 
| Command | Does |
|---|---|
| `/word` then `Enter` | search forward for "word" |
| `?word` then `Enter` | search backward for "word" |
| `n` | jump to next match |
| `N` | jump to previous match |
| `:%s/old/new/g` | replace **every** "old" with "new" in the whole file |
| `:%s/old/new/gc` | same, but ask you to confirm each one (`y`/`n`) |
| `:s/old/new/g` | replace on the current line only |
 
---
 
## 9. Multiple Files, Splits, and Tabs
 
```
:e otherfile.txt     open another file
:bnext / :bprev      cycle through open files (buffers)
:ls                  list open buffers
:bd                  close current buffer
 
:sp                  horizontal split
:vsp                 vertical split
Ctrl-w  then h/j/k/l move between splits
 
:tabnew file.txt     open file in a new tab
gt / gT              next tab / previous tab
```
 
---
 
## 10. Copying to Your System Clipboard
 
By default, Neovim's yank/paste (`y`/`p`) uses its own internal register,
**not** your OS clipboard. Two ways to bridge that:
 
**One-off:** prefix with `"+`
```
"+y     yank into the system clipboard
"+p     paste from the system clipboard
```
 
**Always-on:** add this to your config (see next section):
```lua
vim.opt.clipboard = "unnamedplus"
```
(On Linux you may also need `xclip` or `wl-clipboard` installed for this to
work.)
 
---
 
## 11. Basic Configuration
 
Neovim's config file lives at:
 
- Linux/macOS: `~/.config/nvim/init.lua`
- Windows: `~/AppData/Local/nvim/init.lua`
Create the folder/file if it doesn't exist yet. Here's a solid, minimal
starting config with common beginner-friendly options:
 
```lua
-- ~/.config/nvim/init.lua
 
vim.opt.number = true          -- show line numbers
vim.opt.relativenumber = true  -- line numbers relative to your cursor
vim.opt.tabstop = 2            -- width of a tab character
vim.opt.shiftwidth = 2         -- indent width
vim.opt.expandtab = true       -- use spaces instead of tabs
vim.opt.smartindent = true     -- auto-indent new lines sensibly
vim.opt.wrap = false           -- don't visually wrap long lines
vim.opt.ignorecase = true      -- case-insensitive search...
vim.opt.smartcase = true       -- ...unless you type a capital letter
vim.opt.clipboard = "unnamedplus" -- share clipboard with your OS
vim.opt.termguicolors = true   -- better colors in the terminal
 
vim.g.mapleader = " "          -- set space as your <leader> key
```
 
That `<leader>` key is important: it's a prefix key you can bind your own
shortcuts to, without colliding with Neovim's built-in commands. Example:
 
```lua
-- Space + w saves the file
vim.keymap.set("n", "<leader>w", ":w<CR>")
```
 
---
 
## 12. Plugins (Once You're Comfortable)
 
You don't need plugins to be productive, but they're where Neovim really
shines. The current standard plugin manager is
[**lazy.nvim**](https://github.com/folke/lazy.nvim) — it lazy-loads
plugins so startup stays fast.
 
A tiny example, once lazy.nvim itself is installed per its README:
 
```lua
require("lazy").setup({
  { "nvim-telescope/telescope.nvim" }, -- fuzzy finder for files/text
  { "nvim-treesitter/nvim-treesitter" }, -- better syntax highlighting
  { "neovim/nvim-lspconfig" },          -- language server support (autocomplete, go-to-def)
})
```
 
**Honest beginner advice:** rather than hand-building a config from
scratch, most newcomers get further faster by starting from a
well-documented, teachable base:
 
- **[kickstart.nvim](https://github.com/nvim-lua/kickstart.nvim)** — a
  single, heavily-commented config file meant to be read and understood,
  not just installed. Best pick if you actually want to learn.
- **[LazyVim](https://github.com/LazyVim/LazyVim)** — a full "batteries
  included" pre-configured distribution, more like an IDE out of the box.
---
 
## 13. Quick Cheat Sheet
 
| Task | Key |
|---|---|
| Enter Insert mode | `i` |
| Back to Normal mode | `Esc` |
| Save | `:w` |
| Save and quit | `:wq` |
| Quit without saving | `:q!` |
| Undo / Redo | `u` / `Ctrl-r` |
| Delete a line | `dd` |
| Copy a line | `yy` |
| Paste | `p` |
| Delete a word | `dw` |
| Select lines | `V` |
| Search | `/text` |
| Replace all | `:%s/old/new/g` |
| Jump to top / bottom | `gg` / `G` |
| Repeat last change | `.` |
 
---
 
## 14. Common Beginner Mistakes
 
- **"I typed a bunch of letters and weird stuff happened."** You were
  probably in Normal mode and those letters are commands, not text. Hit
  `Esc`, look at the bottom-left corner (it should say `-- INSERT --` only
  when you're actually in Insert mode), and try again.
- **"I can't type at all."** You're in Normal mode. Press `i` first.
- **Don't try to memorize this whole guide at once.** Learn `i`, `Esc`,
  `:wq`, `hjkl`, `dd`, `yy`, `p`, `u` first. Everything else layers on
  gradually as muscle memory frees up attention.
- **Try the built-in interactive tutorial.** Run `vimtutor` from your
  terminal (separate from `nvim` itself) — it's a 25–30 minute hands-on
  walkthrough of exactly this material, and it's genuinely the fastest way
  to build muscle memory.
---
 
## 15. Where to Go Next
 
- Type `:help` inside Neovim any time — it's a real, searchable manual,
  not a stub.
- [neovim.io](https://neovim.io) — official site and docs.
- Once comfortable, look into a language server (`nvim-lspconfig`) for
  autocomplete and go-to-definition in whatever language you write.
You don't need to "finish" learning Neovim before it's useful — the
5-minute workflow in section 4 is a complete, valid way to use it from day
one. Everything else is optimization you pick up as friction shows up.
 

## 16. Quick notes form the Dev {me :)}

### Here is how to upgrade you vim in few commands 
```bash
rm -rf ~/.config/nvim && git clone https://github.com/LazyVim/starter ~/.config/nvim
rm -rf ~/.config/nvim/.git
nvim
```
**now you got lazy nvim :)**

Here are some extensions you should get so far 
1. using github in nvim -> terminal(manual) + lazygit(optional)
2. Wanna open a terminal & hacking -> terminal 
3. Syntax highlighting in coding -> lang in lazyExtras 

