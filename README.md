# _bamboo_

### What is it?
It is a little game, that I'm working on, as a hobby. I started it when I was learning Go.  
It is a top-down sandbox game, where you are stranded on an island, and need to survive. Craft, kill, build, make yourself home.  
Most of ideas and game mechanics were inspired by [Minecraft](https://minecraft.net/).  

### How to play it?
Clone the repository, and run `go run .` in the terminal.   
[Go](https://go.dev/) must be installed for this to work, of course.

### Progress
See [FEATURES.md](FEATURES.md)

## Naming conventions
You may notice, that some functions are suffixed with `B` (example: `world.ChunkAtB`), and some are not.  
The reason is pretty simple: functions with suffix `B` accept block coordinates, and functions without that suffix accept chunk coordinates.