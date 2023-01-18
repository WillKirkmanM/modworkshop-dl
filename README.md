# modworkshop-dl
<p align=center>
  <img src=https://upload.wikimedia.org/wikipedia/commons/thumb/d/d8/Payday2-logo.png/1200px-Payday2-logo.png>
</p>

## ℹ A Command-Line Utility Tool for Installing Mods from [Mod Workshop](https://modworkshop.net/)

### Get Started
1. Create / Open the [`modlist.txt`](https://github.com/WillKirkmanM/modworkshop-dl/blob/main/modlist.txt) file (In the same directory as the executable).
2. Paste the desired mods in the "`Mods`" header
3. Paste the desired assets in the "`Assets`" header
4. Start the tool
```
$ ./modworkshop-dl
```
5. Launch `PAYDAY 2`

### Building
To build the files run the command

  ### How does it work?
  - Web Scraping with [Colly](http://go-colly.org/)
  - Downloading with [Grab](https://github.com/cavaliergopher/grab)
  - Interactive Terminal with [Uilive](https://github.com/gosuri/uilive)
  - Unarchiving .zip / .rar / .tar with [Archiver v3](https://github.com/mholt/archiver)


### What I've Learned
- The "Fundahmentals" of Golang ⏩
- Command Line Tooling (How they are made) 💿
- Web Scraping ✨
