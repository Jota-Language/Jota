<h3 align="center">
    <img src="https://raw.githubusercontent.com/mattishere/jota-assets/main/banner.png" alt="banner" width="550px" /><br>
    <i>The <b>simple language</b> for you!</i>
    <br><br>
</h3>

# 🚀 What's Jota?
Jota is a new scripting language that makes it **simple to do simple things**, aiming to offer an *easy-to-learn platform* for those who need to do something *fast* (even though the interpreter itself isn't that fast).

<details>
<summary><b>🗒️ Dev's note</b></summary>

Currently implemented in Go and based on the [Lox language (and JLox interpreter) from CraftingInterpreters](https://craftinginterpreters.com/) means it's not super fast, but it definitely gets the job done. For now, bare with a Lox semi-clone made in Go, but I believe that using Lox as the base for the language will offer a solid foundation to build upon later.

Since this is my very first time building and designing a language, the book has been a tremendous help, and I want to greatly thank the author as well as wish him the best of luck in his future endeavours!

Again, this is my first time making a language so if you have any kind of constructive criticism let me know!
</details>
<br><br>

# 🏗️ A W.I.P. Language
While Jota has the basics, it definitely still lacks a lot - from a better standard issue library, all the way to user-defined functions, classes and an import system.

Keeping that in mind, use Jota for fun and **not for any sort of serious production** *(yet!)*.

To learn about what's currently possible, feel free to check out the [examples](https://github.com/mattishere/jota/tree/main/examples)!
<br><br>

# 🔧 Usage
- `jota`: starts a REPL session.
- `jota [file.jota]`: runs a .jota file.
<br><br>

# 💾 Installation
<details>
<summary><b>🐧 Linux & Darwin (macOS)</b></summary>

- Build from source (suggested): 
```git clone https://github.com/mattishere/jota.git && cd jota && sudo make install``` *This one installs to /usr/local/bin, to install to /usr/bin you should run sudo make `global_install` instead. To uninstall, run the makefile options `uninstall` or `global_uninstall` as per your previous choice.*
- Get a binary from the [releases](https://github.com/mattishere/jota/releases).
</details>

<details>
<summary><b>🪟 Windows</b></summary>

- Get an executable from the [releases](https://github.com/mattishere/jota/releases).

**NOTE:** The executables are currently not signed, so you will likely get a **Defender** warning. Make sure to add the executable to your `PATH` to be able to use Jota in any directory.
</details>
<br><br>

# 🎨 Art & Assets
Check out the [official art repo](https://github.com/mattishere/jota-assets)!

<table>
    <tr>
        <th>Type</th>
        <th>Default</th>
        <th>Light</th>
    </tr>
    <tr>
        <td>Logo</td>
        <td><img src="https://raw.githubusercontent.com/mattishere/jota-assets/main/transparent.png" height="50px"></td>
        <td><img src="https://raw.githubusercontent.com/mattishere/jota-assets/main/transparent-light.png" height="50px" />
    </tr>
        <tr>
        <td>Logo (Circular)</td>
        <td><img src="https://raw.githubusercontent.com/mattishere/jota-assets/main/yellow-circular-bg.png" height="50px"></td>
        <td><img src="https://raw.githubusercontent.com/mattishere/jota-assets/main/yellow-circular-bg-light.png" height="50px" />
    </tr>
        <tr>
        <td>Banner</td>
        <td><img src="https://raw.githubusercontent.com/mattishere/jota-assets/main/banner.png" height="50px"></td>
        <td><img src="https://raw.githubusercontent.com/mattishere/jota-assets/main/banner-light.png" height="50px" />
    </tr>
</table>