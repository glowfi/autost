<div align="center">

# ⚡ autost

### Declarative Linux session autostart — reproducible, minimal, and scriptable

<p>
  <img src="https://img.shields.io/github/go-mod/go-version/glowfi/autost?style=for-the-badge">
  <img src="https://img.shields.io/github/license/glowfi/autost?style=for-the-badge">
  <img src="https://img.shields.io/github/stars/glowfi/autost?style=for-the-badge">
  <img src="https://img.shields.io/github/last-commit/glowfi/autost?style=for-the-badge">
</p>

<p>
  <strong>A lightweight startup manager for Linux desktops and window managers.</strong>
</p>

<p>
  Replace fragile shell startup scripts with a clean, validated YAML configuration.
</p>

</div>

---

## ✨ Features

- 🧾 **Declarative config** — YAML-based startup definition
- 🚀 **Parallel-friendly execution** via non-blocking process start
- 🔧 **Environment variable support** with expansion
- 📜 **Inline scripts** without separate files
- 🧠 **Interpreter abstraction** (`sh`, `bash`, `zsh`, etc.)
- 🖥 **XDG-compliant configuration**

---

## 🧠 Why autost?

Typical Linux startup approaches:

| Method        | Problems            |
| ------------- | ------------------- |
| `.xinitrc`    | grows messy quickly |
| WM configs    | tightly coupled     |
| shell scripts | hard to maintain    |
| cron @reboot  | not session-aware   |

`autost` provides:

✅ structure
✅ validation
✅ reproducibility
✅ separation of concerns

---

## 📦 Installation

### Build from source

```bash
git clone https://github.com/glowfi/autost
cd autost
go build -o autost .
sudo mv autost /usr/local/bin/
```

Requires:

- Go ≥ 1.25
- Linux (XDG environment)

---

## ⚙️ Configuration

`autost` looks for:

```
$XDG_CONFIG_HOME/autost/config.yaml
```

or:

```
~/.config/autost/config.yaml
```

---

### Example Config

```yaml
env:
    - EDITOR: nvim
    - XDG_CONFIG_HOME: '${HOME}/.config'

interpreter: /bin/sh

startup:
    - dunst
    - picom --backend glx

scripts:
    - |
        hour=$(date +%H)
        echo "${hour}"

    - |
        echo "hello"
```

---

## 🧩 Configuration Reference

### `env`

Environment variables exported before execution.

```yaml
env:
    - KEY: value
```

✔ Environment variables are expanded automatically:

```yaml
- CONFIG: '${HOME}/.config'
```

---

### `interpreter`

Shell used to run commands and scripts.

```yaml
interpreter: /bin/sh
```

Required when using:

- `startup`
- `scripts`

---

### `startup`

Commands executed during session startup.

```yaml
startup:
    - dunst
    - picom -b
```

Each command is executed using:

```
<interpreter> -c "<command>"
```

---

### `scripts`

Inline executable scripts.

```yaml
scripts:
    - |
        echo "Running setup"
        mkdir -p ~/.cache/app
```

Each script is:

1. Written to a temporary file
2. Marked executable
3. Executed using the interpreter

---

## 🚀 Usage

Run manually:

```bash
autost
```

Output:

```
autost: loading config from ~/.config/autost/config.yaml
autost: startup complete
```

---

## 🖥 Autostart (Desktop Integration)

Example GNOME/XDG autostart file:

```ini
[Desktop Entry]
Type=Application
Name=Autostart Program
Exec=sh -c "/usr/local/bin/autost"
X-GNOME-Autostart-enabled=true
```

Place in:

```
~/.config/autostart/
```

---

## 🧱 Execution Model

```
Load Config
     ↓
Validate
     ↓
Set Environment Variables
     ↓
Start Commands (non-blocking)
     ↓
Execute Scripts
     ↓
Wait for shutdown signal
```

`autost` listens for:

- `SIGINT`
- `SIGTERM`

and shuts down gracefully.

---

## ✅ Validation Guarantees

`autost` prevents common configuration mistakes:

- ❌ Empty environment values
- ❌ Missing interpreter when commands exist
- ❌ Invalid YAML

Errors fail fast at startup.

---

## 🧪 Testing

```bash
go test ./...
```

Includes:

- config parsing tests
- validation checks
- environment expansion tests

---

## 📁 Project Structure

```
internal/
 ├── config/     # YAML parsing & validation
 └── executor/   # process execution engine
```

Design goals:

- small surface area
- composable packages
- easy extension

---

## 🔒 Security Notes

- Scripts execute with **user privileges**
- No privilege escalation
- Temporary scripts created in `/tmp`
- Environment variables validated before export

---

## 🛠 Roadmap

- [ ] Parallel execution groups
- [ ] Dependency ordering
- [ ] Restart policies
- [ ] Logging levels
- [ ] Dry-run mode
- [ ] systemd user integration

---

## 🤝 Contributing

Contributions welcome.

```bash
fork → branch → PR
```

Guidelines:

- keep dependencies minimal
- prefer stdlib
- add tests for behavior changes

---

## 📜 License

GPLv3 — see [LICENSE](LICENSE).

---

## ⭐ Philosophy

`autost` follows a simple idea:

> Desktop startup should be **declarative**, **predictable**, and **boringly reliable**.

No magic.
No hidden state.
Just configuration → execution.
