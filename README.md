<div align="center">
  <img src="https://github.com/user-attachments/assets/db44efe6-c75a-4790-8ba1-39632e99b13a" alt="HelloGang Preview" width="100%">
</div>

# ✨ HelloGang CLI

**HelloGang** is a beautiful, extremely fast terminal greeter written in Go. Instead of a boring, blank prompt when you open your terminal, HelloGang instantly renders a personalized, dynamic ASCII art greeting along with real-time system stats.

Powered by *Cobra* and *Lipgloss*, it features a striking "Claude Orange" aesthetic, auto-generation of your name in graffiti, and automatic shell installation.

---

## ⚡ Features
- **Dynamic ASCII Greeting**: Generates high-quality graffiti ASCII art using your actual name.
- **System Monitoring**: instantly view real-time CPU usage, RAM stats, and exact Date/Time.
- **Premium UI**: Designed with clean spacing, vibrant solid progress bars, and a sharp Claude Orange & White color palette.
- **Cross-Platform Auto-Start**: Installs seamlessly into **PowerShell, CMD, or Bash** profiles, so the greeting appears automatically whenever you open your terminal.

## 🚀 Installation

### 1. Direct Download (Easiest)
1. Go to the [Releases Tab](../../releases/latest).
2. Download the executable for your OS (e.g. `hellogang-windows-amd64.exe`).
3. Rename it to `hellogang.exe` and place it in a folder of your choice (preferably one in your system `PATH`).

### 2. Build via Go
If you have Go installed on your machine:
```bash
git clone https://github.com/GitNimay/cli-app.git
cd cli-app
go install .
```

## 🛠 Usage
Test the application instantly by running:
```bash
hellogang
```

To configure HelloGang to run automatically every time you open a terminal:
```bash
hellogang install
```
*(You will be prompted to enter your name, and HelloGang will smartly handle the rest!)*

To safely remove the startup integration later:
```bash
hellogang uninstall
```

---
*Built with [Cobra](https://github.com/spf13/cobra) and [Bubbletea/Lipgloss](https://github.com/charmbracelet/lipgloss) *
